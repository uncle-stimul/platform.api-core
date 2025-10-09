package configs

import "os"

type DatabaseConfig struct {
	address  string
	port     string
	username string
	password string
	dbname   string
}

func (cfg DatabaseConfig) GetAddress() string {
	return cfg.address
}

func (cfg DatabaseConfig) GetPort() string {
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
	return &DatabaseConfig{
		address:  os.Getenv("PLATFORM_DB_ADDRESS"),
		port:     os.Getenv("PLATFORM_DB_PORT"),
		dbname:   os.Getenv("PLATFORM_DB_NAME"),
		username: os.Getenv("PLATFORM_DB_USER"),
		password: os.Getenv("PLATFORM_DB_PASS"),
	}
}
