package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"worker/internal/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/mmcdole/gofeed"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	FeedId uuid.UUID `json:"feed_id"`
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

func scrapeFeed(db *database.Queries, feed database.Feed) {
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed", err)
		return
	}
	log.Printf("Feed %s collected, %v posts found", rssFeed.Title, len(rssFeed.Items))
	for _, item := range rssFeed.Items {
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

func main() {
	log.Println("Starting worker")
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
	db_conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB", err)
	}
	db := database.New(db_conn)

	// Connect to RabbitMQ server
	rabbit_conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ: %v", err)
		return
	}
	defer func() {
		rabbit_conn.Close()
		db_conn.Close()
	}()
	ch, err := rabbit_conn.Channel()
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	log.Printf("Created queue: %v", q.Name)
	if err != nil {
		fmt.Printf("Failed to declare queue %v", err)
		return
	}
	msgs, err := ch.Consume(
		queueName, // Queue name
		"",        // Consumer
		false,     // Auto-Ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("Unable to consume from %v, err: %v", queueName, err)
	}
	var forever chan struct{}
	go func() {
		for delivery := range msgs {
			var msg Message
			err := json.Unmarshal([]byte(delivery.Body), &msg)
			if err != nil {
				log.Fatal("Failed to unmarshal JSON:", err)
			}
			fmt.Printf("Received feed_id: %s\n", msg.FeedId)
			feed, err := db.GetFeedById(context.Background(), msg.FeedId)
			if err != nil {
				log.Println("Unable to get feed: ", err)
				log.Println("feed: ", feed)
				delivery.Nack(false, true)
				continue
			}
			scrapeFeed(db, feed)
			delivery.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
