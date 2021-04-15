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
		b.sendPaddlePos(channelID, state.Ball.Y)
	}
}

func (b *Bot) sendPaddlePos(channelID int, pos int ) {
	m := GameMessage{Action: "received", Pos: pos}
	msg, _ := json.Marshal(m)
	b.SendMessage("GameChannel", channelID, string(msg))
}

