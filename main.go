package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	ws *websocket.Conn
	rcvCh chan *Message
	sendCh chan *Message
	pongCh chan bool
}

type Message struct {
	Command    string            `json:"command"`
	Data       string   		 `json:"data,omitempty"`
	Identifier string			 `json:"identifier"`
}

type BearerToken struct {
	Token string `json:"access-token"`
	Client string `json:"client"`
	Uid string `json:"uid"`
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
}

type UserEvent struct {
	ID int `json:"id"`
	Action  string `json:"action"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Action string `json:"action"`
}

func formatChatMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID: ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	m := ChatMessage{Message: "yes", Action: "received"}
	msg, _ := json.Marshal(m)

	return &Message{
		Command: "message",
		Identifier: string(data),
		Data: string(msg),
	}
}

func formatSubscribeMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID: ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	return &Message{
		Command: "subscribe",
		Identifier: string(data),
	}
}

func IdentifyChannel(event *Event, ch chan *Message) {
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
		}
		subscribeOnEvent(ch, &personnalEvent)
	case "ChatChannel":
		log.Info("I received a chat message: ", string(event.Message))
		sendChatResponse(ch)
	default:
		log.Info("Unknown chan")
	}
}


func subscribeOnEvent(ch chan *Message, p *UserEvent) {
	switch p.Action {
	case "game_invitation":
		subscribeGame(ch, p.ID)
	case "chat_invitation":
		subscribeChat(ch, p.ID)
	default:
		log.Info("SubscribeOnEvent: Unknown action")
	}
}


func receiveRoutine(c *Client, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		_, message,  err := c.ws.ReadMessage()
		if err != nil {
			log.Fatal("Unable to rcv:", err)
		}
		var e map[string]interface{}
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
		}
		if t, ok := e["type"]; ok {
			switch t {
			case "welcome":
				log.Infof("Connected to PongWebsocket")
			case "confirm_subscription":
				log.Infof("Subscribed to %s", e["identifier"])
			case "ping":
				c.pongCh <- true
			default:
				log.Info("rcv:",t,  string(message))
			}
		}else{
			log.Info("RawMessage:", string(message) )
			var e Event
			err := json.Unmarshal(message, &e)
			if err != nil {
				log.Error("Unable to unmarshal Event:", err)
			}else{
				IdentifyChannel(&e, c.sendCh)
			}
		}
	}
}

func parseMessage(message []byte){
	var m *Message
	err := json.Unmarshal(message, &m)
	if err != nil {
		log.Error("Unable to unmarshal message", err)
	}
}

func sendRoutine(c *Client, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		select {
		case m := <- c.sendCh:
			log.Info("sent:", m)
			if err := c.ws.WriteJSON(m); err != nil {
				log.Error("Unable to send msg:", err)
			}
		case <- c.pongCh:
			if err := c.ws.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				log.Error("Unable to send ping:", err)
			}
		}
	}
}

func sendChatResponse(msg chan *Message) {
	msg <- formatChatMessage("ChatChannel", 1)
}

func subscribeUser(msg chan *Message,  ID int) {
	log.Debug("Subscribing to UserChannel")
	msg <- formatSubscribeMessage("UserChannel", ID)
}

func subscribeGame(msg chan *Message,  ID int) {
	log.Info("Subscribing to Game topic")
	msg <- formatSubscribeMessage("GameChannel", ID)
}

func subscribeChat(msg chan *Message,  ID int) {
	log.Info("Subscribing to Chat topic")
	msg <- formatSubscribeMessage("ChatChannel", ID)
}

func getBearerToken() *BearerToken {
	var resp *http.Response
	for {
		var err error
		resp, err = http.Get("http://pong:3000/two_factor/1?code=" + os.Getenv("ALFRED_CODE"))
		if err != nil {
			time.Sleep(10 * time.Second)
		}else{
			break
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body: ", err)
	}
	var b *BearerToken
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Fatal("Unable to unmarshal auth:", err)
	}
	return b
}

func connect() *websocket.Conn {
	url := "ws://pong:3000/cable"
	bearerToken := getBearerToken()
	req,_ := http.NewRequest("GET", url, nil)
	req.Header.Add("Origin", "http://localhost:3000/")
	req.Header.Add("access-token", bearerToken.Token)
	req.Header.Add("client", bearerToken.Client)
	req.Header.Add("uid", bearerToken.Uid)
	ws, _, err := websocket.DefaultDialer.Dial(url, req.Header)
	if err != nil {
		log.Error("Unable to connect:", err)
	}
	return ws
}

func configLogger() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		TimestampFormat: "[15:04:05]",
	})
}

func main() {
	wg := sync.WaitGroup{}
	configLogger()
	c := &Client{
		sendCh: make(chan *Message, 20),
		rcvCh: make(chan *Message),
		pongCh: make(chan bool) }
	c.ws = connect()
	defer c.ws.Close()
	go receiveRoutine(c, &wg)
	go sendRoutine(c, &wg)
	time.Sleep(2 * time.Second)
	subscribeUser(c.sendCh, 1)
	subscribeChat(c.sendCh, 1)
	wg.Wait()
}
