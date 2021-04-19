package api

import "net/http"

func (p *PongAPI) DoGet(target string) ([]byte, error) {
	url := p.host + "/api" + target
	req, err := http.NewRequest("GET", url, nil)
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
	return readResponse(resp), nil
}
