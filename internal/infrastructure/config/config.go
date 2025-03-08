package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func ReadEnvConfig(envPath string) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ReadIntVal(optName string, defaultVal int) int {
	if optStr, ok := os.LookupEnv(optName); ok {
		if value, err := strconv.Atoi(optStr); err == nil {
			return value
		}
	}
	return defaultVal
}
