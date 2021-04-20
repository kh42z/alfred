package actioncable

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	Message    json.RawMessage `json:"message"`
	Type	   string	       `json:"type"`
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
		var e Event
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Error("Unable to unmarshal rcv", err)
			return
		}
		if len(e.Type) > 0 {
			ac.internalMessage(&e)
		} else {
			ac.dispatchChannel(&e)
		}
	}
	ac.stopCh <- true
	ac.wg.Done()
}

func (ac *ActionCable) internalMessage(e *Event) {

	switch e.Type {
	case "welcome":
		log.Debug("Connected to ActionCable")
		ac.startCh <- true
	case "confirm_subscription":
		var i Identifier
		err := json.Unmarshal([]byte(e.Identifier), &i)
		if err != nil {
			log.Warn("Unable to unmarshal Identifier", i)
			return
		}
		for name, e := range ac.channels {
			if name == i.Channel {
				e.OnSubscription(i.ID)
			}
		}
	case "disconnect":
		log.Warn("We got disconnected.")
		ac.startCh <- false
	case "ping":
	default:
		log.Warn("unknown internal type rcv:", e.Type)
	}
}
