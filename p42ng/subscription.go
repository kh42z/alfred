package p42ng

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) SubscribeUser(ID int) {
	log.Debug("Subscribing to UserChannel")
	b.Ac.RegisterChannel("UserChannel", b.UserNotification)
	b.Ac.Subscribe("UserChannel", ID)
}

func (b *Bot) SubscribeGame(ID int) {
	log.Debug("I joined game_id [", ID, "]")
	b.Ac.RegisterChannel("GameChannel", b.GameUpdate)
	b.Ac.Subscribe("GameChannel", ID)
}

func (b *Bot) SubscribeActivity() {
	log.Debug("Subscribing to ActivtyChannel")
	b.Ac.RegisterChannel("ActivityChannel", b.ActivityUpdate)
	b.Ac.Subscribe("ActivityChannel", 1)
}

func (b *Bot) SubscribeChat(ID int) {
	log.Debug("Subscribing to ChatRoom ", ID)
	b.Ac.RegisterChannel("ChatChannel", b.ChatResponse)
	b.Ac.Subscribe("ChatChannel", ID)
}

func (b *Bot) subscribeOnEvent(p *UserEvent) {
	switch p.Action {
	case "game_won":
		log.Infof("I just won this game [%d]", p.ID)
	case "game_lost":
		log.Infof("I just lost a game [%d]", p.ID)
	case "game_invitation":
		b.SubscribeGame(p.ID)
	case "chat_invitation":
		b.SubscribeChat(p.ID)
	case "guild_invitation":
		b.JoinGuild(p.ID)
	default:
		log.Info("SubscribeOnEvent: Unknown action")
	}
}

func (b *Bot) SubscribeToChatRooms() {
	for _, chat := range b.GetChatRooms() {
		b.SubscribeChat(chat.ID)
	}
}
