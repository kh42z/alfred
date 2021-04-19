package main

import (
	"alfred/p42ng"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func configLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "[15:04:05]",
	})
}

func main() {
	configLogger()
	hal := p42ng.NewBot(
		getEnv("ALFRED_HOST", "pong:3000"),
		getEnv("ALFRED_CODE", "0000"),
		toInt(getEnv("ALFRED_UID", "1")),
		toBool(getEnv("ALFRED_SSL", "false")))
	hal.Ac.Start()
	hal.Api.UpdateNickname("Alfred")
	hal.SubscribeUser(toInt(getEnv("ALFRED_UID", "1")))
	hal.SubscribeActivity()
	hal.SubscribeToChatRooms()
	hal.Ac.Wait()
}
