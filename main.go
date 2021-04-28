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

func sigHandler(bot *p42ng.Bot, sigCh chan os.Signal) {
	sig := <-sigCh
	switch sig {
	case os.Interrupt:
		if err := bot.Ws.Close(); err != nil {
			log.Error("Unable to close websocket ", err)
		}
	}
}

func main() {
	configLogger()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	hal, err := p42ng.NewBot(
		getEnv("ALFRED_HOST", "pong:3000"),
		getEnv("ALFRED_CODE", "0000"),
		toInt(getEnv("ALFRED_UID", "1")),
		toBool(getEnv("ALFRED_SSL", "false")))
	if err != nil {
		log.Fatal("Unable to create bot: ", err)
	}
	go sigHandler(hal, sigCh)
	log.Infof("Alfred at your service")
	hal.UpdateNickname("Alfred")
	go func() {
		hal.SubscribeUser(60)
		hal.SubscribeUser(toInt(getEnv("ALFRED_UID", "1")))
		hal.SubscribeActivity()
		hal.SubscribeToChatRooms()
	}()
	if err := hal.Ac.Run(); err != nil {
		log.Error(err)
	}
	log.Info("Goodbye, Sir.")
}
