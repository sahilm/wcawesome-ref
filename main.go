package main

import (
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sync/atomic"
	"os"
	"github.com/fatih/color"
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

const goal = `
                                       
                                   88  
                                   88  
                                   88  
 ,adPPYb,d8  ,adPPYba,  ,adPPYYba, 88  
a8"     Y88 a8"     "8a ""      Y8 88  
8b       88 8b       d8 ,adPPPPP88 88  
"8a,   ,d88 "8a,   ,a8" 88,    ,88 88  
  "YbbdP"Y8   "YbbdP"'   "8bbdP"Y8 88
aa,    ,88
"Y8bbdP"
`

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
	} else {
		if strings.Contains(strings.ToLower(RefNotification.Event.TypeOfEvent), "yellow") {
			color.Yellow(RefNotification.Event.Time +
				" " + RefNotification.Country +
				" " + RefNotification.Event.Player +
				" " + RefNotification.Event.TypeOfEvent)
		} else {
			fmt.Println(RefNotification.Event.Time +
				" " + RefNotification.Country +
				" " + RefNotification.Event.Player +
				" " + RefNotification.Event.TypeOfEvent)
		}

		if strings.Contains(strings.ToLower(RefNotification.Event.TypeOfEvent), "goal") {
			color.Green(goal)
		}

		w.Write([]byte(message))
	}

	if atomic.LoadInt64(&sendkill) == 2 {
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
}

func main() {
	http.HandleFunc("/", notificationHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
