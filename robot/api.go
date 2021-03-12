package robot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (b *Bot)DoPost(body, target string) {
	jsonStr := []byte(body)
	url := "http://" + b.host + "/api" + target
	log.Info("I'm sending a request on: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Origin", "http://" + b.host)
	req.Header.Add("access-token", b.bearerToken.Token)
	req.Header.Add("content-type","application/json")
	req.Header.Add("client", b.bearerToken.Client)
	req.Header.Add("uid", b.bearerToken.Uid)
	if err != nil {
		log.Error("unable to fmt request", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Unable to execute request:", err)
		return
	}
	defer resp.Body.Close()
	strBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Error("Unable to read request:", err)
		return
	}
	log.Info("I received the following response: "+ string(strBody))
}

func (b *Bot)join_guild(guild_id int){
	b.DoPost(`{ "user_id": 1 }`, "/guilds/" + strconv.Itoa(guild_id)  +"/members")
}
