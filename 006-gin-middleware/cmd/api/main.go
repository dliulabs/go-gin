package main

import (
	"context"
	"fmt"

	"middleware/handlers"

	"github.com/gin-gonic/gin"
)

var app *handlers.Config
var wrapper *handlers.Wrapper
var keyApp *handlers.Apikey

func init() {
	ctx := context.Background()
	app = handlers.NewApp(ctx)
	wrapper = handlers.NewMiddleware(ctx)
	keyApp = handlers.NewApikey(ctx)
}

func main() {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("key", "foo")
	})

	engine.Use(func(c *gin.Context) {
		fmt.Println(c.MustGet("key").(string)) // foo
	})
	/* Step 1
	// engine.GET("/dummy", app.GetDummyEndpoint)
	fmt.Println("curl -X GET http://localhost:5000/dummy")
	*/
	/* Step 2 */
	authorized := engine.Group("/")
	authorized.Use(wrapper.DummyMiddleware())

	/* Step 3 */
	authorized.Use(keyApp.AuthMiddleware())
	{
		authorized.GET("/dummy", app.GetDummyEndpoint)
	}
	// common steps
	engine.Run(":5000")
}
