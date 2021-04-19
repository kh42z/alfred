package p42ng

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Nickname string `json:"nickname"`
}

func (b *Bot) UpdateNickname(name string) {
	target := fmt.Sprintf("/users/%d", b.Api.UserID)
	payload := fmt.Sprintf("{\"nickname\":\"%s\"}", name)
	b.Api.DoPatch(target, payload)
}

func (b *Bot) RetrieveNickname(id int) string {
	target := fmt.Sprintf("/users/%d", id)
	resp, err := b.Api.DoGet(target)
	if err != nil {
		log.Warn("Unable to retrieve user nickname", err)
	}
	u := User{}
	err = json.Unmarshal(resp, &u)
	if err != nil {
		log.Error("Unable to unmarshal:", err)
		return "Unknown"
	}
	return u.Nickname
}
