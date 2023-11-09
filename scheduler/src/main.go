package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"scheduler/internal/database"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error reading feed from db", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("unable to mark feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed", err)
		return
	}
	log.Printf("Feed %s collected, %v posts found", rssFeed.Title, len(rssFeed.Items))
	for _, item := range rssFeed.Items {
		// log.Printf("Title: %s ", item.Title)
		db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			PublishedAt: item.PublishedParsed.UTC(),
			FeedID:      feed.ID,
		})
	}
}

func urlToFeed(url string) (gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Fatal("Unable to get feed: ", err)
		return gofeed.Feed{}, nil
	}
	return *feed, nil
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in env")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB", err)
	}

	db := database.New(conn)
	startScraping(db, 4, time.Minute)

}
