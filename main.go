package main

import (
	"fmt"
	"PokeBot/bot"
	"PokeBot/config"
)

func main(){
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()
	<-make(chan struct{})
	return
}
