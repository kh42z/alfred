package p42ng

import (
	"encoding/json"
	"fmt"
	"github.com/kh42z/actioncable"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ChatMessage struct {
	Content string `json:"content"`
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

func (c *ChatEvent) SubscriptionHandler(_ *actioncable.Client, chatroomID int) {
	log.Infof("I joined the chatroom %d.", chatroomID)
}

func (c *ChatEvent) MessageHandler(_ *actioncable.Client, e []byte, chatroomID int) {
	var content MessageContent
	err := json.Unmarshal(e, &content)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if content.SenderID == c.b.Api.UserID {
		return
	}
	log.Infof("[%d] user_%d: [%s]", chatroomID, content.SenderID, content.Content)
	if len(content.Content) > 2 && content.Content[0] == '!' {
		c.adminCmd(chatroomID, content.Content[1:])
	} else {
		c.b.sendChatResponse(chatroomID, "Do you want to play, Sir?")
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
	case "alert":
		if len(ss) < 1 {
			c.sendUsage(chatroomID)
			return
		}
		c.b.sendChatResponse(chatroomID, "&lt")
	default:
		c.sendUsage(chatroomID)
	}
}

func (b *Bot) sendChatResponse(id int, message string) {
	m := ChatMessage{Content: message, Action: "received"}
	msg, _ := json.Marshal(m)
	b.Ac.SendMessage("ChatChannel", id, string(msg))
}

func (b *Bot) SubscribeChat(ID int) {
	log.Debug("Subscribing to ChatRoom ", ID)
	b.Ac.AddChannelHandler("ChatChannel", b.NewChatEvent())
	b.Ac.Subscribe("ChatChannel", ID)
}

func (b *Bot) SubscribeToChatRooms() {
	for _, chat := range b.GetChatRooms() {
		b.SubscribeChat(chat.ID)
	}
}
