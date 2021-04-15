package robot

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ChatMessage struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

type MessageContent struct {
	Content  string `json:"content"`
	SenderID int    `json:"sender_id"`
}

func (b *Bot) sendUsage(chatRoomID int) {
	b.sendChatResponse(chatRoomID, "!kick <user> <duration>\n")
}

func (b *Bot) adminCmd(chatroomID int, cmdline string) {
	ss := strings.Split(cmdline, " ")
	switch ss[0] {
	case "kick":
		if len(ss) < 2 {
			b.sendUsage(chatroomID)
			return
		}
		_, err := b.api.DoPost(fmt.Sprintf("{\"user_id\": %s, \"duration\": %s}", ss[1], ss[2]), fmt.Sprintf("/chats/%d/mutes", chatroomID))
		if err != nil {
			b.sendChatResponse(chatroomID, err.Error())
		} else {
			b.sendChatResponse(chatroomID, fmt.Sprintf("Well, [%s] is muted for a while, Sir.", ss[1]))
		}
	default:
		b.sendUsage(chatroomID)
	}
}

func (b *Bot) ChatResponse(e []byte, chatroomID int) {
	var content MessageContent
	err := json.Unmarshal(e, &content)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if content.SenderID == b.api.UserID {
		return
	}
	log.Infof("I received a chatMessage > user_%d: [%s]", content.SenderID, content.Content)
	if len(content.Content) > 2 && content.Content[0] == '!' {
		b.adminCmd(chatroomID, content.Content[1:])
	} else {
		b.sendChatResponse(chatroomID, "yes")
	}
}

func (b *Bot) sendChatResponse(id int, message string) {
	m := ChatMessage{Message: message, Action: "received"}
	msg, _ := json.Marshal(m)
	b.SendMessage("ChatChannel", id, string(msg))
}
