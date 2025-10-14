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

func selectSections(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var sections []models.Sections
	err := pgdb.Find(&sections).Error
	if err != nil {
		msg := "При получении выборки конечных точек возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	result := make([]models.GetSectionsResponse, 0, len(sections))
	for _, section := range sections {
		result = append(result, models.GetSectionsResponse{
			Module:   section.Module,
			Endpoint: section.Endpoint,
		})
	}

	c.JSON(http.StatusOK, result)
}

func selectSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для вывода конечной точки"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатор конечной точки"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	var section models.Sections
	if err := pgdb.First(&section, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Конечная точка с ID %d не найдена в базе данных", req.ID)
			log.Error(msg)
			c.JSON(http.StatusNotFound, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		} else {
			msg := fmt.Sprintf("При поиске конечной точки с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		}
		return
	}

	res := models.GetSectionsResponse{
		Module:   section.Module,
		Endpoint: section.Endpoint,
	}

	log.Infof("Предоставлена информация о конечной точке с ID: %d", req.ID)
	c.JSON(http.StatusOK, res)
}
