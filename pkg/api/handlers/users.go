package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func SelectUsers(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var users []models.Users
	if hasErr := utils.ExecuteSelectAll(c, log, pgdb, &users, "Roles", "users"); hasErr {
		return
	}

	result := make([]models.UsersResponse, 0, len(users))
	for _, user := range users {
		roles := make([]string, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
		result = append(result, models.UsersResponse{
			Username: user.Username,
			Roles:    roles,
			Status:   user.Status,
		})
	}

	log.Debug("Предоставлена информационная выборка по объектам из таблицы \"users\"")
	c.JSON(http.StatusOK, result)
}

func SelectUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetUserRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	var user models.Users
	if hasErr := utils.ExecuteSelectByID(c, log, pgdb, &user, "Roles", req.ID, "users"); hasErr {
		return
	}

	res := models.GetUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Status:   user.Status,
		Roles:    utils.ExtractRolesNames(user.Roles),
	}

	log.Debugf("Предоставлена информационная выборка по объекту с \"id:%d\" из таблицы \"users\"", req.ID)
	c.JSON(http.StatusOK, res)
}

func CreateUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.CreateUserRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckStringSlice(c, log, []string{req.Username, req.Password}); hasErr {
		return
	}

	user := models.Users{
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
	}

	if hasErr := utils.ExecuteCreate(c, log, pgdb, &user, "users"); hasErr {
		return
	}

	if len(req.Roles) > 0 {
		var roles []*models.Roles
		if hasErr := utils.ExecuteWhereByID(c, log, pgdb, req.Roles, &roles); hasErr {
			return
		}
		if hasErr := utils.ExecuteAssociationsSet(c, log, pgdb, &user, "users", "Roles", &roles); hasErr {
			return
		}
	}

	utils.CheckAssociations(c, log, pgdb, &user, user.ID, "users", "Roles")
}

func UpdateUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateUserRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingUser models.Users
		if errMsg := utils.CheckExistedObj(tx, log, &existingUser, req.ID, "users"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		user := models.Users{
			ID:       existingUser.ID,
			Username: req.Username,
			Password: req.Password,
			Status:   req.Status,
		}

		if errMsg := utils.ExecuteSave(tx, log, &user, req.ID, "users"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		if errMsg := utils.ExecuteAssociationsClean(tx, log, &user, "Roles", req.ID, "users"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		var rolesToSet []models.Roles
		if len(req.Roles) > 0 {
			for _, id := range req.Roles {
				rolesToSet = append(rolesToSet, models.Roles{ID: id})
			}
			if errMsg := utils.ExecuteAssociationsFind(tx, log, req.Roles, &rolesToSet, "roles", len(rolesToSet), len(req.Roles)); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			if errMsg := utils.ExecuteAssociationsUpdate(tx, log, &user, "Roles", &rolesToSet, req.ID, "users"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
		} else {
			if errMsg := utils.ExecuteAssociationsPrepareRollback(tx, log, "Roles", &existingUser, req.ID, "users"); errMsg != "" {
				return fmt.Errorf("%s", errMsg)
			}
			rolesToSet = existingUser.Roles
			if len(rolesToSet) > 0 {
				if errMsg := utils.ExecuteAssociationsRollback(tx, log, &user, "Roles", &rolesToSet, req.ID, "users"); errMsg != "" {
					return fmt.Errorf("%s", errMsg)
				}
			}
		}
		if errMsg := utils.CheckUpdateObj(tx, log, "Roles", &user, req.ID, "users"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}
		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("Выполнение обновления объекта с \"id:%d\" в таблице \"users\" привело к ошибке", req.ID)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполненио обновление объекта с \"id:%d\" в таблице \"users\"", req.ID)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func DeleteUser(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.DeleteUserRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	if hasErr := utils.CheckBuildinAdmin(c, log, req.ID); hasErr {
		return
	}

	if hasErr := utils.CheckDeleteObj(c, log, pgdb, &models.Users{}, req.ID, "users"); hasErr {
		return
	}

	utils.ExecuteDeleteByID(c, log, pgdb, &models.Users{ID: req.ID}, req.ID, "users")
}
