package p42ng

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

type ChatEvent struct {
	b *Bot
}

func (b *Bot) NewChatEvent() *ChatEvent {
	return &ChatEvent{
		b: b,
	}
}

func (c *ChatEvent) OnSubscription(chatroomID int){
	log.Infof("I joined the chatroom %d.", chatroomID)
}

func (c *ChatEvent) OnMessage(e []byte, chatroomID int){
	var content MessageContent
	err := json.Unmarshal(e, &content)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if content.SenderID == c.b.Api.UserID {
		return
	}
	log.Infof("I received a chatMessage > user_%d: [%s]", content.SenderID, content.Content)
	if len(content.Content) > 2 && content.Content[0] == '!' {
		c.adminCmd(chatroomID, content.Content[1:])
	} else {
		c.b.sendChatResponse(chatroomID, "yes")
	}
}

func (c *ChatEvent) sendUsage(chatRoomID int) {
	c.b.sendChatResponse(chatRoomID, "!kick <user> <duration>\n")
}

func (c *ChatEvent) adminCmd(chatroomID int, cmdline string) {
	ss := strings.Split(cmdline, " ")
	switch ss[0] {
	case "kick":
		if len(ss) < 2 {
			c.sendUsage(chatroomID)
			return
		}
		_, err := c.b.Api.DoPost(fmt.Sprintf("{\"user_id\": %s, \"duration\": %s}", ss[1], ss[2]), fmt.Sprintf("/chats/%d/mutes", chatroomID))
		if err != nil {
			c.b.sendChatResponse(chatroomID, err.Error())
		} else {
			c.b.sendChatResponse(chatroomID, fmt.Sprintf("Well, [%s] is muted for a while, Sir.", ss[1]))
		}
	default:
		c.sendUsage(chatroomID)
	}
}



func (b *Bot) sendChatResponse(id int, message string) {
	m := ChatMessage{Message: message, Action: "received"}
	msg, _ := json.Marshal(m)
	b.Ac.SendMessage("ChatChannel", id, string(msg))
}


