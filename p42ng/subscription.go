package p42ng

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) OnUserChannelEvent(p *UserMessage) {
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
	case "achievement_unlocked":
		log.Infof("Seems like I have unlocked an achievement [%d]", p.ID)
	default:
		log.Infof("I do not know how to react to this [%s]", p.Action)
	}
}
