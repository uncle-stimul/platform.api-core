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

func updatePermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для обновления права доступа"
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

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingPermission models.Permissions
		if err := tx.First(&existingPermission, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("Право доступа с ID: %d не найдено в базе данных", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			msg := fmt.Sprintf("Поиск права доступа с ID: %d завершился ошибкой:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		permission := models.Permissions{
			ID:          existingPermission.ID,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := tx.Save(&permission).Error; err != nil {
			msg := fmt.Sprintf("При обновлении права доступа с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		if err := tx.Model(&permission).Association("Sections").Clear(); err != nil {
			msg := fmt.Sprintf("При очистке конечных точек для права доступа с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		var sectionsToSet []models.Sections
		if len(req.Sections) > 0 {
			for _, id := range req.Sections {
				sectionsToSet = append(sectionsToSet, models.Sections{ID: id})
			}

			if err := tx.Where("id IN ?", req.Sections).Find(&sectionsToSet).Error; err != nil || len(sectionsToSet) != len(req.Sections) {
				msg := "Указанные в запросе конечные точки не существуют"
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}

			if err := tx.Model(&permission).Association("Sections").Append(sectionsToSet); err != nil {
				msg := fmt.Sprintf("При добавлении конечных точек для права доступа с ID: %d возникла ошибка", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
		} else {
			if err := tx.Preload("Sections").First(&existingPermission, req.ID).Error; err != nil {
				msg := fmt.Sprintf("При загрузке текущих конечных точек для права доступа с ID: %d возникла ошибка", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			sectionsToSet = existingPermission.Sections
			if len(sectionsToSet) > 0 {
				if err := tx.Model(&permission).Association("Sections").Append(sectionsToSet); err != nil {
					msg := fmt.Sprintf("При восстановлении конечных точек для права доступа с ID: %d возникла ошибка", req.ID)
					log.WithError(err).Error(msg)
					return fmt.Errorf("%s", msg)
				}
			}
		}

		if err := tx.Preload("Sections").First(&permission, req.ID).Error; err != nil {
			msg := fmt.Sprintf("При восстановлении конечных точек для права доступа с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Обновление права доступа c ID: %d привело к ошибке", req.ID),
		})
		return
	} else {
		msg := fmt.Sprintf("Право доступо c ID: %d успешно обновлено", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
