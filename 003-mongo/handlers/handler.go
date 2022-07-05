package handlers

import (
	"fmt"
	"net/http"
	"time"

	"recipes/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Config struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewApp(ctx context.Context, collection *mongo.Collection) *Config {
	return &Config{
		collection: collection,
		ctx:        ctx,
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
	c.JSON(http.StatusOK, recipe)
}

func (app *Config) ListRecipesHandler(c *gin.Context) {
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
	c.JSON(http.StatusOK, recipes)
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
	result, err := app.collection.UpdateOne(
		app.ctx,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{
				{"name", recipe.Name},
				{"instructions", recipe.Instructions},
				{"ingredients", recipe.Ingredients},
				{"tags", recipe.Tags},
				{"publishedAt", time.Now()},
			}},
		},
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated %v Documents!", result.ModifiedCount)})
}

func (app *Config) GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cursor := app.collection.FindOne(
		app.ctx,
		bson.M{"_id": objectId},
	)

	var recipe models.Recipe // returns "id"
	// var recipe bson.M // returns "_id"
	err := cursor.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (app *Config) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, err := app.collection.DeleteOne(app.ctx, bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted %v document(s)", result.DeletedCount)})
}

func (app *Config) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	cursor, err := app.collection.Find(app.ctx, bson.M{"tags": tag}) // M for a map
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		defer cursor.Close(app.ctx)
		return
	}

	listOfRecipes := make([]models.Recipe, 0)
	if err = cursor.All(app.ctx, &listOfRecipes); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, listOfRecipes)
}
