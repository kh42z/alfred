package robot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) SendMessage(channelName string, channelID int, content string) {
	data, err := json.Marshal(Command{
		Channel: channelName,
		ID: channelID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	b.sendCh <- &Message{
		Command: "message",
		Identifier: string(data),
		Data: content,
	}
}