package file_controller

import (
	"log"
	"net/http"
	"path/filepath"

	"gin-gonic-gorm/constanta"
	"gin-gonic-gorm/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SendStatus(ctx *gin.Context) {
	filename := ctx.MustGet("filename").(string)

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "file uploaded",
		"file_name": filename,
	})
}

func HandleUploadFile(ctx *gin.Context) {
	claimsData := ctx.MustGet("claimsData").(jwt.MapClaims)
	log.Println(claimsData)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if fileHeader == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload data",
		})
		return
	}

	validated := utils.FileValidationExtension(fileHeader, []string{".jpg"})
	if !validated {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "File type not allowed",
		})
		return
	}

	extensionFile := filepath.Ext(fileHeader.Filename)
	filename := utils.RandomFileName(extensionFile)

	err = utils.SaveFile(ctx, fileHeader, filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save file",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file uploaded!",
	})
}

func HandleRemoveFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return
	}

	err := utils.RemoveFile(constanta.DIR_FILE + filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file deleted",
	})
}
