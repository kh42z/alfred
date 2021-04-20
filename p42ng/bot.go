package p42ng

import (
	"alfred/p42ng/actioncable"
	"alfred/p42ng/api"
)

type Bot struct {
	Api     *api.PongAPI
	Ac *actioncable.ActionCable
}


func NewBot(host, code string, uid int, secure bool) *Bot {
	var wsHost, httpHost string
	if secure {
		httpHost = "https://" + host
		wsHost = "wss://" + host
	}else{
		httpHost = "http://" + host
		wsHost = "ws://" + host
	}
	b := Bot{Api: api.NewAPI(httpHost, code, uid)}
	b.Ac = actioncable.NewActionCable(wsHost, b.Api.GenerateAuthHeaders())
	return &b
}