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

type ActivityEvent struct {
	users map[int]string
	b *Bot
}

func (b *Bot) NewActivityEvent() *ActivityEvent {
	return &ActivityEvent{
		users: make(map[int]string),
		b: b,
	}
}

func (a *ActivityEvent) OnSubscription(_ int) {
	log.Infof("I'm listening to Activity!")
}

func (a *ActivityEvent) OnMessage(e []byte, _ int){
	var activityMessage ActivityMessage
	err := json.Unmarshal(e, &activityMessage)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if _, ok := a.users[activityMessage.ID]; !ok {
		a.users[activityMessage.ID] = a.b.RetrieveNickname(activityMessage.ID)
	}
	if activityMessage.Action == "user_update_status" {
		log.Infof("Seems like [%s] status changed to <%s>", a.users[activityMessage.ID], activityMessage.Status)
	}
	//if activityMessage.Action == "user_update_status" && activityMessage.Status == "online" && activityMessage.ID != b.Api.UserID {
	//	chatID, err := b.CreateDMChatroom(activityMessage.ID)
	//	if err == nil {
	//		b.SubscribeChat(chatID)
	//		b.Ac.SendMessage("ChatChannel", chatID, fmt.Sprintf("Have you heard about Pong, %s?", b.users[activityMessage.ID]))
	//	}
	//}
}

func (b *Bot) SubscribeActivity() {
	log.Debug("Subscribing to ActivtyChannel")
	b.Ac.RegisterChannel("ActivityChannel", b.NewActivityEvent())
	b.Ac.Subscribe("ActivityChannel", 1)
}
