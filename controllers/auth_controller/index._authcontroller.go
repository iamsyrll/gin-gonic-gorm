package auth_controller

import (
	"net/http"
	"time"

	"gin-gonic-gorm/database"
	"gin-gonic-gorm/models"
	"gin-gonic-gorm/requests"
	"gin-gonic-gorm/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(ctx *gin.Context) {
	loginRequest := new(requests.LoginRequest)
	if errReq := ctx.ShouldBind(loginRequest); errReq != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"error":   errReq,
		})
		return
	}

	user := new(models.User)

	err := database.DB.Table("users").Where("email = ?", loginRequest.Email).First(user).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email/Passworrd is wrong",
		})
		return
	}

	if loginRequest.Password != "anjaygurinjaymakanbajai" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email/Passworrd is wrong",
		})
		return
	}

	payload := &jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtToken, err := utils.GenerateToken(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"token":   jwtToken,
	})
}
