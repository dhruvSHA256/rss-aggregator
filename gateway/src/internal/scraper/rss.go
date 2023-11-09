package scraper

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func urlToFeed(url string) (gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Fatal("Unable to get feed: ", err)
		return gofeed.Feed{}, nil
	}
	return *feed, nil
}
