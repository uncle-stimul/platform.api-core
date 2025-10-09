package configs

import (
	"os"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	address string
	port    string
	mode    string
}

func (cfg ServerConfig) GetAddress() string {
	return cfg.address
}

func (cfg ServerConfig) GetPort() string {
	return cfg.port
}

func (cfg ServerConfig) GetMode() string {
	if cfg.mode == "prod" {
		return gin.ReleaseMode
	}
	return gin.DebugMode
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		address: os.Getenv("PLATFORM_API_ADDRESS"),
		port:    os.Getenv("PLATFORM_API_PORT"),
		mode:    os.Getenv("PLATFORM_API_MODE"),
	}
}
