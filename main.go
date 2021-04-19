package main

import (
	"alfred/p42ng"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func configLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "[15:04:05]",
	})
}

func sigHandler(bot *p42ng.Bot ,sigCh chan os.Signal) {
	sig := <-sigCh
	switch sig {
	case os.Interrupt:
		bot.Ac.Stop()
	}
}


func main() {
	configLogger()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	hal := p42ng.NewBot(
		getEnv("ALFRED_HOST", "pong:3000"),
		getEnv("ALFRED_CODE", "0000"),
		toInt(getEnv("ALFRED_UID", "1")),
		toBool(getEnv("ALFRED_SSL", "false")))
	hal.Ac.Start()
	log.Infof("Alfred at your service")
	go sigHandler(hal, sigCh)
	hal.UpdateNickname("Alfred")
	hal.SubscribeUser(toInt(getEnv("ALFRED_UID", "1")))
	hal.SubscribeActivity()
	hal.SubscribeToChatRooms()
	hal.Ac.Wait()
	log.Info("Goodbye, Sir.")
}
