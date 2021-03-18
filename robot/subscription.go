package robot

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) SubscribeUser(ID int) {
	log.Debug("Subscribing to UserChannel")
	b.sendCh <- formatSubscribeMessage("UserChannel", ID)
}

func (b *Bot) SubscribeGame(ID int) {
	log.Info("I joined game_id [", ID, "]")
	b.sendCh <- formatSubscribeMessage("GameChannel", ID)
}

func (b *Bot) SubscribeActivity() {
	log.Info("Subscribing to ActivtyChannel")
	b.sendCh <- formatSubscribeMessage("ActivityChannel", 1)
}

func (b *Bot) SubscribeChat(ID int) {
	log.Debug("Subscribing to a ChatRoom")
	b.sendCh <- formatSubscribeMessage("ChatChannel", ID)
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

func (b *Bot) InitChatSubscriptions() {
	for _, chat := range b.retrieveSubscriptions() {
		b.SubscribeChat(chat.ID)
	}
}
