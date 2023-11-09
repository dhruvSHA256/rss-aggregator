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

func sendMessage(message string, queue *amqp.Queue, ch *amqp.Channel) error {
	log.Printf("sending message %v to queue %q\n", message, queue.Name)
	err := ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Println("Failed to publish a message")
		return err
	}
	return nil
}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration, queue *amqp.Queue, ch *amqp.Channel) {
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
			err = sendMessage(message, queue, ch)
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
	db_conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB", err)
	}

	// Connect to RabbitMQ server
	rabbit_conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ: %v", err)
		return
	}
	ch, err := rabbit_conn.Channel()
	if err != nil {
		fmt.Print("failed to open channel")
		return
	}
	defer func() {
		ch.Close()
		rabbit_conn.Close()
		db_conn.Close()
	}()
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Print("Failed to declare a queue")
		return
	}

	db := database.New(db_conn)
	startScraping(db, 4, time.Minute, &q, ch)
}
