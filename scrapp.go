package main

import (
	// "encoding/json"
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// change format of string due to restriction in scrapping proccess
func removeDay(fullDate string) int {
	if strings.HasPrefix(fullDate, "Янва") {
		return January
	} else if strings.HasPrefix(fullDate, "Февр") {
		return February
	} else if strings.HasPrefix(fullDate, "Март") {
		return March
	} else if strings.HasPrefix(fullDate, "Апре") {
		return April
	} else if strings.HasPrefix(fullDate, "Июнь") {
		return June
	} else if strings.HasPrefix(fullDate, "Июль") {
		return July
	} else if strings.HasPrefix(fullDate, "Авгу") {
		return August
	} else if strings.HasPrefix(fullDate, "Сент") {
		return September
	} else if strings.HasPrefix(fullDate, "Октя") {
		return October
	} else if strings.HasPrefix(fullDate, "Нояб") {
		return November
	} else if strings.HasPrefix(fullDate, "Дека") {
		return December
	} else {
		return May
	}
}

// enum for every month
const (
	January   = 1
	February  = 2
	March     = 3
	April     = 4
	May       = 5
	June      = 6
	July      = 7
	August    = 8
	September = 9
	October   = 10
	November  = 11
	December  = 12
)

// main structure of requests
type Play struct {
	Name  string `json:"name"`
	Genre string `json:"genre"`
	Age   string `json:"string"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
}

// slice to store plays collected from web-page
var playsScrapped []Play

func scrapping() []Play {
	file, err := os.Open("plays.json") //opening file (not reading!)
	if err != nil {
		log.Panicln("Failed to create the output JSON file", err)
	}
	defer file.Close()

	c := colly.NewCollector() //main object to work with colly scrapper
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	// c.SetRequestTimeout(40)

	c.OnError(func(e *colly.Response, err error) {
		log.Panicln("Following error occured during scrapping: ", err)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnHTML("div.aff_el", func(e *colly.HTMLElement) {
		month := removeDay(e.ChildText("div.day"))                     //obtaining month of specific play
		currentMonth, err := strconv.Atoi(time.Now().Month().String()) //obtaning current month

		day, err := strconv.Atoi(e.ChildText("div.date")) //obtaining day of specific play
		if err != nil {
			log.Panicln("Error occured during casting from ascii to int")
		}

		//if play is relevant (play is NOT outdated) add to slice of scrapped plays
		if month == 1 && day >= time.Now().Day() {
			playsScrapped = append(playsScrapped, Play{e.ChildText("div.name"), e.ChildText("div.genre"), e.ChildText("div.age_rating"), month, day})
		} else if month >= currentMonth && day >= time.Now().Day() {
			playsScrapped = append(playsScrapped, Play{e.ChildText("div.name"), e.ChildText("div.genre"), e.ChildText("div.age_rating"), month, day})
		}

	})

	c.OnScraped(func(e *colly.Response) {
		// playsRead, _ := os.ReadFile("plays.json")
		// var playsGet []Play

		// err := json.Unmarshal(playsRead, &playsGet)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// fmt.Println(playsGet[len(playsGet)-1])

		// if playsGet[len(playsGet)-1].Month != playsScrapped[len(playsScrapped)-1].Month || playsGet[len(playsGet)-1].Day != playsScrapped[len(playsScrapped)-1].Day {
		// 	fmt.Println("CHANGED")
		// }

		// jsonFormat, _ := json.MarshalIndent(playsScrapped, "", "")
		// file.Write(jsonFormat)

		log.Println("Done scrapping succesfully")
	})

	c.Visit("https://teatrkachalov.ru/affiche/base/")

	return playsScrapped
}
