package main

import (
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sync/atomic"
	"os"
)

type event struct {
	Id          int64  `json:"id"`
	Player      string `json:"player"`
	TypeOfEvent string `json:"type_of_event"`
	Time        string `json:"time"`
}

type refNotification struct {
	Country string `json:"country"`
	Event   event  `json:"event"`
}

var sendkill int64

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	var RefNotification refNotification
	err := json.Unmarshal(bodyBytes, &RefNotification)
	if err != nil {
		fmt.Printf("%s", err)
	}

	if RefNotification.Country == "GAME_OVER" {
		atomic.AddInt64(&sendkill, 1)
	}
	if atomic.LoadInt64(&sendkill) == 2 {
		fmt.Println("KILL! KILL! KILL! KILL! KILL! KILL! KILL! ")
		atomic.StoreInt64(&sendkill, 0)
		req, err := http.NewRequest("POST", orchestratorURL+"/gameover", nil)
		if err != nil {
			fmt.Printf("%s", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	fmt.Println(RefNotification.Event.Time +
		" " + RefNotification.Country +
		" " + RefNotification.Event.Player +
		" " + RefNotification.Event.TypeOfEvent)

	w.Write([]byte(message))
}

func main() {
	http.HandleFunc("/", notificationHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
