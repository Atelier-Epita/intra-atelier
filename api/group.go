package api

import (
	"net/http"

	"github.com/Atelier-Epita/intra-atelier/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleGroup() {
	users := router.Group("/groups")

	users.GET("", GetGroupsHandler)
	users.GET("/:name", GetGroupByNameHandler)

	users.POST("/:name", CreateGroupHandler)
}

func GetGroupsHandler(c *gin.Context) {
	zap.S().Info("Getting all groups...")

	groups, err := models.GetGroups()
	if err != nil {
		Abort(c, err, http.StatusBadRequest, "Couldn't get groups")
		return
	}

	c.JSON(http.StatusOK, groups)
}

// @Summary Create group
// @Tags users
// @Sucess 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 500 "Server error"
// @Router /users [POST]
func CreateGroupHandler(c *gin.Context) {
	g := models.Group{Name: c.Param("name")}
	if err := g.Insert(); err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't create group")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created group",
	})
}

func GetGroupByNameHandler(c *gin.Context) {
	name := c.Param("name")
	group, err := models.GetGroup(name)
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Group "+name+" not found")
		return
	}

	c.JSON(http.StatusOK, group)
}
