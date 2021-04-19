package actioncable

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Channel string `json:"channel"`
	ID      int    `json:"id"`
}

type Message struct {
	Command    string `json:"command"`
	Data       string `json:"data,omitempty"`
	Identifier string `json:"identifier"`
}

func (ac *ActionCable) SendMessage(channelName string, channelID int, content string) {
	data, err := json.Marshal(Command{
		Channel: channelName,
		ID:      channelID,
	})
	if err != nil {
		log.Fatal("Unable to marshal:", err)
	}
	ac.sendCh <- &Message{
		Command:    "message",
		Identifier: string(data),
		Data:       content,
	}
}

func (ac *ActionCable) sendRoutine() {
	ac.wg.Add(1)
	for {
		select {
		case m := <-ac.sendCh:
			log.Debug("sent:", m)
			if err := ac.ws.WriteJSON(m); err != nil {
				log.Error("Unable to send msg:", err)
			}
		case <-ac.pongCh:
			if err := ac.ws.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				log.Error("Unable to send ping:", err)
			}
		case <-ac.stopCh:
			ac.wg.Done()
			return
		}
	}
}
