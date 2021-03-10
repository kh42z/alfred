package main

import (
	"alfred/alfred"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
)

func configLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		TimestampFormat: "[15:04:05]",
	})
}

func main() {
	configLogger()
	hal :=  alfred.NewBot()
	hal.Start("pong:3000", os.Getenv("ALFRED_CODE"))
	hal.SubscribeUser(1)
	hal.InitChatSubscriptions("pong:3000")
	hal.Wait()
}
