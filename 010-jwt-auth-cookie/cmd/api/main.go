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

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const (
	mongoURL          = "mongodb://localhost:27017"
	AuthUserName      = "admin"
	AuthPassword      = "password"
	recipeDatabase    = "recipes"
	recipesCollection = "recipes"
	usersCollection   = "users"
	redisHost         = "localhost:6379"
	redisPassword     = "password"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key      = []byte("super-secret-key")
	store, _ = redisStore.NewStore(10, "tcp", redisHost, redisPassword, key)
)

var client *mongo.Client

// we must create a global variable to access the endpoints handlers.
var app *handlers.RecipesConfig
var authApp *handlers.AuthHandler
var ctx context.Context

func init() {
	// create a context in order to disconnect
	ctx = context.Background()
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
	collection := client.Database(recipeDatabase).Collection(recipesCollection)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	status := redisClient.Ping()
	log.Printf("Redis Ping: %v\n", status)

	app = handlers.NewRecipesHandler(ctx, collection, redisClient)
	collectionUsers := client.Database(recipeDatabase).Collection(usersCollection)
	authApp = handlers.NewAuthHandler(ctx, collectionUsers)
}

// main is the entry point for the application.
func main() {

	engine := gin.Default()
	engine.Use(sessions.Sessions("recipes_api", store))

	engine.POST("/signin", authApp.SignInHandler)
	engine.POST("/signout", authApp.SignOutHandler)
	engine.POST("/refresh", authApp.RefreshHandler)

	authorized := engine.Group("/")
	authorized.Use(authApp.AuthMiddleware())
	{
		engine.GET("/recipes", app.ListRecipesHandler)
	}

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
