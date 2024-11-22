package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jasonlvhit/gocron"
	"log"
	"os"
)

var playsGet []Play

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmsgprefix) //logger setting

	bot, err := tgbotapi.NewBotAPI("6115886452:AAGgRk9mEHlI7hJsF4LU_fUIMbV7ZygDowc") // bot setting
	if err != nil {
		log.Panicln(err) //handling possible error
	}
	bot.Debug = true //set outputing of all technical information related to bot-api

	updateConfig := tgbotapi.NewUpdate(0) //to know number of values handled
	updateConfig.Timeout = 10

	playsRead, readerr := os.ReadFile("plays.json")
	if readerr != nil {
		log.Panicln(readerr)
	}

	errjson := json.Unmarshal(playsRead, &playsGet)
	log.Println("Последняя постановка (из playsGet) после unmrshal:", playsGet[len(
		playsGet)-1])

	if errjson != nil {
		log.Panicln("Failed to unmarshal plays.json: ", errjson) //handling error during getting data from json file
	}

	//In case file plays.json is empty
	if len(playsGet) == 0 {
		file, _ := os.Create("plays.json")
		defer file.Close()
		jsonFormat, _ := json.MarshalIndent(scrapping(), "", "")
		file.Write(jsonFormat)
		playsRead, _ = os.ReadFile("plays.json")
		fmt.Println("Was")
	}

	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Checking every 30 minutes changes on the web-page
	go func() {
		for {
			gocron.Every(10).Seconds().Do(checkUpdates, &playsGet, bot)
			<-gocron.Start()
			log.Println("HERE?")
		}
	}()

	//Handling updates from telegram
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		// if update.Message == nil {
		// 	continue // ignore non-Message updates
		// }
		handlingUpdates(update, bot, playsGet)
		help(update, bot)
		start(update, bot)
	}
}
