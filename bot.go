package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ses *discordgo.Session

var cUUID string

type siteData struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	EpNum      int    `json:"epnum"`
	MagicShort string `json:"magic_short"`
	MagicLong  string `json:"magic_long"`
}

func main() {
	var err error
	ses, err = discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	ready := make(chan bool)
	ses.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) { ready <- true })
	err = ses.Open()
	if err != nil {
		log.Fatal(err)
	}
	<-ready
	ses.UpdateStatus(0, "initializing...")
	ep := rtGrabLatestEpisodeInfo()
	if ep != nil {
		if _, err := os.Stat("rwby_info.json"); os.IsNotExist(err) {
			update()
		} else {
			f, _ := ioutil.ReadFile("rwby_info.json")
			var current siteData
			err := json.Unmarshal(f, &current)
			if err != nil {
				update()
			} else {
				if current.UUID != ep.UUID {
					update()
				} else {
					ses.UpdateStatus(0, "Ep "+strconv.Itoa(ep.EpNum)+" - "+ep.Title)
					cUUID = ep.UUID
				}
			}
		}
	} else {
		ses.UpdateStatus(0, "Waiting for Volume 8")
		send("No episodes detected, waiting for Volume 8...")
	}
	for {
		wait()
		check()
	}
}

func update() {
	ses.UpdateStatus(0, "Updating...")
	ep := rtGrabLatestEpisodeInfo()
	if ep == nil {
		send("<@!360457422181105666> , Failed to grab episode info.")
		panic(errors.New("Could not pull latest episode info."))
	}
	send("New episode detected!")
	send("The title is: " + ep.Title)
	email, password, err := generateRTAccount()
	if err != nil {
		send("<@!360457422181105666> , Failed to generate a new account. The server sent:")
		send(err.Error())
		send("The script will now panic")
		panic(err)
	}
	token, err := rtAuthenticate(email, password)
	if err != nil {
		send("<@!360457422181105666> , Failed to Authenticate with RT.")
		send(err.Error())
		panic(err)
	}
	err = rtActivateFirst(token)
	if err != nil {
		send("<@!360457422181105666> , Failed to active FIRST trial.")
		send(err.Error())
		panic(err)
	}
	magicShort, magicLong, err := rtGrabLatestEpisode(email, password)
	if err != nil {
		send("<@!360457422181105666> , Failed to grab episode tokens.")
		send(err.Error())
		panic(err)
	}
	if _, err := os.Stat("rwby_info.json"); !os.IsNotExist(err) {
		os.Remove("rwby_info.json")
	}
	f, _ := os.Create("rwby_info.json")
	os.Chmod("rwby_info.json", 0777)
	var store siteData
	store.UUID = ep.UUID
	store.Title = ep.Title
	store.EpNum = ep.EpNum
	store.MagicShort = magicShort
	store.MagicLong = magicLong
	JSON, _ := json.Marshal(store)
	f.Write(JSON)
	f.Close()
	cUUID = ep.UUID
	ses.UpdateStatus(0, "Ep "+strconv.Itoa(ep.EpNum)+" - "+ep.Title)
	send("https://xnopyt.info/rwby?tokenshort=" + magicShort + "&tokenlong=" + magicLong + "&ep=" + strconv.Itoa(ep.EpNum) + "&title=" + url.QueryEscape(ep.Title))
}

func send(msg string) {
	ses.ChannelMessageSend(config.Channel, msg)
	fmt.Println(msg)
}

func wait() {
	fmt.Println("Wait loop begin")
	for {
		t := time.Now()
		if t.Weekday() == 6 && t.Hour() > 9 {
			ep := rtGrabLatestEpisodeInfo()
			if ep == nil {
				break
			}
			if !(ep.GoLive.Day() == t.Day() && t.Month() == ep.GoLive.Month() && ep.GoLive.Year() == t.Year()) {
				break
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func check() {
	fmt.Println("Check loop begin...")
	ses.UpdateStatus(0, "Waiting for new episode...")
	for {
		ep := rtGrabLatestEpisodeInfo()
		if ep != nil {
			if ep.UUID != cUUID {
				update()
				send("@everyone")
				break
			}
		}
		time.Sleep(30 * time.Second)
	}
}
