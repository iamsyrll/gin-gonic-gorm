package user_controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"gin-gonic-gorm/database"
	"gin-gonic-gorm/models"
	"gin-gonic-gorm/requests"
	"gin-gonic-gorm/responses"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUser(ctx *gin.Context) {
	users := new([]models.User)

	err := database.DB.Table("users").Find(&users).Error
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "internal server error",
		})
		log.Println(err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"data": users,
	})
}

func GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user := new(responses.UserResponse)

	err := database.DB.Table("users").Where("id = ? ", id).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "data transmitted",
		"data":    user,
	})
}

func Store(ctx *gin.Context) {
	userRequest := new(requests.UserRequest)

	err := ctx.ShouldBind(&userRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var userExists models.User

	err = database.DB.Table("users").Where("email = ? ", userRequest.Email).First(&userExists).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email has been register",
		})

		return
	}

	user := new(models.User)
	user.Name = &userRequest.Name
	user.Email = &userRequest.Email
	user.Address = &userRequest.Address
	user.BornDate = &userRequest.BornDate

	err = database.DB.Table("users").Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed create user",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "data saved successfully",
		"data":    user,
	})
}

func UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user := new(models.User)
	userReq := new(requests.UserRequest)
	userEmailExist := new(models.User)

	if err := ctx.ShouldBind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := database.DB.Table("users").Where("id = ?", id).Find(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if user.ID == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	err = database.DB.Table("users").Where("email = ? ", userReq.Email).Find(&userEmailExist).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if userEmailExist.Email != nil && *user.ID != *userEmailExist.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email already exists",
		})
		return
	}

	user.Name = &userReq.Name
	user.Email = &userReq.Email
	user.Address = &userReq.Address
	user.BornDate = &userReq.BornDate

	err = database.DB.Table("users").Where("id = ?", id).Updates(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	userResponse := responses.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "data updated successfully",
		"data":    userResponse,
	})
}

func DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user := new(models.User)

	err := database.DB.Table("users").Where("id = ? ", id).Find(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if user.ID == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	err = database.DB.Table("users").Where("id = ? ", id).Unscoped().Delete(&models.User{}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "data deleted successfully",
	})
}

func GetUserPaginate(ctx *gin.Context) {
	page := ctx.Query("page")

	if page == "" {
		page = "1"
	}

	perPage := ctx.Query("perPage")
	if perPage == "" {
		page = "10"
	}

	perPageInt, _ := strconv.Atoi(perPage)
	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}

	users := new([]models.User)

	err := database.DB.Table("users").Offset((pageInt - 1) * perPageInt).Limit(perPageInt).Find(&users).Error
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "internal server error",
		})
		log.Println(err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"data":    users,
		"page":    pageInt,
		"perPage": perPageInt,
	})
}
