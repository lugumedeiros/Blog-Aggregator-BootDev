package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"html"
	db "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/database"
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

func FetchFeed(feedUrl string) (*RSSFeed, error){
	req, err := http.NewRequestWithContext(db.AppDB.Ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err_resp := client.Do(req)
	if err_resp != nil {
		return nil, err_resp
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ERR: %v", resp.Status)
	}
	body, errr_body := io.ReadAll(resp.Body)
	if errr_body != nil {
		return nil, errr_body
	}
	defer resp.Body.Close()

	var feed RSSFeed
	err_unmarshal := xml.Unmarshal(body, &feed)
	if err_unmarshal != nil {
		return nil, err_unmarshal
	}
	
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		item := &feed.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &feed, nil
}