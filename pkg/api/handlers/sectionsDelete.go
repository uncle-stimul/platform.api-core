package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func deleteSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для конечной точки"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатора конечной точки"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.First(&models.Sections{}, req.ID).Error; err != nil {
		msg := fmt.Sprintf("Конечная точка с ID: %d не найден в базе данных", req.ID)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusNotFound, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.Delete(&models.Users{ID: req.ID}).Error; err != nil {
		msg := fmt.Sprintf("При удалении конечной точки с ID %d возникла ошибка", req.ID)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Конечная точка с ID %d успешно удалена", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
