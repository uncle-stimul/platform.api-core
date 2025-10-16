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

func CheckReqStruct[R any](c *gin.Context, log *logrus.Logger, req *R) bool {
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен некорректный JSON для выполнения запроса"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{Status: "error", Msg: msg})
		c.Abort()
		return true
	}
	return false
}

func CheckExistID(c *gin.Context, log *logrus.Logger, id uint) bool {
	if id == 0 {
		msg := "Полученный JSON не содержит идентификатора объекта"
		log.Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func CheckStringSlice(c *gin.Context, log *logrus.Logger, data []string) bool {
	msg := "Полученный JSON-запрос не содержит обязательные поля"
	for _, item := range data {
		if item == "" {
			log.Error(msg)
			c.JSON(http.StatusBadRequest, models.DefaultResponse{Status: "error", Msg: msg})
			return true
		}
	}
	return false
}

func CheckBuildinAdmin(c *gin.Context, log *logrus.Logger, id uint) bool {
	if id == 1 {
		msg := "Невозможно выполнить действия над встроенным администратором"
		log.Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}

func CheckAssociations[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, id uint, table string, preload string) {
	if err := pgdb.Preload(preload).First(&data, &id).Error; err != nil {
		msg := fmt.Sprintf("При загрузке связей для объекта с \"id:%d\" из таблицы \"%s\" возникла ошибка", id, table)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{Status: "error", Msg: msg})
	} else {
		msg := fmt.Sprintf("Выполнено создание нового объекта с \"id:%d\" в таблице \"%s\"", id, table)
		log.Debug(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{Status: "success", Msg: msg})
	}
}

func CheckExistedObj[D any](tx *gorm.DB, log *logrus.Logger, data *D, id uint, table string) string {
	if err := tx.First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("Объект с \"id:%d\" отсутствует в таблице \"%s\"", id, table)
			log.WithError(err).Error(msg)
			return msg
		}
		msg := fmt.Sprintf("Поиск объекта с \"id:%d\" в таблице \"%s\" завершился ошибкой", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func CheckUpdateObj[D any](tx *gorm.DB, log *logrus.Logger, preloadTable string, data *D, id uint, table string) string {
	if err := tx.Preload(preloadTable).First(&data, id).Error; err != nil {
		msg := fmt.Sprintf("При проверки обновления связей для объекта с \"id:%d\" в таблице \"%s\" завершился ошибкой", id, table)
		log.WithError(err).Error(msg)
		return msg
	}
	return ""
}

func CheckDeleteObj[D any](c *gin.Context, log *logrus.Logger, pgdb *gorm.DB, data *D, id uint, table string) bool {
	if err := pgdb.First(&data, id).Error; err != nil {
		msg := fmt.Sprintf("Поиск объекта с \"id:%d\" в таблице \"%s\" завершился ошибкой", id, table)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusNotFound, models.DefaultResponse{Status: "error", Msg: msg})
		return true
	}
	return false
}
