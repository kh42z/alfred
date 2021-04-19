package actioncable

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func formatSubscribeMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID:      ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	return &Message{
		Command:    "subscribe",
		Identifier: string(data),
	}
}
func (ac *ActionCable) Subscribe(channel string, ID int) {
	ac.sendCh <- formatSubscribeMessage(channel, ID)
}
