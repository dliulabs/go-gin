package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"recipes/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Config struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewApp(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *Config {
	return &Config{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (app *Config) CreateRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := app.collection.InsertOne(app.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error while inserting a new recipe %v", err.Error())})
		return
	}
	data, _ := json.Marshal(recipe)
	app.redisClient.Set(fmt.Sprintf("recipes:%v", recipe.ID.Hex()), string(data), 1*time.Hour)
	c.JSON(http.StatusOK, recipe)
}

func (app *Config) ListRecipesHandler(c *gin.Context) {
	val, err := app.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Printf("Response documents from MongoDB")
		cursor, err := app.collection.Find(app.ctx, bson.M{}) // M for a map
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			defer cursor.Close(app.ctx)
			return
		}

		recipes := make([]models.Recipe, 0)
		if err = cursor.All(app.ctx, &recipes); err != nil {
			panic(err)
		}
		/*
			for cursor.Next(app.ctx) {
				var recipe models.Recipe
				cursor.Decode(&recipe)
				recipes = append(recipes, recipe)
			}
		*/
		data, _ := json.Marshal(recipes)
		app.redisClient.Set("recipes", string(data), 1*time.Hour)

		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		log.Printf("Response documents from Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}
}

func (app *Config) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	recipe.ID = objectId
	recipe.PublishedAt = time.Now()
	result, err := app.collection.UpdateOne(
		app.ctx,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{
				{"name", recipe.Name},
				{"instructions", recipe.Instructions},
				{"ingredients", recipe.Ingredients},
				{"tags", recipe.Tags},
				{"publishedAt", recipe.PublishedAt},
			}},
		},
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(recipe)
	app.redisClient.Set(fmt.Sprintf("recipes:%v", recipe.ID.Hex()), string(data), 1*time.Hour)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated %v Documents!", result.ModifiedCount)})
}

func (app *Config) GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	val, err := app.redisClient.Get(fmt.Sprintf("recipes:%v", objectId.Hex())).Result()
	var recipe models.Recipe // returns "id"
	if err == redis.Nil {
		log.Printf("Response documents from MongoDB")
		cursor := app.collection.FindOne(
			app.ctx,
			bson.M{"_id": objectId},
		)

		// var recipe bson.M // returns "_id"
		err := cursor.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recipe)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		log.Printf("Response documents from Redis")
		json.Unmarshal([]byte(val), &recipe)
		c.JSON(http.StatusOK, recipe)
	}
}

func (app *Config) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, err := app.collection.DeleteOne(app.ctx, bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	app.redisClient.Del(fmt.Sprintf("recipes:%v", objectId.Hex()))
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted %v document(s)", result.DeletedCount)})
}

func (app *Config) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	val, err := app.redisClient.Get(fmt.Sprintf("recipes:tag:%v", tag)).Result()
	recipes := make([]models.Recipe, 0)
	if err == redis.Nil {
		log.Printf("Response documents from MongoDB")
		cursor, err := app.collection.Find(app.ctx, bson.M{"tags": tag}) // M for a map
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			defer cursor.Close(app.ctx)
			return
		}

		if err = cursor.All(app.ctx, &recipes); err != nil {
			panic(err)
		}
		data, _ := json.Marshal(recipes)
		app.redisClient.Set(fmt.Sprintf("recipes:tag:%v", tag), string(data), 1*time.Hour)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		log.Printf("Response documents from Redis")
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}
}
