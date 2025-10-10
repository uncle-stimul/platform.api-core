package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/db/methods"
	dmodels "platform.api-core/pkg/db/models"
	"platform.api-core/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	users, err := methods.SelectAll[dmodels.Users](pgdb.Preload("Roles"))
	if err != nil {
		log.WithError(err).Error("При получении выборки пользователей возникла ошибка:")
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    "При получении выборки пользователей возникла ошибка",
		})
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
		log.WithError(err).Error("Получен некорректный JSON для вывода пользователя:")
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    "Получен некорректный JSON для вывода пользователя",
		})
		return
	}

	if req.ID == 0 {
		log.Error("Полученный JSON не содержит идентификатора пользователя")
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    "ID пользователя обязателен",
		})
		return
	}

	var user models.Users
	if err := pgdb.Preload("Roles").First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Пользователь с ID %d не найден", req.ID)
			c.JSON(http.StatusNotFound, models.DefaultResponse{
				Status: "error",
				Msg:    fmt.Sprintf("Пользователь с ID %d не найден", req.ID),
			})
		} else {
			log.WithError(err).Errorf("При поиске пользователя с ID: %d возникла ошибка", req.ID)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    fmt.Sprintf("При поиске пользователя с ID: %d возникла ошибка", req.ID),
			})
		}
		return
	}

	res := models.GetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Status:   user.Status,
		Roles:    utils.ExtractRoleIDs(user.Roles),
	}

	log.Infof("Предоставлена информация по пользователю с ID: %d", req.ID)
	c.JSON(http.StatusOK, res)
}
