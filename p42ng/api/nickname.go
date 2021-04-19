package api

import "fmt"

func (api *PongAPI) UpdateNickname(name string) {
	target := fmt.Sprintf("/users/%d", api.UserID)
	payload := fmt.Sprintf("{\"nickname\":\"%s\"}", name)
	api.DoPatch(target, payload)
}
