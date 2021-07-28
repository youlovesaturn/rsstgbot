package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	tb "gopkg.in/tucnak/telebot.v2"
)

var cfg config
var tg *tb.Bot
var lastPubDate string

func main() {
	var err error

	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	tg, err = tb.NewBot(tb.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	if fileExists(cfg.Filename) {
		lastPubDate, err = readFromFile(cfg.Filename)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err := os.Create(cfg.Filename)
		if err != nil {
			log.Println(err)
			return
		}
	}

	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		sendPost()
	}

}

func sendPost() {
	post, err := getLatestPost()
	if err != nil {
		log.Println(err)
		return
	}

	m := fmt.Sprintf("[%v](%v)", post.Title, post.Link)

	if post.PubDate != lastPubDate {
		tg.Send(tb.ChatID(cfg.ChannelID), m, &tb.SendOptions{
			ParseMode:             "Markdown",
			DisableWebPagePreview: cfg.DisablePreview,
		})
		lastPubDate = getLastPubDate()
		writeToFile(lastPubDate)

	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readFromFile(path string) (str string, err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	l, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return
	}
	return string(l), nil
}

func writeToFile(str string) (err error) {
	f, err := os.OpenFile(cfg.Filename, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
