package api

import (
	"github.com/gin-gonic/gin"
)

var router *gin.RouterGroup

func CreateRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(cors())
	r.GET()


	return nil
}