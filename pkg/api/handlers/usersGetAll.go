package handlers

import (
	"net/http"

	"platform.api-core/pkg/db/methods"
	"platform.api-core/pkg/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type usersResponse struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Status   bool     `json:"status"`
}

func getUsers(c *gin.Context) {
	dbCntx, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Контекс не содержит информации о БД"})
		return
	}
	pgdb := dbCntx.(*gorm.DB)

	users, err := methods.SelectAll[models.Users](pgdb.Preload("Roles"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]usersResponse, 0, len(users))
	for _, user := range users {
		roles := make([]string, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
		result = append(result, usersResponse{
			Username: user.Username,
			Roles:    roles,
			Status:   user.Status,
		})
	}

	c.JSON(http.StatusOK, result)
}
