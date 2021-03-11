package robot

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func (b *Bot)addMessage(event *Event){
	b.statsCh <- true
}

type sendMessageFn func()

func (b *Bot) withWS() {
	b.sendCh <- formatChatMessage("ChatChannel", 1)
}

func (b *Bot) withHTTP() {
	jsonStr := []byte(`{"content":"yes"}`)
	req, err := http.NewRequest("POST", "http://" + b.host + "/api/chats/1/messages", bytes.NewBuffer(jsonStr))
	req.Header.Add("Origin", "http://" + b.host)
	req.Header.Add("access-token", b.bearerToken.Token)
	req.Header.Add("content-type","application/json")
	req.Header.Add("client", b.bearerToken.Client)
	req.Header.Add("uid", b.bearerToken.Uid)
	if err != nil {
		log.Fatal("unable to fmt request", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}

func (b *Bot) generateStats(requestNb int,sendMessage sendMessageFn) {
	currentNb := 0
	start := time.Now()
	for {
		sendMessage()
		<- b.statsCh
		currentNb++
		if currentNb == requestNb {
			break
		}
	}
	totalTime := time.Since(start).Milliseconds()
	averageTime := float64(totalTime) / float64(requestNb)
	log.Info("AverageTimePerBrdcastReceive: ", averageTime, " ms")
}

func (b *Bot)Bench(request int, code string) {
	b.connect(code)
	go b.receiveRoutine(b.addMessage)
	go b.sendRoutine()
	b.SubscribeChat(1)
	time.Sleep(2 * time.Second)
	log.Info("Starting benchmark with native Websocket, sending: ", request, " broadcasts")
	b.generateStats(request, b.withWS)
	log.Info("Starting benchmark with POST /api/chats/1/messages, sending: ", request, " requests")
	b.generateStats(request, b.withHTTP)
}
