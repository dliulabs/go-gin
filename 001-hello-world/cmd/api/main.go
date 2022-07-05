package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

var PORT = 3000

type Contact struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"required,e164"`
	CountryCode string `json:"countryCode" binding:"required,iso3166_1_alpha2"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func main() {
	log.Print("This is our first log message in Go.")
	engine := gin.New()
	engine.SetTrustedProxies([]string{"127.0.0.1"})

	engine.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	engine.GET("/:name", func(context *gin.Context) {
		name := context.Params.ByName("name")

		context.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("hello %v", name),
		})
	})

	engine.POST("/contact", func(context *gin.Context) {
		body := Contact{}
		if err := context.ShouldBindJSON(&body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				out := make([]ErrorMsg, len(ve))
				for i, fe := range ve {
					out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
				}
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			}
			return
		}
		fmt.Println(body)
		context.JSON(http.StatusAccepted, &body)
	})

	log.Printf("try me: curl -X GET http://localhost:%d\n", PORT)
	log.Printf("try me: curl -X GET http://localhost:%d/%s\n", PORT, "david")
	log.Printf(`try me:
	curl -X GET curl -X POST http://localhost:%v/contact -d '{"firstname":"David","lastname":"Liu","email":"dliu@example.com","phone":"+16155551212","countrycode":"US"}`, PORT)
	engine.Run(fmt.Sprintf(":%d", PORT))
}
