package main

import "github.com/bwmarrin/discordgo"

func main() {
	data := loadData()
	discord, err := discordgo.New("Bot " + data.botToken)
	if err != nil {
		panic(err)
	}
	ready := make(chan bool)
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	err = discord.Open()
	if err != nil {
		panic(err)
	}
	<-ready
}
