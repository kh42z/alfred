package alfred

import "encoding/json"
import log "github.com/sirupsen/logrus"

func formatChatMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID: ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	m := ChatMessage{Message: "yes", Action: "received"}
	msg, _ := json.Marshal(m)

	return &Message{
		Command: "message",
		Identifier: string(data),
		Data: string(msg),
	}
}

func formatSubscribeMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID: ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	return &Message{
		Command: "subscribe",
		Identifier: string(data),
	}
}

