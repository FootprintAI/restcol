package config

import "os"

func SwagHostPath() string {
	return getEnv("RESTCOL_SWAG_HOST", "")
}

func SwagBasePath() string {
	return getEnv("RESTCOL_SWAG_BASEPATH", "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
