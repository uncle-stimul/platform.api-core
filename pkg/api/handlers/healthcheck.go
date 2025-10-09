package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status string `json:status`
}

func getHealthcheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, healthResponse{Status: "running"})
}
