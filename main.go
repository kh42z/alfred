package main

import (
	"alfred/robot"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func configLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		TimestampFormat: "[15:04:05]",
	})
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func main() {
	configLogger()
	strUid := getEnv("ALFRED_UID", "1")
	code := getEnv("ALFRED_CODE", "0000")
	uid, err := strconv.Atoi(strUid)
	if err != nil {
		log.Fatal("Unable to cast uid to int:", err)
	}
	hal :=  robot.NewBot("pong:3000", uid)
	hal.Start(code)
	hal.SubscribeUser(uid)
	hal.SubscribeActivity()
	hal.InitChatSubscriptions()
	hal.Wait()
}
