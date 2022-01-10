package bot

import (
	"fmt"
	"PokeBot/config"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
)

var BotId string
var goBot *discordgo.Session

type PokemonInfo struct {
	Name string `json:"name"`
	Sprites SpriteUrls `json:sprites`
}

type SpriteUrls struct {
	DefaultFront string `json:"front_default"`
	FemaleFront string `json:"front_female"`
	ShinyFront string `json:"front_shiny"`
	ShinyFemaleFront string `json:"front_shiny_female"`
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is live!")
}

func messageHandler( s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if m.Content != "!pokemon" {
		return
	}

	//setup and get a random pokemon's data
	var responseObject PokemonInfo

	rand.Seed(time.Now().UnixNano())
	minPokeId, maxPokeId := 1, 898
	pokeId := rand.Intn(maxPokeId - minPokeId) + minPokeId
	pokeUrl := "https://pokeapi.co/api/v2/pokemon/"+strconv.Itoa(pokeId)

	response, err := http.Get(pokeUrl)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.Unmarshal(responseData, &responseObject)

	pokeImage := &discordgo.MessageEmbedImage{}
	if rand.Intn(4096) == 0 {
		pokeImage.URL = responseObject.Sprites.ShinyFront
	} else {
		pokeImage.URL = responseObject.Sprites.DefaultFront
	}

	pokeResponse := &discordgo.MessageEmbed{
		Image: pokeImage,
		Title: strings.Title(responseObject.Name),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, pokeResponse)
}

