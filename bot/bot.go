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
	//init the bot
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//set the user info so the bot doesn't self reply later
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	BotId = u.ID

	//pass the bot the function to run
	goBot.AddHandler(messageHandler)

	//open the bot
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is live!")
}

func messageHandler( s *discordgo.Session, m *discordgo.MessageCreate) {
	//don't let it respond to itself
	if m.Author.ID == BotId {
		return
	}

	//if the message is not the trigger, exit fast
	if m.Content != config.BotPrefix + "pokemon" {
		return
	}

	//variable declaration
	var (
		pokeUrl string
		pokemon PokemonInfo
		pokemonBytes []byte
		pokeImage *discordgo.MessageEmbedImage
		pokeResponse *discordgo.MessageEmbed
	)

	rand.Seed(time.Now().UnixNano())

	//get pokemon data
	pokeUrl = getPokeUrl()

	//get response from web url and parse data
	pokemonBytes = getUrlInfo(pokeUrl)
	json.Unmarshal(pokemonBytes, &pokemon)

	//get poke image
	pokeImage = getPokeImage(pokemon)

	//craft a response
	pokeResponse = &discordgo.MessageEmbed{
		Image: pokeImage,
		Title: strings.Title(pokemon.Name),
	}

	//reply
	s.ChannelMessageSendEmbed(m.ChannelID, pokeResponse)
}

func getPokeUrl() string {
	//generate random num within the range of pokemon and append to url
	var (
		minPokeId, maxPokeId, pokeId int
		pokeUrl string
	)

	minPokeId, maxPokeId = 1, 898
	pokeId = rand.Intn(maxPokeId - minPokeId) + minPokeId
	pokeUrl = "https://pokeapi.co/api/v2/pokemon/"+strconv.Itoa(pokeId)
	return pokeUrl
}

func getUrlInfo(url string) []byte {
	//get the initial info
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}
	}

	//read the response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}
	}

	return responseData
}

func getPokeImage(pokeInfo PokemonInfo) *discordgo.MessageEmbedImage {

	pokeImage := &discordgo.MessageEmbedImage{}
	chance := 4096

	//generate random num, 1/chance is the possibility 
	if rand.Intn(chance) == 0 {
		pokeImage.URL = pokeInfo.Sprites.ShinyFront
	} else {
		pokeImage.URL = pokeInfo.Sprites.DefaultFront
	}

	return pokeImage
}
