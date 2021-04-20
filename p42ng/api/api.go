package api

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type PongAPI struct {
	Token  string `json:"access-token"`
	Client string `json:"client"`
	Uid    string `json:"uid"`
	UserID int    `json:"-"`
	host   string `json:"-"`
}

func NewAPI(host, code string, uid int) *PongAPI {
	api := PongAPI{UserID: uid, host: host}
	api.Connect(code)
	return &api
}

func readResponse(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Unable to read request:", err)
		return []byte("")
	}
	return body
}

func (p *PongAPI) setReqHeaders(req *http.Request) {
	req.Header.Add("Origin", p.host)
	req.Header.Add("access-token", p.Token)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("client", p.Client)
	req.Header.Add("uid", p.Uid)
}

func (p *PongAPI) GenerateAuthHeaders() http.Header {
	req, _ := http.NewRequest("GET", p.host, nil)
	p.setReqHeaders(req)
	return req.Header
}
