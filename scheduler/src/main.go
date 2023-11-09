package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scheduler/internal/database"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/mmcdole/gofeed"
	"github.com/streadway/amqp"
)

func sendMessage(message string, conn *amqp.Connection, queueName string) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()
	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	fmt.Printf("Message sent to RabbitMQ: %s\n", message)
	return nil
}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration, rabbit_conn *amqp.Connection, queueName string) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error reading feed from db", err)
			continue
		}
		for _, feed := range feeds {
			data := map[string]interface{}{
				"feed_id": feed.ID,
			}
			message_bytes, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			message := string(message_bytes)
			err = sendMessage(message, rabbit_conn, queueName)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
			if err != nil {
				log.Println("unable to mark feed as fetched", err)
				return
			}
		}
	}
}

// func scrapeFeed(db *database.Queries, feed database.Feed) {

// 	rssFeed, err := urlToFeed(feed.Url)
// 	if err != nil {
// 		log.Println("error fetching feed", err)
// 		return
// 	}
// 	log.Printf("Feed %s collected, %v posts found", rssFeed.Title, len(rssFeed.Items))
// 	for _, item := range rssFeed.Items {
// 		db.CreatePost(context.Background(), database.CreatePostParams{
// 			ID:          uuid.New(),
// 			CreatedAt:   time.Now().UTC(),
// 			UpdatedAt:   time.Now().UTC(),
// 			Title:       item.Title,
// 			Url:         item.Link,
// 			PublishedAt: item.PublishedParsed.UTC(),
// 			FeedID:      feed.ID,
// 		})
// 	}
// }

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
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		log.Fatal("RABBIT_URL not set in env")
	}
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		log.Fatal("QUEUE_NAME not set in env")
	}

	// Connect to RabbitMQ server
	rabbit_conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ: %v", err)
		return
	}
	db_conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB", err)
	}

	defer func() {
		rabbit_conn.Close()
		db_conn.Close()
	}()

	db := database.New(db_conn)
	startScraping(db, 4, time.Minute, rabbit_conn, queueName)
}
