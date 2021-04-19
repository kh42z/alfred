package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Unable to cast int:", err)
	}
	return i
}

func toBool(s string) bool {
	i, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Unable to cast bool:", err)
	}
	return i
}
