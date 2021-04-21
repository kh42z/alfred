package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (api *PongAPI) Connect(code string) {
	api.twoFactorSignIn(code)
}

func (api *PongAPI) twoFactorSignIn(code string) {
	var resp *http.Response
	url := fmt.Sprintf("%s/two_factor/%d/?code=%s", api.host, api.UserID, code)
	var err error
	resp, err = http.Get(url)
	if err != nil {
		log.Fatalf("Connection refused: %s\n", api.host)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Invalid response code (%d) with %s ", resp.StatusCode, url)
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
