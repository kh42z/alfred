package p42ng

import (
	"alfred/p42ng/actioncable"
	"alfred/p42ng/api"
)

type Bot struct {
	Api *api.PongAPI
	Ac  *actioncable.ActionCable
}

func NewBot(host, code string, uid int, secure bool) (*Bot, error) {
	var wsHost, httpHost string
	if secure {
		httpHost = "https://" + host
		wsHost = "wss://" + host
	} else {
		httpHost = "http://" + host
		wsHost = "ws://" + host
	}
	b := Bot{Api: api.NewAPI(httpHost, code, uid)}
	var err error
	b.Ac, err = actioncable.NewActionCable(wsHost, b.Api.GenerateAuthHeaders())
	if err != nil {
		return nil, err
	}
	return &b, nil
}
