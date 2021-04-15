package robot

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Bot struct {
	ws      *websocket.Conn
	rcvCh   chan *Message
	sendCh  chan *Message
	pongCh  chan bool
	statsCh chan bool
	api     *PongAPI
	wg      *sync.WaitGroup
	host    string
}

type Message struct {
	Command    string            `json:"command"`
	Data       string   		 `json:"data,omitempty"`
	Identifier string			 `json:"identifier"`
}

type Command struct {
	Channel string `json:"channel"`
	ID int `json:"id"`
}

type Event struct {
	Message    json.RawMessage    `json:"message"`
	Identifier string `json:"identifier"`
}

type Identifier struct {
	Channel string `json:"channel"`
	ID int `json:"id"`
}

type onMessageFn func(*Event)

func (b *Bot) receiveRoutine(onMessage onMessageFn) {
	b.wg.Add(1)
	for {
		_, message,  err := b.ws.ReadMessage()
		if err != nil {
			log.Fatal("Unable to rcv:", err)
		}
		var e map[string]interface{}
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
			return
		}
		if t, ok := e["type"]; ok {
			switch t {
			case "welcome":
				log.Infof("Connected to PongWebsocket")
			case "confirm_subscription":
				log.Infof("Subscribed to %s", e["identifier"])
			case "ping":
				b.pongCh <- true
			default:
				log.Info("rcv:",t,  string(message))
			}
		}else{
			log.Debug("RawMessage:", string(message) )
			var e Event
			err := json.Unmarshal(message, &e)
			if err != nil {
				log.Error("Unable to unmarshal Event:", err)
				return
			}
			onMessage(&e)
		}
	}
}

func (b *Bot) identifyChannel(event *Event) {
	var i Identifier
	err := json.Unmarshal([]byte(event.Identifier), &i)
	if err != nil {
		log.Error("Unable to unmarshal Identifier", i)
		return
	}
	switch i.Channel {
	case "UserChannel":
		b.UserNotification(event.Message)
	case "ChatChannel":
		b.ChatResponse(event.Message, i.ID)
	case "ActivityChannel":
		ActivityUpdate(event.Message)
	case "GameChannel":
		b.GameUpdate(event.Message, i.ID)
	default:
		log.Info("Unknown chan")
	}
}

func (b *Bot) sendRoutine() {
	b.wg.Add(1)
	for {
		select {
		case m := <- b.sendCh:
			log.Debug("sent:", m)
			if err := b.ws.WriteJSON(m); err != nil {
				log.Error("Unable to send msg:", err)
			}
		case <- b.pongCh:
			if err := b.ws.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				log.Error("Unable to send ping:", err)
			}
		}
	}
}

func NewBot(host string, uid int) *Bot {
	return &Bot{
		host: host,
		sendCh: make(chan *Message, 20),
		rcvCh: make(chan *Message),
		pongCh: make(chan bool),
		statsCh: make(chan bool),
		api: &PongAPI{ UserID: uid},
		wg: &sync.WaitGroup{},
	}
}

func (b *Bot) Start(code string) {
	b.connect(code)
	go b.receiveRoutine(b.identifyChannel)
	go b.sendRoutine()
}

func (b *Bot) Wait() {
	b.wg.Wait()
	b.ws.Close()
}