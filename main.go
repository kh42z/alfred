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

type Message struct {
	Command    string            `json:"command"`
	Data       string   `json:"data,omitempty"`
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
	if i.Channel == "UserChannel" {
		log.Info("I received a personnal event!")
		var personnalEvent UserEvent
		err := json.Unmarshal(event.Message, &personnalEvent)
		if err != nil {
			log.Error("Unable to unmarshal userchannel:", err)
		}
		subscribeGame(ch, personnalEvent.ID)
	}
}


func receiveRoutine(ws *websocket.Conn, ch chan *Message) {
	for {
		_, message,  err := ws.ReadMessage()
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
			default:
				log.Info("rcv:",t,  string(message))
			}
		}else{
			log.Debug("RawMessage:", string(message) )
			var e Event
			err := json.Unmarshal(message, &e)
			if err != nil {
				log.Error("Unable to unmarshal Event:", err)
			}else{
				IdentifyChannel(&e, ch)
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

func sendRoutine(ws *websocket.Conn, msg chan *Message) {

	for {
		m := <- msg
		log.Debug("sent:", m)
		if err := ws.WriteJSON(m); err != nil {
			log.Error("Unable to send msg:", err)
		}
	}
}

func subscribeUser(msg chan *Message,  ID int) {
	log.Debug("Subscribing to UserChannel")
	msg <- formatSubscribeMessage("UserChannel", ID)
}

func subscribeGame(msg chan *Message,  ID int) {
	log.Info("Subscribing to Game topic")
	msg <- formatSubscribeMessage("GameChannel", ID)
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
	ws := connect()
	defer ws.Close()
	sendCh := make(chan *Message)
	go receiveRoutine(ws, sendCh)
	go sendRoutine(ws, sendCh)
	time.Sleep(2 * time.Second)
	subscribeUser(sendCh, 1)
	time.Sleep(100 * time.Second)
	wg.Wait()
}
