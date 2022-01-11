package main

import (
	"fmt"
	"PokeBot/bot"
	"PokeBot/config"
)

func main(){
	//call into reader code and handle error
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//start up the bot
	bot.Start()
	<-make(chan struct{})
	return
}
