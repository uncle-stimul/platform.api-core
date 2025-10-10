package configs

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	address  string
	port     int
	username string
	password string
	dbname   string
}

func (cfg DatabaseConfig) GetAddress() string {
	return cfg.address
}

func (cfg DatabaseConfig) GetPort() int {
	return cfg.port
}

func (cfg DatabaseConfig) GetUser() string {
	return cfg.username
}

func (cfg DatabaseConfig) GetPass() string {
	return cfg.password
}

func (cfg DatabaseConfig) GetDBName() string {
	return cfg.dbname
}

func NewDatabaseConfig() *DatabaseConfig {
	port, err := strconv.Atoi(os.Getenv("PLATFORM_DB_PORT"))
	if err != nil {
		port = 5432
	}
	return &DatabaseConfig{
		address:  os.Getenv("PLATFORM_DB_ADDRESS"),
		port:     port,
		dbname:   os.Getenv("PLATFORM_DB_NAME"),
		username: os.Getenv("PLATFORM_DB_USER"),
		password: os.Getenv("PLATFORM_DB_PASS"),
	}
}
