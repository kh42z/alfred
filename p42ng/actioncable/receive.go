package actioncable

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	Message    json.RawMessage `json:"message"`
	Identifier string          `json:"identifier"`
}

func (ac *ActionCable) receiveRoutine() {
	ac.wg.Add(1)
	for {
		_, message, err := ac.ws.ReadMessage()
		if err != nil {
			log.Debug("Unable to rcv:", err)
			break
		}
		var e map[string]interface{}
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
			return
		}
		if _, ok := e["type"]; ok {
			ac.internalMessage(e)
		} else {
			var e Event
			err := json.Unmarshal(message, &e)
			if err != nil {
				log.Error("Unable to unmarshal Event:", err)
				return
			}
			ac.dispatchChannel(&e)
		}
	}
	ac.stopCh <- true
	ac.wg.Done()
}

func (ac *ActionCable) internalMessage(e map[string]interface{}) {

	switch e["type"] {
	case "welcome":
		log.Debug("Connected to ActionCable")
	case "confirm_subscription":
		log.Debugf("I'm listening, Sir: %s", e["identifier"])
	case "ping":
		ac.pongCh <- true
	default:
		log.Warn("unknown internal type rcv:", e["type"])
	}
}
