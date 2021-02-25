package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Message struct {
	Command    string            `json:"command"`
	Data       string   `json:"data,omitempty"`
	Identifier string			 `json:"identifier"`
	errc       chan error
}

type Command struct {
	Channel string `json:"channel"`
	UserID int `json:"user_id"`
}

type Event struct {
	Type string `json:"type"`
	Message    json.RawMessage    `json:"message"`
	Data       json.RawMessage    `json:"data"`
	Identifier *Command `json:"identifier"`
}

type Data struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func formatSubscribeMessage(channel string, ID int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		UserID: ID,

	})
	if err != nil {
		log.Fatal("Unable to marshal")
	}
	return &Message{
		Command: "subscribe",
		Identifier: string(data),
	}
}

func formatReceivedMessage(text string) *Message {
	d := Data{Message: text, Action: "received"}
	data, err := json.Marshal(d)
	if err != nil {
		log.Println("Unable to marshal")
	}
	id, err := json.Marshal(Command{
		Channel: "ChatChannel",
	})
	if err != nil {
		log.Fatal("Unable to marshal")
	}
	return &Message{
		Command: "message",
		Identifier: string(id),
		Data: string(data),
	}
}

func receiveRoutine(ws *websocket.Conn) {
	for {
		_, message,  err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		var e map[string]interface{}
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
		}
		if t, ok := e["type"]; ok {
			switch t {
			case "welcome":
				log.Infof("Successfuly connected to PongWebsocket")
			case "confirm_subscription":
				log.Infof("Successfuly subscribe to UserChannel")
			case "ping":
				break
			default:
				log.Info("rcv:",t,  string(message))

			}
		}else{
			log.Error("No type key:", string(message) )
		}
	}
}

func sendRoutine(ws *websocket.Conn, msg chan *Message) {
	for {
		m := <- msg
		log.Info("sent:", m)
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

func connect(token, client, uid string) *websocket.Conn {
	url := "ws://localhost:3000/cable"

	req,_ := http.NewRequest("GET", url, nil)
	req.Header.Add("Origin", "http://localhost:3000/")
	req.Header.Add("access-token", token)
	req.Header.Add("client", client)
	req.Header.Add("uid", uid)
	ws, _, err := websocket.DefaultDialer.Dial(url, req.Header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return ws
}



func main() {

	var token, uid, client string

	flag.StringVar(&token, "t", "", "access-token")
	flag.StringVar(&client, "c", "", "client")
	flag.StringVar(&uid, "u","69891", "uid")
	flag.Parse()

	ws := connect(token, client, uid)
	defer ws.Close()
	sendCh := make(chan *Message)
	go receiveRoutine(ws)
	go sendRoutine(ws, sendCh)
	time.Sleep(2 * time.Second)
	subscribeUser(sendCh, 1)
	time.Sleep(100 * time.Second)
}
