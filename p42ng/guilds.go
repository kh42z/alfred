package p42ng

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (b *Bot) JoinGuild(guild_id int) {
	_, err := b.Api.DoPost(`{ "user_id": `+strconv.Itoa(b.Api.UserID)+` }`, "/guilds/"+strconv.Itoa(guild_id)+"/members")
	if err != nil {
		log.Error("Unable to join guild", err)
		return
	}
	resp, err := b.Api.DoGet("/guilds/" + strconv.Itoa(guild_id))

	if err != nil {
		log.Warn("Unable to get Guild name", err)
		return
	}
	var e map[string]interface{}
	json.Unmarshal(resp, &e)
	if name, ok := e["name"]; ok {
		log.Info("I joined the guild [", name, "]")
	}
}
