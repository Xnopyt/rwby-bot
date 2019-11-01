package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var data loadedData
var cUUID string

type siteData struct {
	UUID       string `json:"uuid"`
	Title      string `json:"title"`
	EpNum      int    `json:"epnum"`
	MagicShort string `json:"magic_short"`
	MagicLong  string `json:"magic_long"`
}

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
	data.Session.UpdateStatus(0, "initializing...")
	send("Initializing...")
	send("Checking if an episode is avalible...")
	ep := rtGrabLatestEpisodeInfo()
	if ep != nil {
		send("Seeing if we have stored data...")
		if _, err := os.Stat("rwby_info.json"); os.IsNotExist(err) {
			send("I can't find any previous data, running the updater...")
			update()
		} else {
			send("Checking stored data...")
			f, _ := ioutil.ReadFile("rwby_info.json")
			var current siteData
			err := json.Unmarshal(f, &current)
			if err != nil {
				send("Failed to read old data! Running updater...")
				update()
			} else {
				if current.UUID != ep.UUID {
					send("Data is not current, running updater...")
					update()
				} else {
					send("Everything looks good!")
					data.Session.UpdateStatus(0, "Ep "+strconv.Itoa(ep.EpNum)+" - "+ep.Title)
					cUUID = ep.UUID
				}
			}
		}
	} else {
		send("No episodes are currently avalible...")
		data.Session.UpdateStatus(0, "Waiting for Volume 7")
	}
	send("Init Done!")
	for {
		wait()
		check()
	}
}

func update() {
	data.Session.UpdateStatus(0, "Updating...")
	send("Generating a RT account...")
	email, password, err := generateRTAccount()
	if err != nil {
		send("<@!360457422181105666> , Failed to generate a new account. The server sent:")
		send(err.Error())
		send("The script will now panic")
		panic(err)
	}
	send("Authenticating with RT...")
	token := rtAuthenticate(email, password)
	send("Starting a FIRST trial...")
	rtActivateFirst(token)
	send("Grabbing the lastest episode...")
	ep := rtGrabLatestEpisodeInfo()
	magicShort, magicLong := rtGrabLatestEpisode(email, password)
	send("Deleting old site content...")
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
	send("Writting new data...")
	f.Write(JSON)
	f.Close()
	cUUID = ep.UUID
	data.Session.UpdateStatus(0, "Ep "+strconv.Itoa(ep.EpNum)+" - "+ep.Title)
	send("Done!")
}

func wait() {
	fmt.Println("Wait loop begin")
	for {
		t := time.Now()
		if t.Weekday() == 6 && t.Hour() > 9 {
			ep := rtGrabLatestEpisodeInfo()
			if ep != nil {
				if !(ep.GoLive.Day() == t.Day() && t.Month() == ep.GoLive.Month() && ep.GoLive.Year() == t.Year()) {
					break
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func check() {
	fmt.Println("Check loop begin...")
	data.Session.UpdateStatus(0, "Waiting for new episode...")
	for {
		ep := rtGrabLatestEpisodeInfo()
		if ep != nil {
			if ep.UUID != cUUID {
				send("New episode detected!")
				send("The title is: " + ep.Title)
				send("Running updater...")
				update()
				send("New RWBY, " + ep.Title + " , is now avalible at https://how2trianglemuygud.com/rwbyvol6 @everyone")
				break
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func send(msg string) {
	data.Session.ChannelMessageSend("639789890330034208", msg)
	fmt.Println(msg)
}
