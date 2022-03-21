package api

import (
	"net/http"

	"github.com/Atelier-Epita/intra-atelier/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleUser() {
	users := router.Group("/users")

	users.GET("", GetUsersHandler)
	users.POST("", CreateUserHandler)

	users.GET("/:email", GetUserByEmailHandler)

	users.POST("/:email/:groupName", AddGroupToUserHandler)
	users.GET("/:email/groups", GetUserGroupsHandler)
}

// @Summary Get users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 "Couldn't get users"
// @Router /users [GET]
func GetUsersHandler(c *gin.Context) {
	zap.S().Info("Getting all users...")

	users, err := models.GetUsers()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get users")
		return
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
// @Accept json
// @Param request body UserRequest true "UserRequest"
// @Success 200 "OK"
// @Failure 400 "Bad request"
// @Failure 500 "Couldn't create user"
// @Router /users [POST]
func CreateUserHandler(c *gin.Context) {
	var userRequest UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		Abort(c, err, http.StatusBadRequest, "Invalid user request")
		return
	}

	u := models.User{Email: userRequest.Email, FirstName: userRequest.FirstName, LastName: userRequest.LastName}

	if err := u.Insert(); err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created user",
	})
}

// @Summary Get user by mail
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 "Not found"
// @Router /users/{email} [GET]
// @Param email path string true "email"
func GetUserByEmailHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByMail(mail)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User using "+mail+" not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Unused
func GetUserByNameHandler(c *gin.Context) {
	name := c.Param("name")
	user, err := models.GetUserByMail(name)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User "+name+" not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get user by mail
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 "User of group not found"
// @Failure 500 "Couldn't add group"
// @Router /users/{email}/{GroupName} [POST]
// @Param UserMail path string true "User mail"
// @Param GroupName path string true "GroupName"
func AddGroupToUserHandler(c *gin.Context) {
	email := c.Param("email")
	group_name := c.Param("groupName")
	user, err := models.GetUserByMail(email)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User "+email+" not found")
		return
	}

	group, err := models.GetGroup(group_name)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "Group "+group_name+" not found")
		return
	}

	if err := user.AddGroup(group); err != nil {
		Abort(c, err, http.StatusInternalServerError, "Group "+group_name+" couldnt be inserted")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully added group to user",
	})
}

// @Summary Get user groups by user email
// @Tags users
// @Produce json
// @Success 200 {array} models.Group
// @Failure 404 "User Not found"
// @Failure 500 "Couldn't get groups"
// @Router /users/{email}/groups [GET]
// @Param email path string true "email"
func GetUserGroupsHandler(c *gin.Context) {
	mail := c.Param("email")
	user, err := models.GetUserByMail(mail)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User using "+mail+" not found")
		return
	}

	groups, err := user.GetGroups()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get groups for user "+user.Email)
		return
	}

	c.JSON(http.StatusOK, groups)
}
