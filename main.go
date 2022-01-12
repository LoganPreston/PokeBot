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

	var (
		goBot *discordgo.Session
		err   error
	)

	//call into reader code and handle error
	if err = config.ReadConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	//start up the bot
	//init the bot
	if goBot, err = discordgo.New("Bot " + config.Token); err != nil {
		fmt.Println(err.Error())
		return
	}

	//setup bot and add a message handler to listen on
	bot.BotSetup(goBot)
	goBot.AddHandler(bot.MessageHandler)

	//open the bot. Defer the close so we don't forget to close gracefully
	if err = goBot.Open(); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer goBot.Close()

	fmt.Println("Bot is live!")

	//wait until we get a ctrl+C or some other interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//graceful cleanup, deferred close is run before return.
	fmt.Println("\nBot is shutting down...")
	return
}
