package p42ng

import (
	"alfred/p42ng/actioncable"
	"alfred/p42ng/api"
)

type Bot struct {
	Api     *api.PongAPI
	Ac *actioncable.ActionCable
	users map[int]string
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
	b := Bot{Api: api.NewAPI(httpHost, code, uid), users: make(map[int]string)}
	b.Ac = actioncable.NewActionCable(wsHost, b.Api.GenerateAuthReq())
	return &b
}