package routes

import (
	"gin-gonic-gorm/config/app_config"
	"gin-gonic-gorm/controllers/auth_controller"
	"gin-gonic-gorm/controllers/book_controller"
	"gin-gonic-gorm/controllers/user_controller"

	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine) {
	route := app.Group("")
	route.Static(app_config.STATIC_ROUTE, app_config.STATIC_DIR)
	// ROUTE USER
	userRoute := route.Group("user")
	{
		userRoute.GET("/", user_controller.GetAllUser)
		userRoute.GET("/paginate", user_controller.GetUserPaginate)
		userRoute.POST("/", user_controller.Store)
		userRoute.GET("/:id", user_controller.GetByID)
		userRoute.PATCH("/:id", user_controller.UpdateByID)
		userRoute.DELETE("/:id", user_controller.DeleteByID)
		userRoute.POST("/login", auth_controller.Login)
	}

	// ROUTE BOOK
	route.GET("/book", book_controller.GetAllBook)

	v1Route(route)
}
