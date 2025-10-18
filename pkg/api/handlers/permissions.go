package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func SelectPermissions(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)
	var permissions []models.Permissions
	if hasErr := utils.ExecuteSelectAll(c, log, pgdb, &permissions, "Sections", "permissions"); hasErr {
		return
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
	log.Debug("Предоставлена информационная выборка по объектам из таблицы \"permissions\"")
	c.JSON(http.StatusOK, result)
}

func SelectPermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	var permission models.Permissions
	if hasErr := utils.ExecuteSelectByID(c, log, pgdb, &permission, "Sections", req.ID, "permissions"); hasErr {
		return
	}

	res := models.GetRoleResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		Permissions: utils.ExtractSectionsEndpoints(permission.Sections),
	}

	log.Debugf("Предоставлена информационная выборка по объекту с \"id:%d\" из таблицы \"permissions\"", req.ID)
	c.JSON(http.StatusOK, res)
}

func CreatePermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.CreatePermissionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckStringSlice(c, log, []string{req.Name}); hasErr {
		return
	}

	permission := models.Permissions{
		Name:        req.Name,
		Description: req.Description,
	}

	if hasErr := utils.ExecuteCreate(c, log, pgdb, &permission, "permissions"); hasErr {
		return
	}

	if len(req.Sections) > 0 {
		var sections []*models.Sections
		if hasErr := utils.ExecuteWhereByID(c, log, pgdb, req.Sections, &sections); hasErr {
			return
		}
		if hasErr := utils.ExecuteAssociationsSet(c, log, pgdb, &permission, "permissions", "Sections", &sections); hasErr {
			return
		}
	}

	utils.CheckAssociations(c, log, pgdb, &permission, permission.ID, "permissions", "Sections")
}

func UpdatePermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdatePermissionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingPermission models.Permissions
		if errMsg := utils.CheckExistedObj(tx, log, &existingPermission, req.ID, "permissions"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		permission := models.Permissions{
			ID:          existingPermission.ID,
			Name:        req.Name,
			Description: req.Description,
		}

		if errMsg := utils.ExecuteSave(tx, log, &permission, req.ID, "permissions"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		if errMsg := utils.ExecuteAssociationsClean(tx, log, &permission, "Sections", req.ID, "permissions"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		var sectionsToSet []models.Sections
		if len(req.Sections) > 0 {
			for _, id := range req.Sections {
				sectionsToSet = append(sectionsToSet, models.Sections{ID: id})
			}
			if errMsg := utils.ExecuteAssociationsFind(tx, log, req.Sections, &sectionsToSet, "permissions", len(sectionsToSet), len(req.Sections)); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			if errMsg := utils.ExecuteAssociationsUpdate(tx, log, &permission, "Sections", &sectionsToSet, req.ID, "permissions"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
		} else {
			if errMsg := utils.ExecuteAssociationsPrepareRollback(tx, log, "Sections", &existingPermission, req.ID, "permissions"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			sectionsToSet = existingPermission.Sections
			if len(sectionsToSet) > 0 {
				if errMsg := utils.ExecuteAssociationsRollback(tx, log, &permission, "Sections", &sectionsToSet, req.ID, "permissions"); errMsg != "" {
					return fmt.Errorf("%s", errMsg)
				}
			}
		}
		if errMsg := utils.CheckUpdateObj(tx, log, "Sections", &permission, req.ID, "permissions"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}
		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("Выполнение обновления объекта с \"id:%d\" в таблице \"permissions\" привело к ошибке", req.ID)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполненио обновление объекта с \"id:%d\" в таблице \"permissions\"", req.ID)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func DeletePermission(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	if hasErr := utils.CheckDeleteObj(c, log, pgdb, &models.Permissions{}, req.ID, "permissions"); hasErr {
		return
	}

	utils.ExecuteDeleteByID(c, log, pgdb, &models.Permissions{ID: req.ID}, req.ID, "permissions")
}
