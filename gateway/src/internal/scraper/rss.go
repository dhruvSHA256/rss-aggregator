package scraper

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func urlToFeed(url string) (gofeed.Feed, error) {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	if err != nil {
		log.Fatal("Unable to get feed: ", err)
		return gofeed.Feed{}, nil
	}
	return *feed, nil
}
