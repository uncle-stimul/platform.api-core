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
		msg := "Получен некорректный JSON для удаления пользователя:"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатора пользователя"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 1 {
		msg := "Невозможно удалить встроенного администратора"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.First(&models.Users{}, req.ID).Error; err != nil {
		msg := fmt.Sprintf("При поиске пользователя с ID: %d возникла не предвиденная ошибка", req.ID)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusNotFound, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.Delete(&models.Users{ID: req.ID}).Error; err != nil {
		msg := "При удалении пользователя возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Пользователь с ID %d успешно удален", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
