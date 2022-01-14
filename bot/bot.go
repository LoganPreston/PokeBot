package bot

import (
	"fmt"

	"PokeBot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string

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

	if m.Content == "" {
		return
	}

	//split message into the first char prefix and the rest
	prefix, content := m.Content[:1], m.Content[1:]

	//make sure first character is the bot's flag
	if prefix != config.BotPrefix {
		return
	}

	var (
		reply interface{}
		err   error
	)

	//switch based on the content of the message/request. Can cause bot to do nothing
	switch content {
	//handle pokemon messages
	case "pokemon":
		reply, err = replyToPokemonMessage()
	case "triviaQuestion":
		reply = replyToTriviaQuestionMessage(m.ChannelID)
	case "triviaAnswer":
		reply = replyToTriviaAnswerMessage(m.ChannelID)
	case "commands":
		reply = fmt.Sprintf("I support: \n\t!pokemon\n\t!triviaQuestion\n\t!triviaAnswer\n\t!info")
	case "info":
		reply = "I was created by Logan Preston to practice Go. I don't do much outside of Pokemon..."
	//do nothing, just leave
	default:
		return
	}

	//if we broke somewhere, politely tell user sorry, but inform dev of break
	if err != nil {
		reply = "I'm sorry, I failed somewhere along the way. Try again"
		fmt.Println(err.Error())
	}

	//reply with a switch, handle embedded and simple messages
	switch message := reply.(type) {
	case *discordgo.MessageEmbed:
		s.ChannelMessageSendEmbed(m.ChannelID, message)
	case string:
		s.ChannelMessageSend(m.ChannelID, message)
	//log an error, should never be hit
	default:
		fmt.Printf("Unknown message type %T received\n", message)
	}

	return
}
