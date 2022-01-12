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
	Sprites spriteUrls  `json:"sprites"`
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
	Url  string `json:"url"`
}

//Call to initialize the bot's ID so it doesn't reply to itself
func BotSetup(s *discordgo.Session) {
	//set the user info so the bot doesn't self reply later
	if u, err := s.User("@me"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		BotId = u.ID
	}
	return
}

//Handle to use to process messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//don't let it respond to itself
	if m.Author.ID == BotId {
		return
	}

	var reply interface{}
	//handle pokemon messages
	if m.Content == config.BotPrefix+"pokemon" {
		reply = replyToPokemonMessage()
	}

	//reply with a switch, handle embedded and simple messages
	switch message := reply.(type) {
	case *discordgo.MessageEmbed:
		s.ChannelMessageSendEmbed(m.ChannelID, message)
	case string:
		s.ChannelMessageSend(m.ChannelID, message)
	default:
		fmt.Printf("Unknown message type %T received\n", message)
	}

}

//process a request for a pokemon, returns an embedded message about a random pokemon
func replyToPokemonMessage() *discordgo.MessageEmbed {

	var (
		pokemon      pokemonInfo
		pokeImage    *discordgo.MessageEmbedImage
		pokeResponse *discordgo.MessageEmbed
		pokeDesc     strings.Builder
	)

	//setup for rng later
	rand.Seed(time.Now().UnixNano())
	minPokeId, maxPokeId := 1, 898 //range of pokemon to pick from
	chanceOfShiny := 4096          //chance of a shiny is 1/chanceOfShiny

	//get pokemon data
	var pokeUrl string = getPokeUrl(minPokeId, maxPokeId)

	//get response from web url and parse data
	var pokemonBytes []byte = getUrlInfo(pokeUrl)
	json.Unmarshal(pokemonBytes, &pokemon)

	//get poke image
	var shiny bool = isShiny(chanceOfShiny)
	pokeImage = getPokeImage(pokemon, shiny)
	if shiny {
		pokemon.Name += " (Shiny)"
	}

	//set up weight and height
	var weightKg, heightM float64 = pokemon.Weight / 10, pokemon.Height / 10
	var weightLbs float64 = kgsToLbs(weightKg)
	var heightFt, heightIn int = mToFtIn(heightM)

	//get poke description/flavor text
	pokeDesc.WriteString(getPokeDesc(pokemon.Species.Url))
	pokeDesc.WriteString(fmt.Sprintf("\nWeight: %v kgs / %v lbs", weightKg, weightLbs))
	pokeDesc.WriteString(fmt.Sprintf("\nHeight: %v m / %v ft %v in", heightM, heightFt, heightIn))

	//craft a response
	pokeResponse = &discordgo.MessageEmbed{
		Image:       pokeImage,
		Title:       strings.Title(pokemon.Name),
		Description: pokeDesc.String(),
		URL:         "https://pokemondb.net/pokedex/" + pokemon.Name,
		//also have footer available
	}

	return pokeResponse
}

//identify a random pokemon by ID between min and max, return a link to that pokemon
func getPokeUrl(minPokeId, maxPokeId int) string {
	//generate random num within the range of pokemon and append to url
	pokeId := getRandomIntBetween(minPokeId, maxPokeId)
	urlBase := "https://pokeapi.co/api/v2/pokemon/"
	return urlBase + strconv.Itoa(pokeId)
}

func getUrlInfo(url string) []byte {
	var (
		response     *http.Response
		responseData []byte
		err          error
	)
	//get the initial info
	if response, err = http.Get(url); err != nil {
		fmt.Println(err.Error())
		return []byte{}
	}

	//read the response
	if responseData, err = ioutil.ReadAll(response.Body); err != nil {
		fmt.Println(err.Error())
		return []byte{}
	}

	return responseData
}

//determine if a pokemon is shiny or not by rolling a die of a specific chance. 1/chance makes a shiny.
func isShiny(chance int) bool {
	return getRandomIntBetween(0, chance) == 0
}

//pull out the pokemon image and create a struct for the image and allow handling of shiny vs non-shiny
//returns a struct with the image
func getPokeImage(pokeInfo pokemonInfo, shiny bool) *discordgo.MessageEmbedImage {
	pokeImage := &discordgo.MessageEmbedImage{
		URL: pokeInfo.Sprites.DefaultFront,
	}

	if shiny {
		pokeImage.URL = pokeInfo.Sprites.ShinyFront
	}

	return pokeImage
}

//Get a pokemon description from a given url. Gets info, parses, and returns.
//Return a string of the description for the mon.
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

//random int between min and max, bounds are INCLUSIVE
func getRandomIntBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

//simple check to see if the entry is in english.
func isEnglish(entry entryInfo) bool {
	return entry.Language.Name == "en"
}

//convert kg to lbs
func kgsToLbs(weightKg float64) float64 {
	return roundTo(weightKg*2.204623, 2)
}

//round n to given decimal
//shamelessly stolen from stack overflow
//	 https://stackoverflow.com/questions/52048218/round-all-decimal-points-in-golang#52048478
func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

//convert height in M to height in ft and inches. Not super accurate, only integer values
func mToFtIn(heightM float64) (int, int) {
	heightInches := int(heightM * 100 / 2.54)   //loss of precision, but small
	return heightInches / 12, heightInches % 12 //more loss of precision, only ints.
}
