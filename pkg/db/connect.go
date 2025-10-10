package db

import (
	"fmt"
	"time"

	"platform.api-core/pkg/configs"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connect(cfg configs.DatabaseConfig, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.GetAddress(),
		cfg.GetUser(),
		cfg.GetPass(),
		cfg.GetDBName(),
		cfg.GetPort(),
	)

	pgdb, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)

	if err != nil {
		log.WithError(err).Fatal("Во время подключения к БД возникла ошибка:")
	}

	pollDB, err := pgdb.DB()
	if err != nil {
		log.WithError(err).Fatal("Ошибка создания пула соединений")
	}

	pollDB.SetMaxIdleConns(10)
	pollDB.SetMaxOpenConns(100)
	pollDB.SetConnMaxLifetime(time.Hour)

	return pgdb
}
