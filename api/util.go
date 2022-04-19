package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
)

func Abort(c *gin.Context, err error, status int, message string) {
	if err != nil {
		zap.S().Error(err)
	}
	c.AbortWithStatusJSON(status, gin.H{
		"status":  status,
		"message": message,
	})
}

var SC *securecookie.SecureCookie

func initSecureCookies() {
	SC = securecookie.New([]byte(os.Getenv("COOKIE_SECRET")), nil)
}

func GenerateCookie(id string) string {
	encoded, err := SC.Encode(os.Getenv("COOKIE_NAME"), id)
	if err != nil {
		return ""
	}
	return encoded
}
