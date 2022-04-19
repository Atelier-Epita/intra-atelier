package api

import (
	"os"

	"github.com/Atelier-Epita/intra-atelier/models"
	"github.com/gin-gonic/gin"
)

// SessionMiddleware handles cookie
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := os.Getenv("COOKIE_NAME")
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			return
		}

		var mail string
		err = SC.Decode(cookieName, cookie, &mail)
		if err != nil {
			return
		}

		user, err := models.GetUserByMail(mail)
		if err != nil {
			return
		}

		c.Set("user", user.Email)
	}
}
