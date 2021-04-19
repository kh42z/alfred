package actioncable

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

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
		}
	}
}
