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

type GameState struct {
	PlayerLeft *Player `json:"player_left"`
	PlayerRight *Player `json:"player_right"`
	Ball *Ball `json:"ball"`
}

func GameUpdate(e []byte) {
	var state GameState
	err := json.Unmarshal(e, &state)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if state.Ball != nil {
		log.Infof("The ball (%d, %d) is moving (up: %t, left: %t)",state.Ball.X, state.Ball.Y, state.Ball.Up, state.Ball.Left)
	}
}
