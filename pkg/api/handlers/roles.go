package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func SelectRoles(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var roles []models.Roles
	if hasErr := utils.ExecuteSelectAll(c, log, pgdb, &roles, "Permissions", "roles"); hasErr {
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

	log.Debug("Предоставлена информационная выборка по объектам из таблицы \"roles\"")
	c.JSON(http.StatusOK, result)
}

func SelectRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetRoleRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	var role models.Roles
	if hasErr := utils.ExecuteSelectByID(c, log, pgdb, &role, "Permissions", req.ID, "roles"); hasErr {
		return
	}

	res := models.GetRoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: utils.ExtractPermissionsNames(role.Permissions),
	}

	log.Debugf("Предоставлена информационная выборка по объекту с \"id:%d\" из таблицы \"roles\"", req.ID)
	c.JSON(http.StatusOK, res)
}

func CreateRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.CreateRoleRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckStringSlice(c, log, []string{req.Name}); hasErr {
		return
	}

	role := models.Roles{
		Name:        req.Name,
		Description: req.Description,
	}

	if hasErr := utils.ExecuteCreate(c, log, pgdb, &role, "roles"); hasErr {
		return
	}

	if len(req.Permissions) > 0 {
		var permissions []*models.Permissions
		if hasErr := utils.ExecuteWhereByID(c, log, pgdb, req.Permissions, &permissions); hasErr {
			return
		}
		if hasErr := utils.ExecuteAssociationsSet(c, log, pgdb, &role, "roles", "Permissions", &permissions); hasErr {
			return
		}
	}

	utils.CheckAssociations(c, log, pgdb, &role, role.ID, "roles", "Permissions")
}

func UpdateRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateRoleRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingRole models.Roles
		if errMsg := utils.CheckExistedObj(tx, log, &existingRole, req.ID, "roles"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		role := models.Roles{
			ID:          existingRole.ID,
			Name:        req.Name,
			Description: req.Description,
		}

		if errMsg := utils.ExecuteSave(tx, log, &role, req.ID, "roles"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		if errMsg := utils.ExecuteAssociationsClean(tx, log, &role, "Permissions", req.ID, "roles"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		var permissionsToSet []models.Permissions
		if len(req.Permissions) > 0 {
			for _, id := range req.Permissions {
				permissionsToSet = append(permissionsToSet, models.Permissions{ID: id})
			}
			if errMsg := utils.ExecuteAssociationsFind(
				tx, log, req.Permissions, &permissionsToSet,
				"permissions", len(permissionsToSet), len(req.Permissions)); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			if errMsg := utils.ExecuteAssociationsUpdate(tx, log, &role, "Permissions", &permissionsToSet, req.ID, "roles"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
		} else {
			if errMsg := utils.ExecuteAssociationsPrepareRollback(tx, log, "Permissions", &existingRole, req.ID, "roles"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			permissionsToSet = existingRole.Permissions
			if len(permissionsToSet) > 0 {
				if errMsg := utils.ExecuteAssociationsRollback(tx, log, &role, "Permissions", &permissionsToSet, req.ID, "roles"); errMsg != "" {
					return fmt.Errorf("%s", errMsg)
				}
			}
		}
		if errMsg := utils.CheckUpdateObj(tx, log, "Permissions", &role, req.ID, "roles"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}
		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("Выполнение обновления объекта с \"id:%d\" в таблице \"roles\" привело к ошибке", req.ID)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполненио обновление объекта с \"id:%d\" в таблице \"roles\"", req.ID)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func DeleteRole(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetRoleRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	if hasErr := utils.CheckDeleteObj(c, log, pgdb, &models.Roles{}, req.ID, "roles"); hasErr {
		return
	}

	utils.ExecuteDeleteByID(c, log, pgdb, &models.Roles{ID: req.ID}, req.ID, "roles")
}
