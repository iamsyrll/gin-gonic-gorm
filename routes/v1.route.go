package routes

import (
	"gin-gonic-gorm/controllers/file_controller"
	"gin-gonic-gorm/middleware"

	"github.com/gin-gonic/gin"
)

func v1Route(app *gin.RouterGroup) {
	// ROUTE FILE
	authRoute := app.Group("file", middleware.AuthMiddleware)
	authRoute.POST("/upload_file", file_controller.HandleUploadFile)
	authRoute.POST("/middleware", middleware.UploadFile, file_controller.SendStatus)
	authRoute.DELETE("/:filename", file_controller.HandleRemoveFile)
}
