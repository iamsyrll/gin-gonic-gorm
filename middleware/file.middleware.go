package middleware

import (
	"net/http"
	"path/filepath"

	"gin-gonic-gorm/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if fileHeader == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload data",
		})
		return
	}

	validated := utils.FileValidationExtension(fileHeader, []string{".jpg"})
	if !validated {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "File type not allowed",
		})
		return
	}

	extensionFile := filepath.Ext(fileHeader.Filename)
	filename := utils.RandomFileName(extensionFile)

	err = utils.SaveFile(ctx, fileHeader, filename)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save file",
		})
		return
	}

	ctx.Set("filename", filename)

	ctx.Next()
}
