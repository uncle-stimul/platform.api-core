package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/models"
)

func GetHealthcheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "running"})
}

func GetModuleInfo(c *gin.Context) {
	var res = models.ModuleInfoResponse{
		Status:  "success",
		Module:  models.MetadataModule,
		Version: models.MetadataVersion,
		Message: models.MetadataMessage,
	}
	c.IndentedJSON(http.StatusOK, res)
}
