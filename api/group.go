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

// @Summary Get all groups
// @Tags groups
// @Success 200 {array} models.Group
// @Failure 500 "Couldn't get groups"
// @Router /groups [GET]
func GetGroupsHandler(c *gin.Context) {
	zap.S().Info("Getting all groups...")

	groups, err := models.GetGroups()
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get groups")
		return
	}

	c.JSON(http.StatusOK, groups)
}

// @Summary Create group
// @Tags groups
// @Success 200 "OK"
// @Failure 500 "Couldn't create group"
// @Router /groups/{GroupName} [POST]
// @Param GroupName path string true "GroupName"
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

// @Summary Get group by name
// @Tags groups
// @Success 200 {object} models.Group
// @Failure 404 "Group not found"
// @Router /groups/{GroupName} [GET]
// @Param GroupName path string true "GroupName"
func GetGroupByNameHandler(c *gin.Context) {
	name := c.Param("name")
	group, err := models.GetGroup(name)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "Group "+name+" not found")
		return
	}

	c.JSON(http.StatusOK, group)
}
