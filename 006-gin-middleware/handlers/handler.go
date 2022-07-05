package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Config struct {
	ctx context.Context
}

type Wrapper struct {
	ctx context.Context
}

func NewApp(ctx context.Context) *Config {
	return &Config{
		ctx: ctx,
	}
}

func NewMiddleware(ctx context.Context) *Wrapper {
	return &Wrapper{
		ctx: ctx,
	}
}

// Step 1: basic handler
func (app *Config) GetDummyEndpoint(c *gin.Context) {
	resp := map[string]string{"hello": "world"}
	c.JSON(200, resp)
}

// Step 2: add middleware
func (app *Wrapper) DummyMiddleware() gin.HandlerFunc {
	// Do some run-once initialization logic here Foo()
	return func(c *gin.Context) {
		fmt.Println("Im a dummy!")
		c.Next()
	}
}
