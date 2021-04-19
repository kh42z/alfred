package actioncable

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Identifier struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}

type OnEventFn func([]byte, int)

func (ac *ActionCable) RegisterChannel(name string, fn OnEventFn) {
	ac.channels[name] = fn
}

func (ac *ActionCable) dispatchChannel(event *Event) {
	var i Identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		log.Error("Unable to unmarshal Identifier", i)
		return
	}
	for name, fn := range ac.channels {
		if name == i.Channel {
			fn(event.Message, i.ID)
		}
	}
}
