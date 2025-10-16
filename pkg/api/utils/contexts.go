package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AddContextDB(c *gin.Context) *gorm.DB {
	cdb, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "Контекс не содержит информации о БД",
		})
		return nil
	}
	return cdb.(*gorm.DB)
}

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
