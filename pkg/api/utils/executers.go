package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"platform.api-core/pkg/models"
)

func ExecuteSelectAll[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, preload string, table string) bool {
	if err := pgdb.Preload(preload).Find(&data).Error; err != nil {
		msg := fmt.Sprintf("При получении выборки объектов из таблицы \"%s\" возникла ошибка", table)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func ExecuteSelectByID[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, preload string, id uint, table string) bool {
	if err := pgdb.Preload(preload).First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Объект с \"id:%d\" отсутствует в таблице \"%s\"", id, table)
			log.Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		} else {
			msg := fmt.Sprintf("При поиске объекта с \"id:%d\" в таблице \"%s\" возникла ошибка", id, table)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		}
		return true
	}
	return false
}

func ExecuteCreate[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, table string) bool {
	if err := pgdb.Create(&data).Error; err != nil {
		msg := "При создании пользователя возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func ExecuteWhereByID[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, ids []uint, data *D) bool {
	if err := pgdb.Where("id IN ?", ids).Find(&data).Error; err != nil {
		msg := "Указанные объекты связи не найдены в базе данных \"platform_core\""
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func ExecuteAssociationsSet[D any, A any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, table string, preloadTable string, preloadData *A) bool {
	if err := pgdb.Model(&data).Association(preloadTable).Append(preloadData); err != nil {
		msg := fmt.Sprintf("При назначении связей между объеками из таблиц \"%s\" и \"%s\" возникла непредвиденная ошибка", table, preloadTable)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func ExecuteAssociationsClean[D any](tx *gorm.DB, log *logrus.Logger, data *D, preload string, id uint, table string) string {
	if err := tx.Model(&data).Association(preload).Clear(); err != nil {
		msg := fmt.Sprintf("При очистке связей для объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteAssociationsFind[A any](tx *gorm.DB, log *logrus.Logger, data []uint, preloadData *A, preloadTable string, dataL int, preloadL int) string {
	if err := tx.Where("id IN ?", data).Find(&preloadData).Error; err != nil || dataL != preloadL {
		msg := fmt.Sprintf("Указанные в JSON-запросе объекты связи отсутствуют в таблице \"%s\"", preloadTable)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteAssociationsUpdate[D any, A any](tx *gorm.DB, log *logrus.Logger, data *D, preloadTable string, preloadData *A, id uint, table string) string {
	if err := tx.Model(&data).Association(preloadTable).Append(preloadData); err != nil {
		msg := fmt.Sprintf("При обновлении связей для объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteAssociationsPrepareRollback[D any](tx *gorm.DB, log *logrus.Logger, preloadTable string, data *D, id uint, table string) string {
	if err := tx.Preload(preloadTable).First(&data, id).Error; err != nil {
		msg := fmt.Sprintf("При загрузке связей для восстановления объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteAssociationsRollback[D any, A any](tx *gorm.DB, log *logrus.Logger, data *D, preloadTable string, preloadData *A, id uint, table string) string {
	if err := tx.Model(&data).Association(preloadTable).Append(preloadData); err != nil {
		msg := fmt.Sprintf("При восстановлении связей для объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteSave[D any](tx *gorm.DB, log *logrus.Logger, data *D, id uint, table string) string {
	if err := tx.Save(&data).Error; err != nil {
		msg := fmt.Sprintf("При обновлении объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteUpdateByID[D any, A any](tx *gorm.DB, log *logrus.Logger, data *D, preloadTable string, preloadData *A, id uint, table string) string {
	if err := tx.Model(&data).Association(preloadTable).Append(preloadData); err != nil {
		msg := fmt.Sprintf("При добавлении связей для обекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func ExecuteDeleteByID[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, id uint, table string) {
	if err := pgdb.Delete(&data).Error; err != nil {
		msg := fmt.Sprintf("При удалении обекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполнено удаление обекта с \"id:%d\" из таблицы \"%s\"", id, table)
		log.Debug(msg)
		c.JSON(http.StatusOK, models.DefaultResponse{Status: "success", Msg: msg})
	}
}
