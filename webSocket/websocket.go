package webSocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"../traderBot/Worker"
	"../traderInfo"
	"time"
)

var usersWs = make(map[*websocket.Conn]*Clients)

var upgrader = websocket.Upgrader{}

type Clients struct {
	Id int
}

type Status struct {
	Event   string `json:"event"`
	Workers map[string]*Worker.Worker
	Account *traderInfo.Account
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	AddNewUser(ws)
}

func AddNewUser(ws *websocket.Conn) {
	usersWs[ws] = &Clients{}
}

func Sender() {
	for {
		resp := Status{Event: "UpdateStatus", Workers: Worker.GetPoolWolker()}

		for ws := range usersWs {
			err := ws.WriteJSON(resp)
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(usersWs, ws)
			}
		}

		time.Sleep(time.Millisecond * 300)
	}
}
