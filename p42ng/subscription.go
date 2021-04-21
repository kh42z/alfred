package p42ng

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) subscribeOnEvent(p *UserMessage) {
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
