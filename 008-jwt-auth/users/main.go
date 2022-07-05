package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoURL          = "mongodb://localhost:27017"
	AuthUserName      = "admin"
	AuthPassword      = "password"
	recipeDatabase    = "recipes"
	redisHost         = "localhost:6379"
	recipesCollection = "recipes"
	usersCollection   = "users"
)

func main() {
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"mlabouardy": "L3nSFRcZzNQ67bcc",
	}

	ctx := context.Background()
	client, err := connectToMongo(ctx)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collectionUsers := client.Database(recipeDatabase).Collection(usersCollection)
	h := sha256.New()

	for username, password := range users {
		// https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
		fmt.Printf("Inssert user: %v\n", username)
		sha := base64.URLEncoding.EncodeToString(h.Sum([]byte(password)))
		fmt.Printf("Inssert password: %v\n", sha)
		result, err := collectionUsers.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(sha),
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		} else {
			fmt.Printf("Added %v document(s)\n", result)
		}
	}
}

func connectToMongo(ctx context.Context) (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: AuthUserName,
		Password: AuthPassword,
	})

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
