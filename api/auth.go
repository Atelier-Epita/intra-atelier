package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Atelier-Epita/intra-atelier/models"
	"github.com/aureleoules/epitaf/lib/microsoft"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleAuth() {
	users := router.Group("/users")

	users.GET("/authenticate", authenticateHandler)
	users.GET("/callback", callbackHandler)
}

// @Summary Authenticate URL
// @Tags auth
// @Description Build Microsoft oauth url
// @Param   redirect_uri	body	string	true	"redirect_uri"  default(http://localhost:8080/v0.1/users/callback)
// @Success 200	"OK"
// @Failure 406	"Not acceptable"
// @Router /users/authenticate [POST]
func authenticateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, microsoft.SignInURL("http://localhost:8080/v0.1/users/callback"))
}

// @Summary OAuth Callback
// @Description Authenticate user
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
func callbackHandler(c *gin.Context) {
	var code = c.Request.URL.Query().Get("code")
	var redirect_uri = "http://localhost:8080/v0.1/users/callback"
	fmt.Printf(code, redirect_uri)

	token, err := microsoft.GetAccessToken(code, redirect_uri)
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't authenticate")
		return
	}

	client := microsoft.NewClient(token, nil)
	// Retrieve microsoft profile
	profile, err := client.GetProfile()
	if err != nil {
		zap.S().Error(err)
		Abort(c, err, http.StatusInternalServerError, "Couldn't authenticate")
		return
	}

	// Check if user exists in database
	u, err := models.GetUserByMail(profile.Mail)
	if err != nil {
		user := models.User{
			Email:     profile.Mail,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
		}

		// Insert new user and login
		err = user.Insert()
		if err != nil {
			zap.S().Error(err)
			Abort(c, err, http.StatusInternalServerError, "Couldn't authenticate")
			return
		}

		onLogin(c, &user)
	}

	// User already exists, login
	onLogin(c, u)
}

// Attach cookie to user
func onLogin(c *gin.Context, u *models.User) {
	cookieName := os.Getenv("COOKIE_NAME")
	cookieExpiration, err := strconv.Atoi(os.Getenv("COOKIE_EXPIRATION"))
	if err != nil {
		zap.S().Error(err)
		Abort(c, err, http.StatusInternalServerError, err.Error())
		return
	}

	cookieSecure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		zap.S().Error(err)
		Abort(c, err, http.StatusInternalServerError, err.Error())
		return
	}

	encoded, err := SC.Encode(cookieName, u.Email)
	if err != nil {
		zap.S().Error(err)
		Abort(c, err, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(
		cookieName,
		encoded,
		cookieExpiration,
		os.Getenv("COOKIE_DOMAIN"),
		"/",
		cookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully login",
	})
}
