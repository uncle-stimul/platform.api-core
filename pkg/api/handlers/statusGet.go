package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	Status  int    `json:"status"`
	Module  string `json:"module"`
	Version string `json:"version"`
	Message string `json:"msg"`
}

func getStatus(c *gin.Context) {
	var res = statusResponse{
		Status:  http.StatusOK,
		Module:  "api-platform",
		Version: "0.7.d",
		Message: "Модуль управляет пользователями и предоставлением доступа, а также реализует функционал для авторизации",
	}
	c.IndentedJSON(http.StatusOK, res)
}
