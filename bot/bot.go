package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"PokeBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string

//these structs are nested, top is outer and next two are within it

//basic info, and where to find sprites/other info
type pokemonInfo struct {
	Name    string      `json:"name"`
	Height  float64     `json:"height"`
	Weight  float64     `json:"weight"`
	Sprites spriteUrls  `json:sprites`
	Species speciesUrls `json:"species"`
}

//basic sprite url links
type spriteUrls struct {
	DefaultFront     string `json:"front_default"`
	FemaleFront      string `json:"front_female"`
	ShinyFront       string `json:"front_shiny"`
	ShinyFemaleFront string `json:"front_shiny_female"`
}

//basic species url info. Note, name is redundantly stored
type speciesUrls struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//these structs are nested, top is more outer, and are unrelated to above
type speciesInfo struct {
	TextEntries []entryInfo `json:"flavor_text_entries"`
}

type entryInfo struct {
	FlavorText string   `json:"flavor_text"`
	Language   langInfo `json:"language"`
}

type langInfo struct {
	Name string `json:"name"`
	url  string `json:"url"`
}

func BotSetup(s *discordgo.Session) {
	//set the user info so the bot doesn't self reply later
	u, err := s.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	BotId = u.ID
	return
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//don't let it respond to itself
	if m.Author.ID == BotId {
		return
	}

	//if the message is not the trigger, exit fast
	if m.Content != config.BotPrefix+"pokemon" {
		return
	}

	//vars
	var (
		pokeUrl                      string
		pokemon                      pokemonInfo
		pokemonBytes                 []byte
		pokeImage                    *discordgo.MessageEmbedImage
		pokeResponse                 *discordgo.MessageEmbed
		pokeDesc                     string
		heightM, weightKg, weightLbs float64
		heightFt, heightIn           int
		shiny                        bool
	)

	rand.Seed(time.Now().UnixNano())

	//get pokemon data
	pokeUrl = getPokeUrl()

	//get response from web url and parse data
	pokemonBytes = getUrlInfo(pokeUrl)
	json.Unmarshal(pokemonBytes, &pokemon)

	//get poke image
	shiny = isShiny() //rng!
	pokeImage = getPokeImage(pokemon, shiny)
	if shiny {
		pokemon.Name += " (Shiny)"
	}

	//set up weight and height
	weightKg = pokemon.Weight / 10
	weightLbs = kgsToLbs(weightKg)
	heightM = pokemon.Height / 10
	heightFt, heightIn = mToFtIn(heightM)

	//get poke description/flavor text
	pokeDesc = getPokeDesc(pokemon.Species.Url)
	pokeDesc += fmt.Sprintf("\nWeight: %v kgs / %v lbs", weightKg, weightLbs)
	pokeDesc += fmt.Sprintf("\nHeight: %v m / %v ft %v in", heightM, heightFt, heightIn)

	//craft a response
	pokeResponse = &discordgo.MessageEmbed{
		Image:       pokeImage,
		Title:       strings.Title(pokemon.Name),
		Description: pokeDesc,
		URL:         "https://pokemondb.net/pokedex/" + pokemon.Name,
		//also have footer
	}

	//reply
	s.ChannelMessageSendEmbed(m.ChannelID, pokeResponse)
}

func getPokeUrl() string {
	//generate random num within the range of pokemon and append to url
	var (
		minPokeId, maxPokeId, pokeId int
		pokeUrl                      string
	)

	minPokeId, maxPokeId = 1, 898
	pokeId = rand.Intn(maxPokeId-minPokeId) + minPokeId
	pokeUrl = "https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(pokeId)
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

func isShiny() bool {
	chance := 4096
	return rand.Intn(chance) == 0
}

func getPokeImage(pokeInfo pokemonInfo, shiny bool) *discordgo.MessageEmbedImage {
	pokeImage := &discordgo.MessageEmbedImage{}

	if shiny {
		pokeImage.URL = pokeInfo.Sprites.ShinyFront
	} else {
		pokeImage.URL = pokeInfo.Sprites.DefaultFront
	}

	if pokeImage.URL == "" {
		fmt.Println("Couldn't find sprite for this mon")
		fmt.Println("Is shiny? ", shiny)
		fmt.Println(pokeInfo)
	}

	return pokeImage
}

func getPokeDesc(pokeSpeciesUrl string) string {
	var (
		species speciesInfo
		desc    string
	)

	speciesBytes := getUrlInfo(pokeSpeciesUrl)
	json.Unmarshal(speciesBytes, &species)

	//go in reverse order, since more recent entries are more readable
	//clean forwards way: for _, entry := range species.TextEntries {}
	for i := len(species.TextEntries) - 1; i >= 0; i-- {
		entry := species.TextEntries[i]

		if isEnglish(entry) {
			desc = entry.FlavorText
			break
		}
	}
	return desc
}

func isEnglish(entry entryInfo) bool {
	return entry.Language.Name == "en"
}

//convert kg to lbs
func kgsToLbs(weightKg float64) float64 {
	return roundTo(weightKg*2.204623, 2)
}

//shamelessly stolen from stack overflow https://stackoverflow.com/questions/52048218/round-all-decimal-points-in-golang#52048478
func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

//convert height in M to height in ft and inches. Not super accurate.
func mToFtIn(heightM float64) (int, int) {
	heightInches := int(heightM * 100 / 2.54)
	return heightInches / 12, heightInches % 12
}
