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
			log.Fatal("Unable to rcv:", err)
		}
		var e map[string]interface{}
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
			return
		}
		if t, ok := e["type"]; ok {
			switch t {
			case "welcome":
				log.Infof("Alfred at your service")
			case "confirm_subscription":
				log.Infof("I'm listening, Sir: %s", e["identifier"])
			case "ping":
				ac.pongCh <- true
			default:
				log.Warn("rcv:", t, string(message))
			}
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
}
