package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid" // Globally Unique ID Generator
	"golang.org/x/exp/slices"
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

// initialize recipes store
func init() {
	recipes = make([]Recipe, 0)

	data, err := ioutil.ReadFile("recipes.json")
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		log.Fatalf("Unable to unmarshal data: %v", err)
	}

}

// define recipe endpoint handler
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusAccepted, recipe)

}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, &recipes)
}

// using slices.IndexFunc https://cs.opensource.google/go/x/exp/+/master:slices/slices_test.go;l=343;drc=b0d781184e0d33570e9c1a7ea1b0dcdbe5113b78#:~:text=342-,343,-344
func GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	if found := slices.IndexFunc[Recipe](recipes, func(r Recipe) bool { return r.ID == id }); found >= 0 {
		c.JSON(http.StatusOK, recipes[found])
		return
	}
	c.JSON(http.StatusNotFound, id)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	if found := slices.IndexFunc[Recipe](recipes, func(r Recipe) bool { return r.ID == id }); found >= 0 {
		recipe.PublishedAt = time.Now()
		recipe.ID = id
		recipes[found] = recipe
		c.JSON(http.StatusOK, recipe)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	if found := slices.IndexFunc[Recipe](recipes, func(r Recipe) bool { return r.ID == id }); found >= 0 {
		recipe := recipes[found]
		recipes = slices.Delete(recipes, found, found+1)
		c.JSON(http.StatusOK, recipe)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": fmt.Sprintf("ID %v not found", id),
	})

}

func (r Recipe) foundTag(tag string) bool {
	found := slices.IndexFunc[string](r.Tags, func(s string) bool { return strings.EqualFold(s, tag) })
	return found > 0
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for _, recipe := range recipes {
		if recipe.foundTag(tag) {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	engine := gin.New()

	engine.GET("/recipes/search", SearchRecipesHandler)
	engine.POST("/recipes", NewRecipeHandler)
	engine.GET("/recipes", ListRecipesHandler)
	engine.GET("/recipes/:id", GetRecipeHandler)
	engine.PUT("/recipes/:id", UpdateRecipeHandler)
	engine.DELETE("/recipes/:id", DeleteRecipeHandler)
	engine.Run()
}
