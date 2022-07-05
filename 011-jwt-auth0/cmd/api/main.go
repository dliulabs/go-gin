package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"recipes/handlers"
	"recipes/models"

	"github.com/gin-gonic/gin"
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

var client *mongo.Client
var redisClient *redis.Client

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

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "password",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Printf("Redis Ping: %v\n", status)

	// loadRecipes(ctx) // run once when initializing

	app = handlers.NewRecipesHandler(ctx, collection, redisClient)
	collectionUsers := client.Database(recipeDatabase).Collection(usersCollection)
	authApp = handlers.NewAuthHandler(ctx, collectionUsers)
}

// main is the entry point for the application.
func main() {
	engine := gin.Default()

	engine.GET("/recipes", app.ListRecipesHandler)

	authorized := engine.Group("/")
	authorized.Use(authApp.AuthMiddleware())
	{
		authorized.POST("/recipes", app.CreateRecipeHandler)
		authorized.PUT("/recipes/:id", app.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", app.DeleteRecipeHandler)
		authorized.GET("/recipes/:id", app.GetRecipeHandler)
	}

	// engine.Run()
	engine.RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")
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

type JsonRecipe struct {
	ID           string    `json:"id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}

func loadRecipes(ctx context.Context) {
	var recipes = make([]JsonRecipe, 0)

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
		mongoRecipe := models.Recipe{
			ID:           primitive.NewObjectID(),
			Name:         recipe.Name,
			Tags:         recipe.Tags,
			Ingredients:  recipe.Ingredients,
			Instructions: recipe.Instructions,
			PublishedAt:  time.Now(),
		}
		listOfRecipes = append(listOfRecipes, mongoRecipe)
		obj, _ := json.Marshal(mongoRecipe)
		redisClient.Set(fmt.Sprintf("recipes:%v", mongoRecipe.ID.Hex()), string(obj), 1*time.Hour)
	}
	collection := client.Database(recipeDatabase).Collection("recipes")
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
	data, _ = json.Marshal(listOfRecipes)
	redisClient.Set("recipes", string(data), 1*time.Hour)
}
