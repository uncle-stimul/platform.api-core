package handlers

import (
	"net/http"

	"platform.api-core/pkg/db/methods"
	"platform.api-core/pkg/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type rolesResponse struct {
	Name        string   `json:"role"`
	Description string   `json:"descriptions"`
	Permissions []string `json:"permissions"`
}

func getRoles(c *gin.Context) {
	dbCntx, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Контекс не содержит информации о БД"})
		return
	}
	pgdb := dbCntx.(*gorm.DB)

	roles, err := methods.SelectAll[models.Roles](pgdb.Preload("Permissions"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]rolesResponse, 0, len(roles))
	for _, role := range roles {
		permissions := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission.Name)
		}
		result = append(result, rolesResponse{
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissions,
		})
	}

	c.JSON(http.StatusOK, result)
}
