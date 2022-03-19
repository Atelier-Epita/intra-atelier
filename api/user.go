package api

import (
	"intra/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleUser() {
	users := router.Group("/users")

	users.GET("", GetUsersHandler)
	users.GET("/:email", GetUserByEmailHandler)
	users.GET("/:email/groups", GetUserGroupsHandler)

	users.POST("", CreateUserHandler)
	users.POST("/:name/:groupName", AddGroupToUserHandler)
}

// @Summary Get users
// @Tags users
// @Sucess 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 500 "Server error"
// @Router /users [GET]
func GetUsersHandler(c *gin.Context) {
	zap.S().Info("Getting all users...")

	users, err := models.GetUsers()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Couldn't get users",
		})
	}

	c.JSON(http.StatusOK, users)
}

type UserRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
}

// @Summary Create user
// @Tags users
// @Sucess 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 500 "Server error"
// @Router /users [POST]
func CreateUserHandler(c *gin.Context) {
	var userRequest UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid user request",
		})
	} else {
		u := models.User{Email: userRequest.Email, FirstName: userRequest.FirstName, LastName: userRequest.LastName}
		u.Insert()

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Successfully created user",
		})
	}
}

func GetUserByEmailHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByEmail(mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User using " + mail + " not found",
		})
	}

	c.JSON(http.StatusOK, user)
}

func GetUserByNameHandler(c *gin.Context) {
	name := c.Param("name")
	user, err := models.GetUserByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User " + name + " not found",
		})
	}

	c.JSON(http.StatusOK, user)
}

func AddGroupToUserHandler(c *gin.Context) {
	name := c.Param("name")
	group_name := c.Param("groupName")
	user, err := models.GetUserByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User " + name + " not found",
		})
	}
	group, err := models.GetGroupByName(group_name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Group " + name + " not found",
		})
	}
	err = user.AddGroup(*group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusNotFound,
			"message": "Group " + name + " couldnt be inserted",
		})
	}

	c.JSON(http.StatusOK, user)
}

func GetUserGroupsHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByEmail(mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User using " + mail + " not found",
		})
	}

	c.JSON(http.StatusOK, user.Groups)
}
