package api

import (
	"errors"
	"os"
	"time"

	"github.com/Atelier-Epita/intra-atelier/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentications
func AuthMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "epitaf",
		Key:        []byte(os.Getenv("JWT_SECRET")),
		Timeout:    time.Hour * 48,
		MaxRefresh: time.Hour * 48,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			u := data.(*models.User)
			// What we put in the JWT claims
			return jwt.MapClaims{
				"email": u.Email,
				"id":    u.Id,
			}
		},
		Authenticator: callbackHandler,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			_ = c.AbortWithError(code, errors.New(message))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	// Authmiddleware must be active
	if err != nil {
		panic(err)
	}

	return authMiddleware
}
