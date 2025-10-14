package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var users []models.Users
	err := pgdb.Preload("Roles").Find(&users).Error
	if err != nil {
		msg := "При получении выборки пользователей возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return
	}

	result := make([]models.UsersResponse, 0, len(users))
	for _, user := range users {
		roles := make([]string, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
		result = append(result, models.UsersResponse{
			Username: user.Username,
			Roles:    roles,
			Status:   user.Status,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для вывода пользователя"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатора пользователя"
		log.Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return
	}

	var user models.Users
	if err := pgdb.Preload("Roles").First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Пользователь с ID %d не найден", req.ID)
			log.Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		} else {
			msg := fmt.Sprintf("При поиске пользователя с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		}
		return
	}

	res := models.GetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Status:   user.Status,
		Roles:    utils.ExtractRolesNames(user.Roles),
	}

	log.Infof("Предоставлена информация по пользователю с ID: %d", req.ID)
	c.JSON(http.StatusOK, res)
}
