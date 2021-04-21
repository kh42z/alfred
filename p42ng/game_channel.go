package p42ng

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Player struct {
	Y     int `json:"pos"`
	Score int `json:"score"`
}

type Ball struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Up   bool `json:up`
	Left bool `json:left`
}

type GameMessage struct {
	Pos    int    `json:"position"`
	Action string `json:"action"`
}

type GameState struct {
	PlayerLeft  *Player `json:"player_left"`
	PlayerRight *Player `json:"player_right"`
	Ball        *Ball   `json:"ball"`
}

type GameEvent struct {
	b *Bot
}

func (b *Bot) NewGameEvent() *GameEvent {
	return &GameEvent{ b: b}
}

func (g *GameEvent) OnSubscription(channelID int){
	log.Infof("Let's play this game [%d]", channelID)
}

func (g *GameEvent) OnMessage(e []byte, channelID int){
	var state GameState
	err := json.Unmarshal(e, &state)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if state.Ball != nil {
		g.sendPaddlePos(channelID, state.Ball.Y)
	}
}

func (g *GameEvent) sendPaddlePos(channelID int, pos int) {
	m := GameMessage{Action: "received", Pos: pos}
	msg, _ := json.Marshal(m)
	g.b.Ac.SendMessage("GameChannel", channelID, string(msg))
}

func (b *Bot) SubscribeGame(ID int) {
	log.Debug("I joined game_id [", ID, "]")
	b.Ac.RegisterChannel("GameChannel", b.NewGameEvent())
	b.Ac.Subscribe("GameChannel", ID)
}
