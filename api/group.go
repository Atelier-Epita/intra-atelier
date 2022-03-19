package api

import (
	"intra/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleGroup() {
	users := router.Group("/groups")

	users.GET("", GetGroupsHandler)
	users.GET("/:name", GetGroupByNameHandler)

	users.POST("", CreateGroupHandler)
}

func GetGroupsHandler(c *gin.Context) {
	zap.S().Info("Getting all groups...")

	groups, err := models.GetGroups()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Couldn't get groups",
		})
	}

	c.JSON(http.StatusOK, groups)
}

type GroupRequest struct {
	Name string `json:"name"`
}

// @Summary Create group
// @Tags users
// @Sucess 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 500 "Server error"
// @Router /users [POST]
func CreateGroupHandler(c *gin.Context) {
	var groupRequest GroupRequest
	if err := c.BindJSON(&groupRequest); err != nil {
		zap.S().Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid group request",
		})
	} else {
		g := models.Group{Name: groupRequest.Name}
		g.Insert()

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Successfully created group",
		})
	}
}

func GetGroupByNameHandler(c *gin.Context) {
	name := c.Param("name")
	group, err := models.GetGroupByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Group " + name + " not found",
		})
	}

	c.JSON(http.StatusOK, group)
}
