package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "time"
	// "github.com/sirupsen/logrus"
)

func toNameMonth(month int) string {
	switch month {
	case 1:
		return "Январь"
	case 2:
		return "Февраль"
	case 3:
		return "Март"
	case 4:
		return "Апрель"
	case 5:
		return "Май"
	case 6:
		return "Июнь"
	case 7:
		return "Июль"
	case 8:
		return "Август"
	case 9:
		return "Сентябрь"
	case 10:
		return "Октябрь"
	case 11:
		return "Ноябрь"
	case 12:
		return "Декабрь"
	default:
		return "Нет"
	}
}

var userID int64

func help(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(userID, "")
	if update.Message.Command() == "help" {
		msg.Text = "\"Check\" - проверить новые постановки\n\"Get\" - все доступные постановки\n"
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

}

func start(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(userID, "")
	if update.Message.Command() == "start" {
		msg.Text = "\"/help\" - все команды"
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	userID = update.Message.Chat.ID
}

func checkUpdates(playsGet *[]Play, bot *tgbotapi.BotAPI) {
	playsScrapped := scrapping()

	log.Println("Последняя постановка (собранная):", playsScrapped[len(playsScrapped)-1])
	log.Println("Последняя постановка (из playsGet):", (*playsGet)[len(*playsGet)-1])

	//check whether dates of last plays from scrapped and plays.json slices are equal
	if (*playsGet)[len(*playsGet)-1].Month != playsScrapped[len(playsScrapped)-1].Month ||
		(*playsGet)[len(*playsGet)-1].Day != playsScrapped[len(playsScrapped)-1].Day {
		file, _ := os.OpenFile("plays.json", os.O_WRONLY|os.O_TRUNC, 0666)

		jsonFormat, _ := json.MarshalIndent(playsScrapped, "", "") //transforming slice to json format object
		file.Write(jsonFormat)
		file.Close()
		// *playsGet = playsScrapped

		playsRead, readerr := os.ReadFile("plays.json")
		if readerr != nil {
			log.Panicln(readerr)
		}

		json.Unmarshal(playsRead, &playsGet)

		log.Println("Последняя постановка (из playsGet) внутри цикла:", (*playsGet)[len(*playsGet)-1])

		msg := tgbotapi.NewMessage(userID, "!!!Появились новые постановки!!!")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}

		log.Println("File plays.json was rewritten")
		return
	}
	log.Println("Nothing has changed in play.json")
	time.Sleep(time.Second * 3)
}

func handlingUpdates(update tgbotapi.Update, bot *tgbotapi.BotAPI, playsGet []Play) {
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text) //getting new message

		//one specific button. may be changed
		if update.Message.Command() == "get" {
			for i := 0; i < len(playsGet); i++ {
				//sending all availible plays. should be changed/deleted
				toSend := strings.ToUpper(playsGet[i].Name) + "\nЖанр: " + playsGet[i].Genre + "\nОграничение: " + 
					playsGet[i].Age + "\nДата: " + toNameMonth(playsGet[i].Month) + " " + 
					strconv.Itoa(playsGet[i].Day) + "\n"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, toSend)
				if _, err := bot.Send(msg); err != nil {
					log.Panicln(err)
				}
			}
		}
}
