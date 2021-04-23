package actioncable

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type ActionCable struct {
	ws       *websocket.Conn
	wg       *sync.WaitGroup
	sendCh   chan *Message
	stopCh   chan bool
	startCh  chan bool
	channels map[string]ChannelEvent
}

func NewActionCable(host string, headers http.Header) (*ActionCable, error) {
	ac := ActionCable{
		sendCh:   make(chan *Message, 20),
		stopCh:   make(chan bool),
		startCh:  make(chan bool),
		channels: make(map[string]ChannelEvent),
		wg:       &sync.WaitGroup{},
	}
	var err error
	ac.ws, _, err = websocket.DefaultDialer.Dial(host+"/cable", headers)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (ac *ActionCable) Start() {
	go ac.receiveRoutine()
	go ac.sendRoutine()
	<-ac.startCh
}

func (ac *ActionCable) Wait() {
	ac.wg.Wait()
}

func (ac *ActionCable) Stop() {
	ac.ws.Close()
}
