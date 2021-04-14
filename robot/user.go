package robot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type UserEvent struct {
	ID int `json:"id"`
	Action  string `json:"action"`
}

func (b *Bot) UserNotification(e []byte) {
	log.Debug("I received a personnal event!")
	var personnalEvent UserEvent
	err := json.Unmarshal(e, &personnalEvent)
	if err != nil {
		log.Error("Unable to unmarshal userchannel:", err)
		return
	}
	b.subscribeOnEvent(&personnalEvent)
}
