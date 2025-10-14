package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func createSection(c *gin.Context) {
	var req models.CreateSectionRequest
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для создания конечной точки"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.Module == "" || req.Endpoint == "" {
		msg := "В полученном JSON не указаны модуль и конечная точки"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	section := models.Sections{
		Module:   req.Module,
		Endpoint: req.Endpoint,
	}

	if err := pgdb.Create(&section).Error; err != nil {
		msg := "При создании конечной точки возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Было выполнено создание конечной точки \"%s\"", section.Endpoint)
		log.Info(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
