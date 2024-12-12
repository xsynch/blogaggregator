package main

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	request, err := http.NewRequestWithContext(ctx,http.MethodGet,feedURL,nil)
	if err != nil {
		return nil, err 
	}
	client := http.DefaultClient
	request.Header.Add("User-Agent","gator")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err 
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err 
	}
	newRSSFeed := RSSFeed{}
	err = xml.Unmarshal(body,&newRSSFeed)
	if err != nil {
		return nil, err 
	}
	return &newRSSFeed,nil 
}
