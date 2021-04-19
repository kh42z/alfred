package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func (p *PongAPI) DoPatch(target, payload string) {
	url := fmt.Sprintf("%s/api%s", p.host, target)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(payload))
	if err != nil {
		log.Error("Unable to build patch request ", err)
		return
	}
	p.setReqHeaders(request)
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		log.Warn("Unable to change nickname, unfortunate", err)
		return
	}
	ioutil.ReadAll(response.Body)
	response.Body.Close()
}
