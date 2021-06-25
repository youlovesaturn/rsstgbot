package main

import (
	"encoding/xml"
	"errors"
	"net/http"
)

type Item struct {
	Link    string `xml:"link"`
	Title   string `xml:"title"`
	PubDate string `xml:"pubDate"`
}

type Channel struct {
	Link  string `xml:"link"`
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

func getLatestPost() (post *Item, err error) {
	var rss RSS

	resp, err := http.Get(cfg.FeedURL)
	if err != nil {
		return
	}

	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&rss)
	if err != nil {
		return
	}

	if len(rss.Channel.Items) == 0 {
		return nil, errors.New("no result")
	}

	return &rss.Channel.Items[0], nil
}

func getLastPubDate() string {
	post, _ := getLatestPost()
	return post.PubDate
}
