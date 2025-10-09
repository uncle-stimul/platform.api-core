package configs

import (
	"os"
)

type LoggerConfig struct {
	level  string
	format string
}

func (cfg LoggerConfig) GetLevel() string {
	return cfg.level
}

func (cfg LoggerConfig) GetFormat() string {
	return cfg.format
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		level:  os.Getenv("PLATFORM_LOG_LEVEL"),
		format: os.Getenv("PLATFORM_LOG_FORMAT"),
	}
}
