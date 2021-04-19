package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	ID int `json:"id"`
}

func (api *PongAPI) GetChatRooms() []*Chat {
	body, err := api.DoGet(fmt.Sprintf("/chats?participant_id=%b", api.UserID))
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
