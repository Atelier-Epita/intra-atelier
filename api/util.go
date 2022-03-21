package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Abort(c *gin.Context, err error, status int, message string) {
	if err != nil {
		zap.S().Error(err)
	}
	c.AbortWithStatusJSON(status, gin.H{
		"status":  status,
		"message": message,
	})
}
