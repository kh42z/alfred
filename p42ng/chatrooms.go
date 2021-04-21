package p42ng

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	ID int `json:"id"`
	Privacy string `json:"privacy"`
	ParticipantIds []int `json:"participant_ids"`
	OwnerId	int `json:"owner_id"`
	Name string `json:"name"`
}

func (b *Bot) GetChatRooms() []*Chat {
	body, err := b.Api.DoGet(fmt.Sprintf("/chats?participant_id=%d", b.Api.UserID))
	if err != nil {
		log.Fatal("Unable to retrieve chatrooms subscriptions", err)
	}
	var chats []*Chat
	err = json.Unmarshal(body, &chats)
	if err != nil {
		log.Error("Unable to unmarshal:", err)
		return nil
	}
	return chats
}

func (b *Bot) CreateDMChatroom(invitedId int) (int, error) {
	newChat := Chat{Privacy: "direct_message",
		ParticipantIds: []int{invitedId},
		OwnerId: b.Api.UserID,
	}
	payload, _ := json.Marshal(newChat)
	resp , err := b.Api.DoPost(string(payload), "/chats")
	if err != nil {
		log.Warn("Unable to create chatroom", err)
		return 0, err
	}
	var chat Chat
	err = json.Unmarshal(resp, &chat)
	if err != nil {
		log.Warn("Unable to unmarshal chat", err)
		return 0, err
	}
	return chat.ID, nil
}