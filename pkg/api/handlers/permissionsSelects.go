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

func selectPermissions(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var permissions []models.Permissions
	err := pgdb.Preload("Sections").Find(&permissions).Error
	if err != nil {
		msg := "При получении выборки прав доступа возникла ошибка:"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
	}

	result := make([]models.GetPermissionsResponse, 0, len(permissions))
	for _, permission := range permissions {
		sections := make([]string, 0, len(permission.Sections))
		for _, section := range permission.Sections {
			sections = append(sections, section.Endpoint)
		}
		result = append(result, models.GetPermissionsResponse{
			Name:        permission.Name,
			Description: permission.Description,
			Sections:    sections,
		})
	}

	c.JSON(http.StatusOK, result)
}

func selectPermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для вывода прав доступа"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатор прав доступа"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	var permission models.Permissions
	if err := pgdb.Preload("Sections").First(&permission, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Право досутпа с ID %d не найдено в базе данных", req.ID)
			log.Error(msg)
			c.JSON(http.StatusNotFound, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		} else {
			msg := fmt.Sprintf("При поиске прав доступа с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		}
		return
	}

	res := models.GetRoleResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		Permissions: utils.ExtractSectionsEndpoints(permission.Sections),
	}

	log.Infof("Предоставлена информация о правах доступа с ID: %d", req.ID)
	c.JSON(http.StatusOK, res)
}
