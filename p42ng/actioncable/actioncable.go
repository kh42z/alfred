package actioncable

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type ActionCable struct {
	ws       *websocket.Conn
	wg       *sync.WaitGroup
	rcvCh    chan *Message
	sendCh   chan *Message
	pongCh   chan bool
	stopCh	 chan bool
	channels map[string]OnEventFn
}

func NewActionCable(host string, headers http.Header) *ActionCable {
	ac := ActionCable{
		sendCh:   make(chan *Message, 20),
		rcvCh:    make(chan *Message),
		pongCh:   make(chan bool),
		stopCh:   make(chan bool),
		channels: make(map[string]OnEventFn),
		wg:      &sync.WaitGroup{},
	}
	var err error
	ac.ws, _, err = websocket.DefaultDialer.Dial(host+"/cable", headers)
	if err != nil {
		log.Fatal("Unable to connect to websocket:", err)
	}
	return &ac
}

func (ac *ActionCable) Start() {
	go ac.receiveRoutine()
	go ac.sendRoutine()
}

func (ac *ActionCable) Wait() {
	ac.wg.Wait()
}

func (ac *ActionCable) Stop() {
	ac.ws.Close()
}
