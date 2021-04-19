package api

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (p *PongAPI) DoPost(body, target string) ([]byte, error) {
	jsonStr := []byte(body)
	url := p.host + "/api" + target
	log.Debug("I'm sending a request on: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	p.setReqHeaders(req)
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
