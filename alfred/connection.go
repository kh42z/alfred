package alfred

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type BearerToken struct {
	Token string `json:"access-token"`
	Client string `json:"client"`
	Uid string `json:"uid"`
}

type Chat struct {
	ID int `json:"id"`
}

func (b *Bot) connect(host string, code string) {
	b.bearerToken = getBearerToken(host, code)
	req,_ := http.NewRequest("GET", "http://"+ host, nil)
	req.Header.Add("Origin", "http://localhost:3000/")
	req.Header.Add("access-token", b.bearerToken.Token)
	req.Header.Add("client", b.bearerToken.Client)
	req.Header.Add("uid", b.bearerToken.Uid)
	var err error
	b.ws, _, err = websocket.DefaultDialer.Dial("ws://"+ host + "/cable", req.Header)
	if err != nil {
		log.Fatal("Unable to connect to websocket:", err)
	}
}

func getBearerToken(host, code string) *BearerToken {
	var resp *http.Response
	for {
		var err error
		resp, err = http.Get("http://" + host + "/two_factor/1?code=" + code)
		if err != nil {
			time.Sleep(10 * time.Second)
		}else{
			break
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body: ", err)
	}
	var b *BearerToken
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Fatal("Unable to unmarshal auth:", err)
	}
	return b
}

func (b *Bot) retrieveSubscriptions(host string) []*Chat {
	client := &http.Client{}
	req,_ := http.NewRequest("GET", "http://"+  host + "/api/chats?participant_id=1", nil)
	req.Header.Add("Origin", "http://" + host)
	req.Header.Add("access-token", b.bearerToken.Token)
	req.Header.Add("client", b.bearerToken.Client)
	req.Header.Add("uid", b.bearerToken.Uid)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Unable to get Channels:", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Unable to read body: ", err)
		return nil
	}
	var chats []*Chat
	err = json.Unmarshal(body, &chats)
	if err != nil {
		log.Error("Unable to unmarshal:", err)
		return nil
	}
	return chats
}
