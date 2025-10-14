package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func deleteRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для удаления роли"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "В полеченном JSON-запросе отсутствует идентификатор роли"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.First(&models.Roles{}, req.ID).Error; err != nil {
		msg := fmt.Sprintf("При поиске роли с ID: %d возникла не предвиденная ошибка", req.ID)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusNotFound, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Роль с ID: %d не найдена в базе данных", req.ID),
		})
		return
	}

	if err := pgdb.Delete(&models.Roles{ID: req.ID}).Error; err != nil {
		msg := "При удалении пользовательской роли возникла ошибка:"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Пользовательская роль с ID %d успешно удалена", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
