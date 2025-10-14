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

func selectRoles(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var roles []models.Roles
	err := pgdb.Preload("Permissions").Find(&roles).Error
	if err != nil {
		log.WithError(err).Error("При получении выборки ролей возникла ошибка:")
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    "При получении выборки ролей возникла ошибка",
		})
		return
	}

	result := make([]models.GetRolesResponse, 0, len(roles))
	for _, role := range roles {
		permissions := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission.Name)
		}
		result = append(result, models.GetRolesResponse{
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissions,
		})
	}

	c.JSON(http.StatusOK, result)
}

func selectRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для вывода роли"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.ID == 0 {
		msg := "Полученный JSON не содержит идентификатор роли"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	var role models.Roles
	if err := pgdb.Preload("Permissions").First(&role, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Роль с ID %d не найдена в базе данных", req.ID)
			log.Error(msg)
			c.JSON(http.StatusNotFound, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		} else {
			msg := fmt.Sprintf("При поиске роли с ID: %d возникла ошибка", req.ID)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
		}
		return
	}

	res := models.GetRoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: utils.ExtractPermissionsNames(role.Permissions),
	}

	log.Infof("Предоставлена информация о роле с ID: %d", req.ID)
	c.JSON(http.StatusOK, res)
}
