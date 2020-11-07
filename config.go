package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type configLayout struct {
	BotToken         string `json:"bot_token"`
	Channel          string `json:"channel"`
	AnticaptchaToken string `json:"anticaptcha_token"`
	CardInfo         struct {
		FName string `json:"fname"`
		LName string `json:"lname"`
		Num   string `json:"num"`
		Mon   string `json:"mon"`
		Yea   string `json:"yea"`
		CVV   string `json:"cvv"`
		PCode string `json:"pcode"`
	} `json:"card_info"`
}

var config *configLayout

func init() {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	config = new(configLayout)
	err = json.Unmarshal(configFile, config)
	if err != nil {
		log.Fatal(err)
	}
}
