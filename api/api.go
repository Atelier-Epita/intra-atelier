package api

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	// GoSwagger
	_ "github.com/Atelier-Epita/intra-atelier/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

const version = "v0.1"

var router *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

func CreateRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(cors())

	// Swagger
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})

	router = r.Group("/" + version)
	// JWT
	// auth = AuthMiddleware()

	// handleAuth()

	// router.Use(auth.MiddlewareFunc())

	handleUser()
	handleGroup()
	handleFile()
	handleEvent()

	return r
}
