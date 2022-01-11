package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	config *configStruct
	Token, BotPrefix string
)

type configStruct struct {
	Token string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
}

func ReadConfig() error {

	//Read the configuration file, token and prefix. Handle the err
	file, err := ioutil.ReadFile("./config.json")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	//unmarshall the file into the main struct. handle err.
	err = json.Unmarshal(file, &config)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	//unpack struct
	Token = config.Token
	BotPrefix = config.BotPrefix

	//bot info set up successfully at this point, return nothing if success.
	return nil
}
