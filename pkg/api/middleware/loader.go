package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AddMiddleware(router *gin.Engine, log *logrus.Logger, pgdb *gorm.DB) {
	router.Use(setLoggerMiddleware(log))
	router.Use(setContextLoggerMiddleware(log))
	router.Use(DBMiddleware(pgdb, log))
	router.Use(gin.Recovery())
}
