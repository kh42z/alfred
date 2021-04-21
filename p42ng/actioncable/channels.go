package actioncable

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Identifier struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}

type ChannelEvent interface {
	 OnSubscription(int)
	 OnMessage([]byte, int)
}

func (ac *ActionCable) RegisterChannel(name string, event ChannelEvent) {
	ac.channels[name] = event
}

func (ac *ActionCable) dispatchChannel(event *Event) {
	var i Identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		log.Error("Unable to unmarshal Identifier", i)
		return
	}
	for name, e := range ac.channels {
		if name == i.Channel {
			e.OnMessage(event.Message, i.ID)
		}
	}
}
