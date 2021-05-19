package p42ng

import (
	"alfred/p42ng/api"
	"github.com/gorilla/websocket"
	"github.com/kh42z/actioncable"
	"net/http"
)

type Bot struct {
	Api *api.PongAPI
	Ac  *actioncable.Client
	Ws  *websocket.Conn
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
	ws, err := connectWebsocket(wsHost, b.Api.GenerateAuthHeaders())
	if err != nil {
		return nil, err
	}
	b.Ws = ws
	b.Ac = actioncable.NewClient(b.Ws)
	return &b, nil
}

func connectWebsocket(host string, headers http.Header) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(host+"/cable", headers)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
