package handlers

import (
	"fmt"
	"net/http"

	"platform.api-core/pkg/api/utils"
	"platform.api-core/pkg/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	log := utils.AddContextLogger(c)
	pgdb := utils.AddContextDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "Получен не корректрный JSON для создания пользователя"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		msg := "В полученном JSON не указано имя пользователя или пароль"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	user := models.Users{
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
	}

	if err := pgdb.Create(&user).Error; err != nil {
		msg := "При создании пользователя возникла ошибка"
		log.WithError(err).Error(msg)
		c.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	}

	if len(req.Roles) > 0 {
		var roles []*models.Roles
		if err := pgdb.Where("id IN ?", req.Roles).Find(&roles).Error; err != nil {
			msg := "Указанные роли не найдены в базе данных"
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
		if err := pgdb.Model(&user).Association("Roles").Append(roles); err != nil {
			msg := fmt.Sprintf("При назначении ролей пользователю \"%s\" возникла не предвиденная ошибка", user.Username)
			log.WithError(err).Error(msg)
			c.JSON(http.StatusInternalServerError, models.DefaultResponse{
				Status: "error",
				Msg:    msg,
			})
			return
		}
	}

	if err := pgdb.Preload("Roles").First(&user, user.ID).Error; err != nil {
		msg := fmt.Sprintf("При загрузке ролей для пользователя \"%s\" возникла ошибка", user.Username)
		log.WithError(err).Error(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "error",
			Msg:    msg,
		})
		return
	} else {
		msg := fmt.Sprintf("Было выполнено создание пользователя \"%s\"", user.Username)
		log.Info(msg)
		c.JSON(http.StatusCreated, models.DefaultResponse{
			Status: "success",
			Msg:    msg,
		})
	}
}
