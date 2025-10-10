package handlers

import (
	"net/http"

	"platform.api-core/pkg/db/methods"
	"platform.api-core/pkg/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type permissionsResponse struct {
	Name        string   `json:"permissions"`
	Description string   `json:"descriptions"`
	Sections    []string `json:"sections"`
}

func getPermissions(c *gin.Context) {
	dbCntx, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Контекс не содержит информации о БД"})
		return
	}
	pgdb := dbCntx.(*gorm.DB)

	permissions, err := methods.SelectAll[models.Permissions](pgdb.Preload("Sections"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]permissionsResponse, 0, len(permissions))
	for _, permission := range permissions {
		sections := make([]string, 0, len(permission.Sections))
		for _, section := range permission.Sections {
			sections = append(sections, section.URL)
		}
		result = append(result, permissionsResponse{
			Name:        permission.Name,
			Description: permission.Description,
			Sections:    sections,
		})
	}

	c.JSON(http.StatusOK, result)
}
