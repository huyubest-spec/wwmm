package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	UploadDir  string
	StaticPath string
	ChainDiff  int
}

var App = &Config{
	DBHost:     getEnv("DB_HOST", "localhost"),
	DBPort:     getEnvInt("DB_PORT", 3306),
	DBUser:     getEnv("DB_USER", "root"),
	DBPassword: getEnv("DB_PASSWORD", "123456"),
	DBName:     getEnv("DB_NAME", "wwmm_db"),
	ServerPort: getEnv("SERVER_PORT", "8080"),
	UploadDir:  getEnv("UPLOAD_DIR", "./uploads"),
	StaticPath: getEnv("STATIC_PATH", "/static"),
	ChainDiff:  getEnvInt("CHAIN_DIFF", 4),
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
