package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func DeleteUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Получен некорректный JSON для удаления пользователя:")
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    "Получен некорректный JSON для удаления пользователя",
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

	if req.ID == 1 {
		log.Error("Невозможно удалить встроенного администратора")
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    "Невозможно удалить встроенного администратора",
		})
		return
	}

	if err := pgdb.First(&models.Users{}, req.ID).Error; err != nil {
		log.WithError(err).Errorf("При поиске пользователя с ID: %d возникла не предвиденная ошибка", req.ID)
		c.JSON(http.StatusNotFound, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Пользователь с ID: %d не найден в базе данных", req.ID),
		})
		return
	}

	if err := pgdb.Delete(&models.Users{ID: req.ID}).Error; err != nil {
		log.WithError(err).Error("При создании пользователя возникла ошибка:")
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    "При создании пользователя возникла ошибка",
		})
		return
	} else {
		log.Infof("Пользователь с ID %d успешно удален", req.ID)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    fmt.Sprintf("Пользователь с ID %d успешно удален", req.ID),
		})
	}
}
