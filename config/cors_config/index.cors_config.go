package cors_config

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var origins = []string{
	"https://domain-saya.com",
	"https://sub.domain-saya.com",
}

func CorsConfig(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Controll-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Controll-Allow-Credential", "true")
	ctx.Writer.Header().Set("Access-Controll-Allow-Headers", "Content-Type, Content-Length")
	ctx.Writer.Header().Set("Access-Controll-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}

func CorsConfigContrib() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = origins

	return cors.New(config)
}
