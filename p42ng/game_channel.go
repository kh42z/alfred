package p42ng

import (
	"encoding/json"
	"github.com/kh42z/actioncable"
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

type GameEvent struct{}

func (b *Bot) NewGameEvent() *GameEvent {
	return &GameEvent{}
}

func (g *GameEvent) SubscriptionHandler(_ *actioncable.Client, channelID int) {
	log.Infof("Let's play this game [%d]", channelID)
}

func (g *GameEvent) MessageHandler(c *actioncable.Client, e []byte, channelID int) {
	var state GameState
	err := json.Unmarshal(e, &state)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if state.Ball != nil && state.Ball.X > 480 {
		sendPaddlePos(c, channelID, state.Ball.Y)
	}
}

func sendPaddlePos(c *actioncable.Client, channelID int, pos int) {
	m := GameMessage{Action: "received", Pos: pos}
	msg, _ := json.Marshal(m)
	c.SendMessage("GameChannel", channelID, string(msg))
}

func (b *Bot) SubscribeGame(ID int) {
	log.Debug("I joined game_id [", ID, "]")
	b.Ac.AddChannelHandler("GameChannel", b.NewGameEvent())
	b.Ac.Subscribe("GameChannel", ID)
}
