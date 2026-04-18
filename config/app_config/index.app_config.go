package app_config

import "os"

var (
	PORT         = ":8000"
	STATIC_ROUTE = "/public"
	STATIC_DIR   = "./public"
	SECRET_KEY   = "SECRET_KEY"
)

func InitAppConfig() {
	portEnv := os.Getenv("APP_PORT")
	if portEnv != "" {
		PORT = portEnv
	}

	staticRouteEnv := os.Getenv("STATIC_ROUTE")
	if staticRouteEnv != "" {
		STATIC_ROUTE = staticRouteEnv
	}

	staticDirEnv := os.Getenv("STATIC_ROUTE")
	if staticDirEnv != "" {
		STATIC_DIR = staticDirEnv
	}

	secretKeyEnv := os.Getenv("SECRET_KEY")
	if secretKeyEnv != "" {
		STATIC_DIR = secretKeyEnv
	}
}
