package api

import (
	"github.com/gin-gonic/gin"
)

var router *gin.RouterGroup

func CreateRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router = r.Group("/")
	handleUser()
	handleGroup()

	return r
}
