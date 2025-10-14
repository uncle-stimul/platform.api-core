package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func createRole(c *gin.Context) {
	var req models.CreateRoleRequest
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для создания роли"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.Name == "" {
		msg := "В полученном JSON не указано наименование роли"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	role := models.Roles{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := pgdb.Create(&role).Error; err != nil {
		msg := "При создании роли возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if len(req.Permissions) > 0 {
		var permissions []*models.Permissions
		if err := pgdb.Where("id IN ?", req.Permissions).Find(&permissions).Error; err != nil {
			msg := "Указанные права доступа не найдены в базе данных"
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
		if err := pgdb.Model(&role).Association("Permissions").Append(permissions); err != nil {
			msg := fmt.Sprintf("При назначении при назначении прав доступа для роли \"%s\" возникла не предвиденная ошибка", role.Name)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
	}

	if err := pgdb.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		msg := fmt.Sprintf("При загрузке прав доступа для роли \"%s\" возникла ошибка", role.Name)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Было выполнено создание роли \"%s\"", role.Name)
		log.Info(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
