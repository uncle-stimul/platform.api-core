package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AddContextLogger(c *gin.Context) *logrus.Logger {
	clog, exists := c.Get("clog")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "Контекс не содержит информации о фиксации событий"},
		)
		return nil
	}
	return clog.(*logrus.Logger)
}
