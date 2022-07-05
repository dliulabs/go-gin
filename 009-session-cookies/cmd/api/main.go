package main

import (
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const (
	redisHost     = "localhost:6379"
	redisPassword = "password"
	redisDatabase = 10
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key      = []byte("super-secret-key")
	store, _ = redisStore.NewStore(redisDatabase, "tcp", redisHost, redisPassword, key)
)

func main() {
	engine := gin.Default()
	engine.Use(sessions.Sessions("recipes_api", store))

	engine.GET("/login", func(c *gin.Context) {
		sessionToken := xid.New().String()
		session := sessions.Default(c)
		session.Options(sessions.Options{
			MaxAge: 3600 * 1, // 1hrs
		})
		session.Set("token", sessionToken)

		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})

	engine.Run(":8000")

}
