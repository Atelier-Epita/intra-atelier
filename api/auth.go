package api

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/epitaf/lib/microsoft"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleAuth() {
	users := router.Group("/users")

	users.GET("/authenticate", authenticateHandler)
	users.GET("/callback", auth.LoginHandler)
}

// @Summary Authenticate URL
// @Tags auth
// @Description Build Microsoft oauth url
// @Param   redirect_uri	body	string	true	"redirect_uri"  default(http://localhost:8080/callback)
// @Success 200	"OK"
// @Failure 406	"Not acceptable"
// @Router /users/authenticate [POST]
func authenticateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, microsoft.SignInURL("http://localhost:8080/v0.1/users/callback"))
}

// @Summary OAuth Callback
// @Description Authenticate user and return JWT
// @Tags auth
// @Param   code	body	string	true	"code"
// @Param   redirect_uri	body	string	true	"redirect_uri"
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /users/callback [GET]
func callbackHandler(c *gin.Context) (interface{}, error) {
	var code = c.Request.URL.Query().Get("code")
	var redirect_uri = "http://localhost:8080/v0.1/users/callback"
	fmt.Printf(code, redirect_uri)

	token, err := microsoft.GetAccessToken(code, redirect_uri)
	if err != nil {
		return nil, err
	}

	client := microsoft.NewClient(token, nil)
	// Retrieve microsoft profile
	profile, err := client.GetProfile()
	if err != nil {
		zap.S().Error(err)
		return nil, jwt.ErrFailedAuthentication
	}

	// Check if user exists in database
	u, err := models.GetUserByEmail(profile.Mail)
	if err != nil {
		// If the user does not exists, we must create a new one using the CRI.
		user, err := models.PrepareUser(profile.Mail)
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}

		// Insert new user and return user data
		err = user.Insert()
		if err != nil {
			zap.S().Error(err)
			return nil, jwt.ErrFailedAuthentication
		}

		return &user, nil
	}

	// User already exists, return user data
	return u, nil
}
