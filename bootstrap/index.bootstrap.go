package bootstrap

import (
	"log"

	"gin-gonic-gorm/config"
	"gin-gonic-gorm/config/app_config"
	"gin-gonic-gorm/config/cors_config"
	"gin-gonic-gorm/config/log_config"
	"gin-gonic-gorm/database"
	"gin-gonic-gorm/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func BoostrapApp() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
	// config init
	log_config.DefaultLogging("logs/file/app.log")
	config.InitConfig()
	// database init
	database.ConnectDatabase()
	// gin engine
	app := gin.Default()
	// CORS
	app.Use(cors_config.CorsConfigContrib())
	// routes init
	routes.InitRoute(app)
	// run server
	app.Run(app_config.PORT)
}
