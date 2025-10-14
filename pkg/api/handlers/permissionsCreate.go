package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func createPermission(c *gin.Context) {
	var req models.CreatePermissionRequest
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для создания права доступа"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.Name == "" {
		msg := "В полученном JSON не указано наименование права доступа"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	permission := models.Permissions{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := pgdb.Create(&permission).Error; err != nil {
		msg := "При создании права доступа возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if len(req.Sections) > 0 {
		var sections []*models.Sections
		if err := pgdb.Where("id IN ?", req.Sections).Find(&sections).Error; err != nil {
			msg := "Указанные конечные точки не найдены в базе данных"
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
		if err := pgdb.Model(&permission).Association("Sections").Append(sections); err != nil {
			msg := fmt.Sprintf("При назначении при назначении конечных точек для права доступа \"%s\" возникла не предвиденная ошибка", permission.Name)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
	}

	if err := pgdb.Preload("Sections").First(&permission, permission.ID).Error; err != nil {
		msg := fmt.Sprintf("При загрузке конечных точек для права доступа \"%s\" возникла ошибка", permission.Name)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Было выполнено создание права доступа \"%s\"", permission.Name)
		log.Info(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
