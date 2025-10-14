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

func UpdateUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для обновления пользователя:"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := fmt.Sprintf("Полученный JSON, который не содедржит идентификатора пользователя %s", req.Username)
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingUser models.Users
		if err := tx.First(&existingUser, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("Пользователь с ID: %d не найден в базе данных", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			msg := fmt.Sprintf("Поиск пользователя с ID: %d завершился ошибкой:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		user := models.Users{
			ID:       existingUser.ID,
			Username: req.Username,
			Password: req.Password,
			Status:   req.Status,
		}

		if err := tx.Save(&user).Error; err != nil {
			msg := fmt.Sprintf("При обновлении пользователя с ID: %d возникла ошибка:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			msg := fmt.Sprintf("При очистке ролей пользователя с ID: %d возникла ошибка:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		var rolesToSet []models.Roles
		if len(req.Roles) > 0 {
			for _, id := range req.Roles {
				rolesToSet = append(rolesToSet, models.Roles{ID: id})
			}

			if err := tx.Where("id IN ?", req.Roles).Find(&rolesToSet).Error; err != nil || len(rolesToSet) != len(req.Roles) {
				msg := "Указанные в запросе роли не существуют"
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}

			if err := tx.Model(&user).Association("Roles").Append(rolesToSet); err != nil {
				msg := fmt.Sprintf("При добавлении ролей пользователю с ID: %d возникла ошибка:", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
		} else {
			if err := tx.Preload("Roles").First(&existingUser, req.ID).Error; err != nil {
				msg := fmt.Sprintf("При загрузке текущих ролей пользователя с ID: %d возникла ошибка:", req.ID)
				log.WithError(err).Error(msg)
				return fmt.Errorf("%s", msg)
			}
			rolesToSet = existingUser.Roles
			if len(rolesToSet) > 0 {
				if err := tx.Model(&user).Association("Roles").Append(rolesToSet); err != nil {
					msg := fmt.Sprintf("При восстановлении ролей для пользователя с ID: %d возникла ошибка:", req.ID)
					log.WithError(err).Error(msg)
					return fmt.Errorf("%s", msg)
				}
			}
		}

		if err := tx.Preload("Roles").First(&user, req.ID).Error; err != nil {
			msg := fmt.Sprintf("При восстановлении ролей для пользователя с ID: %d возникла ошибка:", req.ID)
			log.WithError(err).Error(msg)
			return fmt.Errorf("%s", msg)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    fmt.Sprintf("Обновление пользователя c ID: %d привело к ошибке", req.ID),
		})
		return
	} else {
		msg := fmt.Sprintf("Пользователь c ID: %d успешно обновлен", req.ID)
		log.Info(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
