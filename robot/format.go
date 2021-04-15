package robot

import "encoding/json"
import log "github.com/sirupsen/logrus"

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
