package robot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type PongAPI struct {
	Token  string `json:"access-token"`
	Client string `json:"client"`
	Uid    string `json:"uid"`
	UserID int    `json:"-"`
	host   string
}

func (b *Bot) JoinGuild(guild_id int) {
	_, err := b.api.DoPost(`{ "user_id": `+strconv.Itoa(b.api.UserID)+` }`, "/guilds/"+strconv.Itoa(guild_id)+"/members")
	if err != nil {
		log.Error("Unable to join guild", err)
		return
	}
	resp, err := b.api.DoGet("/guilds/" + strconv.Itoa(guild_id))

	var e map[string]interface{}
	json.Unmarshal(resp, &e)
	if name, ok := e["name"]; ok {
		log.Info("I joined the guild [", name, "]")
	}
}

func (b *Bot) UpdateNickname(name string) {
	target := fmt.Sprintf("/users/%d", b.api.UserID)
	payload := fmt.Sprintf("{\"nickname\":\"%s\"}", name)
	b.api.DoPatch(target, payload)
}

func (p *PongAPI) DoPatch(target, payload string) {
	url := fmt.Sprintf("http://%s/api%s", p.host, target)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(payload))
	if err != nil {
		log.Error("Unable to build patch request ", err)
		return
	}
	p.setReqHeaders(request, p.host)
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		log.Warn("Unable to change nickname, unfortunate", err)
		return
	}
	ioutil.ReadAll(response.Body)
	response.Body.Close()
}

func (p *PongAPI) DoGet(target string) ([]byte, error) {
	url := "http://" + p.host + "/api" + target
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	p.setReqHeaders(req, p.host)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readResponse(resp), nil
}

func (p *PongAPI) DoPost(body, target string) ([]byte, error) {
	jsonStr := []byte(body)
	url := "http://" + p.host + "/api" + target
	log.Debug("I'm sending a request on: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	p.setReqHeaders(req, p.host)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return nil, errors.New("API responded with status code: " + strconv.Itoa(resp.StatusCode) + " > " + string(readResponse(resp)))
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

func (p *PongAPI) setReqHeaders(req *http.Request, host string) {
	req.Header.Add("Origin", "http://"+host)
	req.Header.Add("access-token", p.Token)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("client", p.Client)
	req.Header.Add("uid", p.Uid)
}
