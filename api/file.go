package api

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Atelier-Epita/intra-atelier/models"
	"github.com/gin-gonic/gin"
)

func handleFile() {
	files := router.Group("/files")

	files.GET("/:id", GetFileById)
	files.POST("/:email/:filename", CreateFileHandler)
}

// @Summary Get file
// @Tags files
// @Produce json
// @Success 200 {object} models.File
// @Failure 400 "File id invalid"
// @Failure 500 "Couldn't get file"
// @Router /files/{id} [GET]
// @Param id path string true "File id"
func GetFileById(c *gin.Context) {
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		Abort(c, err, http.StatusBadRequest, "file id invalid")
		return
	}

	file, err := models.GetFileById(id)
	if err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't get file")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File("./files/" + hex.EncodeToString(file.Filehash))
}

// @Summary Create file
// @Tags files
// @Produce json
// @Success 200 "OK"
// @Failure 400 "Invalid filename or couldn't read file request"
// @Failure 404 "User not found"
// @Failure 500 "Couldn't create file"
// @Router /files/{email}/{filename} [POST]
// @Param email path string true "Email"
// @Param filename path string true "Filename"
// @Param "File Request"
func CreateFileHandler(c *gin.Context) {
	email := c.Param("email")
	filename := c.Param("filename")

	if strings.Contains(filename, "/") {
		Abort(c, nil, http.StatusBadRequest, "Filename invalid")
		return
	}

	u, err := models.GetUserByMail(email)
	if err != nil {
		Abort(c, err, http.StatusNotFound, "User "+email+" not found")
		return
	}

	form, err := c.FormFile("file")
	if err != nil {
		Abort(c, err, http.StatusBadRequest, "Couldn't read file request form")
		return
	}
	openedFile, err := form.Open()
	if err != nil {
		Abort(c, err, http.StatusBadRequest, "Couldn't open file request")
		return
	}

	fileContent, err := ioutil.ReadAll(openedFile)
	if err != nil {
		Abort(c, err, http.StatusBadRequest, "Couldn't read file request")
		return
	}

	var f = models.File{Permission: 0, OwnerID: u.Id, Filename: filename}
	if err := models.CreateFile(f, fileContent); err != nil {
		Abort(c, err, http.StatusInternalServerError, "Couldn't create file")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created file",
	})
}
