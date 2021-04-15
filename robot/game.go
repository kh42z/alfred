package robot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	Y int `json:"pos"`
	Score int `json:"score"`
}

type Ball struct {
	X int `json:"x"`
	Y int `json:"y"`
	Up bool `json:up`
	Left bool `json:left`
}

type GameMessage struct {
	Pos int `json:"position"`
	Action string `json:"action"`
}

type GameState struct {
	PlayerLeft *Player `json:"player_left"`
	PlayerRight *Player `json:"player_right"`
	Ball *Ball `json:"ball"`
}

func (b *Bot) GameUpdate(e []byte, channelID int) {
	var state GameState
	err := json.Unmarshal(e, &state)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if state.Ball != nil {
		log.Infof("The ball (%d, %d) is moving (up: %t, left: %t)",state.Ball.X, state.Ball.Y, state.Ball.Up, state.Ball.Left)
		sendPaddlePos(b.sendCh, channelID, state.Ball.Y)
	}
}

func sendPaddlePos(msg chan *Message, channelID int, pos int ) {
	msg <- formatGameMessage("GameChannel", channelID, pos)
}

func formatGameMessage(channel string, ID int, pos int) *Message {
	data, err := json.Marshal(Command{
		Channel: channel,
		ID: ID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	m := GameMessage{Action: "received", Pos: pos}
	msg, _ := json.Marshal(m)
	return &Message{
		Command: "message",
		Identifier: string(data),
		Data: string(msg),
	}
}
