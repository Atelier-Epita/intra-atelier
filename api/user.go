package api

import (
	"fmt"
	"net/http"

	"github.com/Atelier-Epita/intra-atelier/models"

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
	user, err := models.GetUserByMail(mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User using " + mail + " not found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserByNameHandler(c *gin.Context) {
	name := c.Param("name")
	user, err := models.GetUserByMail(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User " + name + " not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AddGroupToUserHandler(c *gin.Context) {
	email := c.Param("email")
	group_name := c.Param("groupName")
	user, err := models.GetUserByMail(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User " + email + " not found",
		})
		return
	}
	group, err := models.GetGroup(group_name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Group " + email + " not found",
		})
		return
	}
	err = user.AddGroup(group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Group " + email + " couldnt be inserted",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserGroupsHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByMail(mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User using " + mail + " not found",
		})
	}
	groups, err := user.GetGroups()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Couldn't get groups for user " + user.Email,
		})
	} else {
		c.JSON(http.StatusOK, groups)
	}

}
