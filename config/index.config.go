package config

import (
	"gin-gonic-gorm/config/app_config"
	"gin-gonic-gorm/config/db_config"
)

func InitConfig() {
	app_config.InitAppConfig()
	db_config.InitDBConfig()
	// log_config.DefaultLogging()
}
