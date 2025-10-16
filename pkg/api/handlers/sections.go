package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"
)

func SelectSections(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var sections []models.Sections
	if hasErr := utils.ExecuteSelectAll(c, log, pgdb, &sections, "", "sections"); hasErr {
		return
	}

	result := make([]models.GetSectionsResponse, 0, len(sections))
	for _, section := range sections {
		result = append(result, models.GetSectionsResponse{
			Module:   section.Module,
			Endpoint: section.Endpoint,
		})
	}

	log.Debug("Предоставлена информационная выборка по объектам из таблицы \"sections\"")
	c.JSON(http.StatusOK, result)
}

func SelectSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelPermissionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	var section models.Sections
	if hasErr := utils.ExecuteSelectByID(c, log, pgdb, &section, "", req.ID, "sections"); hasErr {
		return
	}

	res := models.GetSectionsResponse{
		Module:   section.Module,
		Endpoint: section.Endpoint,
	}

	log.Debugf("Предоставлена информационная выборка по объекту с \"id:%d\" из таблицы \"sections\"", req.ID)
	c.JSON(http.StatusOK, res)
}

func CreateSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.CreateSectionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckStringSlice(c, log, []string{req.Module, req.Endpoint}); hasErr {
		return
	}

	section := models.Sections{
		Module:   req.Module,
		Endpoint: req.Endpoint,
	}

	if hasErr := utils.ExecuteCreate(c, log, pgdb, &section, "sections"); hasErr {
		return
	} else {
		msg := fmt.Sprintf("Выполнено создание обекта с \"id:%d\" из таблицы \"sections\"", section.ID)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func UpdateSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.UpdateSectionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	err := pgdb.Transaction(func(tx *gorm.DB) error {
		var existingSection models.Sections
		if errMsg := utils.CheckExistedObj(tx, log, &existingSection, req.ID, "sections"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		role := models.Sections{
			ID:       existingSection.ID,
			Module:   req.Module,
			Endpoint: req.Endpoint,
		}

		if errMsg := utils.ExecuteSave(tx, log, &role, req.ID, "sections"); errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}

		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("Выполнение обновления объекта с \"id:%d\" в таблице \"sections\" привело к ошибке", req.ID)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполненио обновление объекта с \"id:%d\" в таблице \"sections\"", req.ID)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func DeleteSection(c *gin.Context) {
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	var req models.GetDelSectionRequest
	if hasErr := utils.CheckReqStruct(c, log, &req); hasErr {
		return
	}

	if hasErr := utils.CheckExistID(c, log, req.ID); hasErr {
		return
	}

	if hasErr := utils.CheckDeleteObj(c, log, pgdb, &models.Sections{}, req.ID, "sections"); hasErr {
		return
	}

	utils.ExecuteDeleteByID(c, log, pgdb, &models.Sections{ID: req.ID}, req.ID, "sections")
}
