package robot

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Chat struct {
	ID int `json:"id"`
}

func (b *Bot) connect(code string) {
	b.twoFactorSignIn(code)
	req,_ := http.NewRequest("GET", "http://" + b.host, nil)
	b.api.setReqHeaders(req, b.host)
	var err error
	b.ws, _, err = websocket.DefaultDialer.Dial("ws://" + b.host + "/cable", req.Header)
	if err != nil {
		log.Fatal("Unable to connect to websocket:", err)
	}
}

func (b *Bot) twoFactorSignIn(code string) {
	var resp *http.Response
	url := fmt.Sprintf("http://%s/two_factor/%d/?code=%s" , b.host, b.api.UserID, code)
	for {
		var err error
		resp, err = http.Get(url)
		if err != nil {
			log.Infof("Waiting for API to be up and running at %s", url)
			time.Sleep(10 * time.Second)
		}else{
			break
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Unable to authentifcate: ", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	log.Infof("BODY [%s]", body);
	if err != nil {
		log.Fatal("Unable to read body: ", err)
	}
	err = json.Unmarshal(body, &b.api)
	if err != nil {
		log.Fatal("Unable to unmarshal auth:", err)
	}
}

func (b *Bot) retrieveSubscriptions() []*Chat {
	body, err := b.api.DoGet(b.host, fmt.Sprintf("/chats?participant_id=%b", b.api.UserID))
	if err != nil {
		log.Fatal("Unable to retrieve chatrooms subscriptions", err)
	}
	var chats []*Chat
	err = json.Unmarshal(body, &chats)
	if err != nil {
		log.Error("Unable to unmarshal:", err)
		return nil
	}
	return chats
}
