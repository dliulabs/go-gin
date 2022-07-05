package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

var mongoClient *mongo.Client
var ctx context.Context
var channelAmqp *amqp.Channel
var queue amqp.Queue

const (
	mongoURL       = "mongodb://mongo:27017" // "mongodb://localhost:27017"
	AuthUserName   = "admin"
	AuthPassword   = "password"
	recipeDatabase = "recipes"
)

type Request struct {
	URL string `json:"url"`
}

type Feed struct {
	Entry []Entry `xml:"entry"`
}

type Entry struct {
	Link struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Title     string `xml:"title"`
	Thumbnail struct {
		Text string `xml:",chardata"`
		URL  string `xml:"url,attr"`
	} `xml:"thumbnail"`
}

func init() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	channelAmqp, _ = amqpConnection.Channel()

	queue, _ = channelAmqp.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"), // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)

	ctx = context.Background()
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: AuthUserName,
		Password: AuthPassword,
	})
	mongoClient, _ = mongo.Connect(ctx, clientOptions)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer amqpConnection.Close()
	channelAmqp, err := amqpConnection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channelAmqp.Close()

	queue, err := channelAmqp.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"), // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channelAmqp.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	forever := make(chan bool)
	msgs, err := channelAmqp.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack: false for manual ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var request Request
			json.Unmarshal(d.Body, &request)
			log.Println("RSS URL:", request.URL)
			entries, _ := GetFeedEntries(request.URL)
			collection := mongoClient.Database(recipeDatabase).Collection("recipes")
			for _, entry := range entries {
				collection.InsertOne(ctx, bson.M{
					"title":     entry.Title,
					"thumbnail": entry.Thumbnail.URL,
					"url":       entry.Link.Href,
				})
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func GetFeedEntries(url string) ([]Entry, error) {
	data, err := ioutil.ReadFile("rss.xml")
	var feed Feed
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err.Error())
		return nil, err
	}
	err = xml.Unmarshal([]byte(data), &feed)
	if err != nil {
		fmt.Printf("failed to load data: %v\n", err.Error())
		return nil, err
	}

	return feed.Entry, nil

}
