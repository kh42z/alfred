package p42ng

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type UserMessage struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}

type UserEvent struct {
	b *Bot
}

func (b *Bot) NewUserEvent() *UserEvent {
	return &UserEvent{
		b: b,
	}
}

func (u *UserEvent) OnSubscription(_ int) {
	log.Info("I'm listening to my personnal event!")
}

func (u *UserEvent) OnMessage(e []byte, _ int) {
	log.Debug("I received a personnal event!")
	var personalEvent UserMessage
	err := json.Unmarshal(e, &personalEvent)
	if err != nil {
		log.Error("Unable to unmarshal userchannel:", err)
		return
	}
	u.b.subscribeOnEvent(&personalEvent)
}
