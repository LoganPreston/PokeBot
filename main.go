package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"PokeBot/bot"
	"PokeBot/config"

	"github.com/bwmarrin/discordgo"
)

func main() {

	//env vars

	var goBot *discordgo.Session
	//call into reader code and handle error
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//start up the bot
	//init the bot
	goBot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//setup bot and add a message handler to listen on
	bot.BotSetup(goBot)
	goBot.AddHandler(bot.MessageHandler)

	//open the bot
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is live!")

	//wait until we get a ctrl+C or some other interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//graceful cleanup
	fmt.Println("Bot is shutting down...")
	goBot.Close()
	return
}
