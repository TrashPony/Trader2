package webSocket

import (
	"github.com/TrashPony/Trader2/traderBot/Worker"
	"github.com/TrashPony/Trader2/traderInfo"
	"github.com/gorilla/websocket"
	"net/http"
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
	}
	usersWs[ws] = &Clients{}
}

func Sender() {
	for {
		resp := Status{Event: "UpdateStatus", Workers: Worker.GetPoolWorker()}

		for ws := range usersWs {
			err := ws.WriteJSON(resp)
			if err != nil {
				ws.Close()
				delete(usersWs, ws)
			}
		}

		time.Sleep(time.Millisecond * 300)
	}
}
