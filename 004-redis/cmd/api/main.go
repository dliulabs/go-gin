package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"recipes/handlers"

	"github.com/gin-gonic/gin"
)

// define recipe struct
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// define recipes store
var recipes []Recipe

const (
	mongoURL       = "mongodb://localhost:27017"
	AuthUserName   = "admin"
	AuthPassword   = "password"
	recipeDatabase = "recipes"
)

var client *mongo.Client
var redisClient *redis.Client

// we must create a global variable to access the endpoints handlers.
var app *handlers.Config
var ctx context.Context

func loadRecipes(ctx context.Context) {
	recipes = make([]Recipe, 0)

	data, err := ioutil.ReadFile("recipes.json")
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		log.Fatalf("Unable to unmarshal data: %v", err)
	}
	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	collection := client.Database(recipeDatabase).Collection("recipes")
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ",
		len(insertManyResult.InsertedIDs))
}

func init() {
	// create a context in order to disconnect
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()

	// connect to mongo
	mongoClient, err := connectToMongo(ctx)
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	collection := client.Database(recipeDatabase).Collection("recipes")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Printf("Redis Ping: %v\n", status)

	app = handlers.NewApp(ctx, collection, redisClient)
}

// main is the entry point for the application.
func main() {
	engine := gin.New()

	engine.POST("/recipes", app.CreateRecipeHandler)
	engine.GET("/recipes", app.ListRecipesHandler)
	engine.GET("/recipes/:id", app.GetRecipeHandler)
	engine.PUT("/recipes/:id", app.UpdateRecipeHandler)
	engine.DELETE("/recipes/:id", app.DeleteRecipeHandler)
	engine.GET("/recipes/search", app.SearchRecipesHandler)

	// close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	engine.Run()
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
