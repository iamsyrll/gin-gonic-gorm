package log_config

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var defaultLogFilePath = "logs/gin.log"

func createLogFolderIfNotExist(path string) {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("Creating ", dir, " directory")
		err := os.MkdirAll(dir, 0o644)
		if err != nil {
			log.Println("Fail to create ", dir)
		} else {
			log.Println(dir, "directory created")
		}
	}
}

func openOrCreateLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		var errCreateFile error
		logFile, errCreateFile = os.Create(path)
		if errCreateFile != nil {
			log.Println("Can't create log file ", errCreateFile)
		}
	}

	return logFile, nil
}

func DefaultLogging(path ...string) {
	gin.DisableConsoleColor()

	if len(path) > 0 && path[0] != "" {
		defaultLogFilePath = path[0]
	}

	createLogFolderIfNotExist(defaultLogFilePath)
	file, _ := openOrCreateLogFile(defaultLogFilePath)

	gin.DefaultWriter = io.MultiWriter(file)
	log.SetOutput(gin.DefaultWriter)
}
