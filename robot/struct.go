package robot

import "encoding/json"

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

type UserEvent struct {
	ID int `json:"id"`
	Action  string `json:"action"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Action string `json:"action"`
}

type MessageContent struct {
	Content string `json:"content"`
	SenderID int `json:"sender_id"`
}
