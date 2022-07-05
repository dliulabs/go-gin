package handlers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func (app *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := os.Getenv("X_API_KEY")
		if api_key == "" {
			fmt.Println("curl -X GET http://localhost:5000/dummy -H 'X-API-KEY:1233445'")
			respondWithError(c, 401, "Please set X_API_KEY environment variable")
		}
		token := c.GetHeader("X-API-KEY")
		if token == "" {
			respondWithError(c, 401, "API token required")
		}
		if api_key != token {
			fmt.Printf("Api Key: %v\n", api_key)
			fmt.Printf("Api token: %v\n", token)
			respondWithError(c, 401, "Authentication failed")
		}
		c.Next()
	}
}
