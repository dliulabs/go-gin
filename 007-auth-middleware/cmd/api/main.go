package main

import (
	"context"
	"log"

	"github.com/go-redis/redis"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"recipes/handlers"

	"github.com/gin-gonic/gin"
)

const (
	mongoURL       = "mongodb://localhost:27017"
	AuthUserName   = "admin"
	AuthPassword   = "password"
	recipeDatabase = "recipes"
)

var client *mongo.Client

// we must create a global variable to access the endpoints handlers.
var app *handlers.RecipesConfig
var authApp *handlers.AuthHandler
var ctx context.Context

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Printf("Redis Ping: %v\n", status)

	app = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

// main is the entry point for the application.
func main() {
	engine := gin.Default()

	// engine.POST("/recipes", app.CreateRecipeHandler)
	engine.GET("/recipes", app.ListRecipesHandler)
	// engine.GET("/recipes/:id", app.GetRecipeHandler)
	// engine.PUT("/recipes/:id", app.UpdateRecipeHandler)
	// engine.DELETE("/recipes/:id", app.DeleteRecipeHandler)
	engine.GET("/recipes/search", app.SearchRecipesHandler)

	authorized := engine.Group("/")
	authorized.Use(authApp.AuthMiddleware())
	{
		authorized.POST("/recipes", app.CreateRecipeHandler)
		authorized.PUT("/recipes/:id", app.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", app.DeleteRecipeHandler)
		authorized.GET("/recipes/:id", app.GetRecipeHandler)
	}

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
