package api

import (
	"fmt"
	"io"

	"platform.api-core/pkg/api/handlers"
	"platform.api-core/pkg/api/middleware"
	"platform.api-core/pkg/configs"
	"platform.api-core/pkg/db"
	"platform.api-core/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Run() {
	logCfg := configs.NewLoggerConfig()
	log := logger.New(logCfg.GetLevel(), logCfg.GetFormat())

	srvCfg := configs.NewServerConfig()
	log.Debug("Инициализация конфигурации компонента [api-core] успешно завершина")

	pgdb := db.Run()
	if pgdb == nil {
		log.Info("Не предвиденная ошибка при подключении к БД")
	}

	gin.SetMode(srvCfg.GetMode())
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	middleware.AddMiddleware(router, log, pgdb)
	handlers.AddRoutes(router)

	err := router.Run(fmt.Sprintf(
		"%s:%s", srvCfg.GetAddress(), srvCfg.GetPort(),
	))
	if err != nil {
		log.WithError(err).Fatal("Не удалось запустить API-сервер: ")
	}

	err = db.Stop(pgdb)
	if err != nil {
		log.WithError(err).Error("При отключении от БД возникла ошибка:")
	}
}
