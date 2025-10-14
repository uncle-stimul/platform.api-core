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

func updateRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для обновления пользовательской роли"
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

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingRole models.Roles
		if err := tx.First(&existingRole, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("Пользовательская роль с ID: %d не найдена в базе данных", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			msg := fmt.Sprintf("Поиск пользовательской роли с ID: %d завершился ошибкой:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		role := models.Roles{
			ID:          existingRole.ID,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := tx.Save(&role).Error; err != nil {
			msg := fmt.Sprintf("При обновлении пользователя с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			msg := fmt.Sprintf("При очистке прав доступа для пользовательской роли с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		var permissionsToSet []models.Permissions
		if len(req.Permissions) > 0 {
			for _, id := range req.Permissions {
				permissionsToSet = append(permissionsToSet, models.Permissions{ID: id})
			}

			if err := tx.Where("id IN ?", req.Permissions).Find(&permissionsToSet).Error; err != nil || len(permissionsToSet) != len(req.Permissions) {
				msg := "Указанные в запросе права доступа не существуют"
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}

			if err := tx.Model(&role).Association("Permissions").Append(permissionsToSet); err != nil {
				msg := fmt.Sprintf("При добавлении прав доступа для пользовательской роли с ID: %d возникла ошибка", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
		} else {
			if err := tx.Preload("Permissions").First(&existingRole, req.ID).Error; err != nil {
				msg := fmt.Sprintf("При загрузке текущих прав доступа для пользовательской роли с ID: %d возникла ошибка", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			permissionsToSet = existingRole.Permissions
			if len(permissionsToSet) > 0 {
				if err := tx.Model(&role).Association("Permissions").Append(permissionsToSet); err != nil {
					msg := fmt.Sprintf("При восстановлении прав доступа для пользовательской роли с ID: %d возникла ошибка", req.ID)
					log.WithError(err).Error(msg)
					return fmt.Errorf("%s", msg)
				}
			}
		}

		if err := tx.Preload("Permissions").First(&role, req.ID).Error; err != nil {
			msg := fmt.Sprintf("При восстановлении прав доступа для пользовательской роли с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Обновление пользовательской роли c ID: %d привело к ошибке", req.ID),
		})
		return
	} else {
		msg := fmt.Sprintf("Пользовательская роль c ID: %d успешно обновлена", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}

}
