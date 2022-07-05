package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

var client *mongo.Client
var ctx context.Context

const (
	mongoURL       = "mongodb://localhost:27017"
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

func ParserHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("rss.xml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Printf("Failed to read file: %v\n", err.Error())
		return
	}
	var feed Feed
	err = xml.Unmarshal([]byte(data), &feed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Printf("failed to load data: %v\n", err.Error())
		return
	}
	collection := client.Database(recipeDatabase).Collection("recipes")
	for _, entry := range feed.Entry {
		collection.InsertOne(ctx, bson.M{
			"title":     entry.Title,
			"thumbnail": entry.Thumbnail.URL,
			"url":       entry.Link.Href,
		})
	}

	c.JSON(http.StatusOK, feed.Entry)
}

func init() {
	ctx = context.Background()
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: AuthUserName,
		Password: AuthPassword,
	})
	client, _ = mongo.Connect(ctx, clientOptions)
}

func main() {
	engine := gin.Default()
	engine.POST("/parse", ParserHandler)
	engine.Run()
}
