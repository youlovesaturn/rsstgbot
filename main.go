package main

import (
	"fmt"
	"log"
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

	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {
		sendPost()
	}

}

func sendPost() {
	post, err := getLatestPost()
	if err != nil {
		return
	}

	m := fmt.Sprintf("[%v](%v)", post.Title, post.Link)

	if post.PubDate != lastPubDate {
		tg.Send(tb.ChatID(cfg.ChannelID), m, &tb.SendOptions{
			ParseMode:             "Markdown",
			DisableWebPagePreview: cfg.DisablePreview,
		})
		lastPubDate = getLastPubDate()
	}
}
