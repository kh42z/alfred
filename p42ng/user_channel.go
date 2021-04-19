package p42ng

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type UserEvent struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}

func (b *Bot) UserNotification(e []byte, _ int) {
	log.Debug("I received a personnal event!")
	var personalEvent UserEvent
	err := json.Unmarshal(e, &personalEvent)
	if err != nil {
		log.Error("Unable to unmarshal userchannel:", err)
		return
	}
	b.subscribeOnEvent(&personalEvent)
}
