package repository

import (
	"platform.api-core/pkg/configs"
	"platform.api-core/pkg/logger"

	"gorm.io/gorm"
)

func Run() *gorm.DB {
	logCfg := configs.NewLoggerConfig()
	log := logger.New(logCfg.GetLevel(), logCfg.GetFormat())

	dbCfg := configs.NewDatabaseConfig()
	log.Debug("Инициализация конфигурации базы данных [api-core] успешно завершина")

	pgdb := connect(*dbCfg, log)

	initShema(pgdb, log)

	return pgdb
}

func Stop(pgdb *gorm.DB) error {
	con, err := pgdb.DB()
	if err != nil {
		return err
	}
	return con.Close()
}
