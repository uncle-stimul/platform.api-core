package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func updateSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для обновления конечной точки"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "В полеченном JSON-запросе отсутствует идентификатор конечной точки"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingSection models.Sections
		if err := tx.First(&existingSection, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("Конечная точка с ID: %d не найдена в базе данных", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			msg := fmt.Sprintf("Поиск конечной точки с ID: %d завершился ошибкой", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		role := models.Sections{
			ID:       existingSection.ID,
			Module:   req.Module,
			Endpoint: req.Endpoint,
		}

		if err := tx.Save(&role).Error; err != nil {
			msg := fmt.Sprintf("При обновлении конечной точки с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Обновление конечной точки c ID: %d привело к ошибке", req.ID),
		})
		return
	} else {
		msg := fmt.Sprintf("Конечная точка c ID: %d успешно обновлена", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
