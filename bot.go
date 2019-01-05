package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var data loadedData

func main() {
	data = loadData()
	var err error
	data.Session, err = discordgo.New("Bot " + data.botToken)
	if err != nil {
		panic(err)
	}
	ready := make(chan bool)
	data.Session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	err = data.Session.Open()
	if err != nil {
		panic(err)
	}
	defer data.Session.Close()
	<-ready
}

func send(msg string) {
	data.Session.ChannelMessageSend("445190274902261770", msg)
	fmt.Println(msg)
}
