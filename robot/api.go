package robot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type PongAPI struct {
	Token string `json:"access-token"`
	Client string `json:"client"`
	Uid string `json:"uid"`
}


func (b *Bot) JoinGuild(guild_id int){
	_, err := b.api.DoPost(`{ "user_id": 1 }`, b.host,"/guilds/" + strconv.Itoa(guild_id)  +"/members")
	if err != nil {
		log.Error("Unable to join guild", err)
		return
	}
	resp, err := b.api.DoGet(b.host, "/api/guilds/" + strconv.Itoa(guild_id))
	var e map[string]interface{}
	json.Unmarshal(resp, e)
	if name, ok := e["name"]; ok {
		log.Info("I joined the guild", name)
	}
}

func (p *PongAPI)DoGet(host, target string) ([]byte, error) {
	url := "http://" + host + "/api" + target
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	p.setReqHeaders(req, host)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readResponse(resp), nil
}

func (p *PongAPI) DoPost(body, host, target string) ([]byte, error) {
	jsonStr := []byte(body)
	url := "http://" + host + "/api" + target
	log.Info("I'm sending a request on: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	p.setReqHeaders(req, host)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return nil, errors.New("API responded with status code: "+ strconv.Itoa(resp.StatusCode))
	}
	return readResponse(resp), nil
}


func readResponse(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Unable to read request:", err)
		return []byte("")
	}
	return body
}

func (p *PongAPI) setReqHeaders(req *http.Request, host string){
	req.Header.Add("Origin", "http://" + host)
	req.Header.Add("access-token", p.Token)
	req.Header.Add("content-type","application/json")
	req.Header.Add("client", p.Client)
	req.Header.Add("uid", p.Uid)
}
