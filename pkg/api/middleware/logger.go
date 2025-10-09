package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func setLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof(
			"Обработка %s-запроса завершилась с кодом [%d] при обращении к \"%s\" с IP-адреса %s",
			c.Request.Method,
			c.Writer.Status(),
			c.Request.URL.Path,
			c.ClientIP(),
		)
		c.Next()
	}
}

func setContextLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("clog", logger)
		c.Next()
	}
}
