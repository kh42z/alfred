package p42ng

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ActivityMessage struct {
	Action string `json:"action"`
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func (b *Bot) ActivityUpdate(e []byte, _ int) {
	var activityMessage ActivityMessage
	err := json.Unmarshal(e, &activityMessage)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if _, ok := b.users[activityMessage.ID]; !ok {
		b.users[activityMessage.ID] = b.RetrieveNickname(activityMessage.ID)
	}
	if activityMessage.Action == "user_update_status" {
		log.Infof("Seems like [%s] status changed to <%s>", b.users[activityMessage.ID], activityMessage.Status)
	}
	//if activityMessage.Action == "user_update_status" && activityMessage.Status == "online" && activityMessage.ID != b.Api.UserID {
	//	chatID, err := b.CreateDMChatroom(activityMessage.ID)
	//	if err == nil {
	//		b.SubscribeChat(chatID)
	//		b.Ac.SendMessage("ChatChannel", chatID, fmt.Sprintf("Have you heard about Pong, %s?", b.users[activityMessage.ID]))
	//	}
	//}
}
