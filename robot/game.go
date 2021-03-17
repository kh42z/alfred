package robot

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
