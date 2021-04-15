package robot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ActivityMessage struct {
	Action string `json:"action"`
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func ActivityUpdate(e []byte) {
	var activityMessage ActivityMessage
	err := json.Unmarshal(e, &activityMessage)
	if err != nil {
		log.Error("Unable to unmarshal content", err)
		return
	}
	if activityMessage.Action == "user_update_status" {
		log.Infof("Seems like [%d] status changed to <%s>", activityMessage.ID, activityMessage.Status)
	}
}
