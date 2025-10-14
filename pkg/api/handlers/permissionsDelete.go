package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func deletePermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для удаления права доступа"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "В полеченном JSON-запросе отсутствует идентификатор права доступа"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if err := pgdb.First(&models.Permissions{}, req.ID).Error; err != nil {
		msg := fmt.Sprintf("При поиске права доступа с ID: %d возникла не предвиденная ошибка", req.ID)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusNotFound, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Право доступа с ID: %d не найдено в базе данных", req.ID),
		})
		return
	}

	if err := pgdb.Delete(&models.Permissions{ID: req.ID}).Error; err != nil {
		msg := "При удалении права доступа возникла ошибка:"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Право доступа с ID %d успешно удалено", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
