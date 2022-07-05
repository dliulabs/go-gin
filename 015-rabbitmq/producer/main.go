package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Request struct {
	URL string `json:"url"`
}

var channelAmqp *amqp.Channel
var queue amqp.Queue

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
}

func ParserHandler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	data, _ := json.Marshal(request)
	err := channelAmqp.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while publishing to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"message": "success"})
}

func main() {
	engine := gin.Default()
	engine.POST("/parse", ParserHandler)
	engine.Run(":5000")
}
