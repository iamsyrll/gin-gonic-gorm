package middleware

import (
	"net/http"
	"strings"

	"gin-gonic-gorm/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")

	if !strings.Contains(bearerToken, "Bearer") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Token",
		})
		return
	}

	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Token",
		})
		return
	}
	token := parts[1]

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	claims, err := utils.DecodeToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	ctx.Set("claimsData", claims)
	ctx.Set("user_id", claims["id"])
	ctx.Set("user_email", claims["email"])
	ctx.Set("user_name", claims["name"])

	ctx.Next()
}

func TokenMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	if token != "123" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token not valid",
		})
		return
	}

	ctx.Next()
}
