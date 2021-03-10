package alfred

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Bot struct {
	ws *websocket.Conn
	rcvCh chan *Message
	sendCh chan *Message
	pongCh chan bool
	bearerToken *BearerToken
	wg *sync.WaitGroup
}


func (b *Bot) receiveRoutine() {
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
			b.identifyChannel(&e)
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
		log.Info("I received a personnal event!")
		var personnalEvent UserEvent
		err := json.Unmarshal(event.Message, &personnalEvent)
		if err != nil {
			log.Error("Unable to unmarshal userchannel:", err)
			return
		}
		b.subscribeOnEvent(&personnalEvent)
	case "ChatChannel":
		var content MessageContent
		err := json.Unmarshal(event.Message, &content)
		if err != nil {
			log.Error("Unable to unmarshal content", err)
			return
		}
		log.Infof("I received a chatMessage > user_%d: [%s]", content.SenderID, content.Content)
		if content.SenderID != 1 {
			sendChatResponse(b.sendCh, i.ID)
		}
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

func NewBot() *Bot {
	return &Bot{
		sendCh: make(chan *Message, 20),
		rcvCh: make(chan *Message),
		pongCh: make(chan bool),
		wg: &sync.WaitGroup{},
	}
}

func (b *Bot) Start(url string, code string) {
	b.connect(url, code)
	go b.receiveRoutine()
	go b.sendRoutine()
}

func (b *Bot) Wait() {
	b.wg.Wait()
	b.ws.Close()
}