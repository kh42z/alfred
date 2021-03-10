package alfred

import log "github.com/sirupsen/logrus"


func (b *Bot) SubscribeUser(ID int) {
	log.Debug("Subscribing to UserChannel")
	b.sendCh <- formatSubscribeMessage("UserChannel", ID)
}

func (b *Bot) SubscribeGame(ID int) {
	log.Info("Subscribing to Game topic")
	b.sendCh <- formatSubscribeMessage("GameChannel", ID)
}

func (b *Bot) SubscribeChat(ID int) {
	log.Debug("Subscribing to Chat topic")
	b.sendCh <- formatSubscribeMessage("ChatChannel", ID)
}

func (b *Bot) subscribeOnEvent(p *UserEvent) {
	switch p.Action {
	case "game_invitation":
		b.SubscribeGame(p.ID)
	case "chat_invitation":
		b.SubscribeChat(p.ID)
	default:
		log.Info("SubscribeOnEvent: Unknown action")
	}
}

func (b *Bot) InitChatSubscriptions(host string) {
	for _, chat := range b.retrieveSubscriptions(host) {
		b.SubscribeChat(chat.ID)
	}
}
