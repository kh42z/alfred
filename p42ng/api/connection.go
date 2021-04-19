package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func (api *PongAPI) Connect(code string) {
	api.twoFactorSignIn(code)
}

func (api *PongAPI) twoFactorSignIn(code string) {
	var resp *http.Response
	url := fmt.Sprintf("%s/two_factor/%d/?code=%s", api.host, api.UserID, code)
	log.Infof("Waiting for API to be up and running at %s", url)
	for {
		var err error
		resp, err = http.Get(url)
		if err != nil {
			log.Infof("Waiting for API to be up and running at %s", url)
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Unable to authenticate: ", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body: ", err)
	}
	err = json.Unmarshal(body, &api)
	if err != nil {
		log.Fatal("Unable to unmarshal auth:", err)
	}
}
