package robot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ChatMessage struct {
	Message string `json:"message"`
	Action string `json:"action"`
}

type MessageContent struct {
	Content string `json:"content"`
	SenderID int `json:"sender_id"`
}

func (b *Bot) ChatResponse(e []byte, chatroomID int) {
	var content MessageContent
	err := json.Unmarshal(e, &content)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	log.Infof("I received a chatMessage > user_%d: [%s]", content.SenderID, content.Content)
	if content.SenderID != b.api.UserID {
		sendChatResponse(b.sendCh, chatroomID)
	}
}

func sendChatResponse(msg chan *Message, id int) {
	msg <- formatChatMessage("ChatChannel", id)
}

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
