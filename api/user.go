package api

import (
	"intra/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleUser() {
	users := router.Group("/users")

	users.GET("", GetUsers)
	users.POST("", CreateUser)
}

func GetUsers(c *gin.Context) {
	zap.S().Info("Getting all users...")

	users, err := models.GetUsers()
	if (err != nil) {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "Couldn't get users",
		})
	}

	c.JSON(http.StatusOK, users)
}

type UserRequest struct {
	Login string `json:"login" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	Name string `json:"name" binding:"required"`
	Promotion uint16 `json:"promotion" binding:"required"`
}

func ValidateRequest(userRequest UserRequest) bool {
	substrings := strings.Split(userRequest.Login, ".")
	if len(substrings) != 2 || substrings[0] == "" || substrings[1] == "" {
		return false
	}
	if userRequest.Promotion > uint16(time.Now().Year()) + 5 {
		return false
	}

	return true
}

func CreateUser(c *gin.Context) {
	var userRequest UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "Invalid user request",
		})
	} else {
		if !ValidateRequest(userRequest) {
			zap.S().Info("Invalid user request %#v", userRequest)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"message": "Invalid user request",
			})
		}

		models.Insert(userRequest.Login,
			userRequest.FirstName,
			userRequest.Name,
			userRequest.Promotion,
		)

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "Successfully created user " + userRequest.Login,
		})
	}

}