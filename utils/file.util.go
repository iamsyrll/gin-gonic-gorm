package utils

import (
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func RandomString(n int) string {
	rand.Seed(time.Now().UnixMilli())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func FileValidation(fileHeader *multipart.FileHeader, fileType []string) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	result := false

	log.Println("Content Type File : ", contentType)

	for _, typeFile := range fileType {
		if contentType == typeFile {
			result = true
			break
		}
	}

	return result
}

func FileValidationExtension(fileHeader *multipart.FileHeader, fileExtension []string) bool {
	extenstion := filepath.Ext(fileHeader.Filename)
	log.Println("Extention : ", extenstion)
	result := false

	for _, typeFile := range fileExtension {
		if extenstion == typeFile {
			result = true
			break
		}
	}

	return result
}

func RandomFileName(ext string, prefix ...string) string {
	currentPrefix := "file"
	if len(prefix) > 0 {
		if prefix[0] != "" {
			currentPrefix = prefix[0]
		}
	}

	currentTime := time.Now().UTC().Format("20061206")
	filename := fmt.Sprintf("%s-%s-%s%s", currentPrefix, currentTime, RandomString(5), ext)

	return filename
}

func SaveFile(ctx *gin.Context, fileHeader *multipart.FileHeader, filename string) error {
	err := ctx.SaveUploadedFile(fileHeader, fmt.Sprintf("./public/files/%v", filename))
	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Println("Failed to remove file")
		return err
	}
	return nil
}
