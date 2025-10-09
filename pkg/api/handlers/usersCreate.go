package handlers

import (
	"net/http"

	"platform.api-core/pkg/db/methods"
	"platform.api-core/pkg/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var request models.Users

	dbCntx, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Контекс не содержит информации о БД"})
		return
	}
	pgdb := dbCntx.(*gorm.DB)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат данных"})
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле username обязательно"})
		return
	}

	checker := map[string]interface{}{"username": request.Username}
	created, err := methods.Create[models.Users](pgdb, checker, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !created {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким username уже существует"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": request})
}
