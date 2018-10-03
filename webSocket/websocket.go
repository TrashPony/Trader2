package webSocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}

var usersWs = make(map[*websocket.Conn]*Clients)
var wsPipe = make(chan chatResponse)

var upgrader = websocket.Upgrader{}

type Clients struct {
	Id int
}

type chatMessage struct {
	Event    string `json:"event"`
	UserName string `json:"user_name"`
	Message  string `json:"message"`
}

type chatResponse struct {
	Event    string `json:"event"`
	UserName string `json:"user_name"`
	GameUser string `json:"game_user"`
	Message  string `json:"message"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	AddNewUser(ws)
}

func AddNewUser(ws *websocket.Conn) {
	usersWs[ws] = &Clients{Id: 1}
	defer ws.Close()

	Reader(ws)
}

func Reader(ws *websocket.Conn) {
	for {
		var msg chatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(usersWs, ws)
			break
		}
		if msg.Event == "NewMessage" {
			var resp = chatResponse{Event: msg.Event, Message: msg.Message}
			wsPipe <- resp
		}
	}
}

func Sender() {
	for {
		resp := <-wsPipe
		mutex.Lock()
		for ws := range usersWs {
			err := ws.WriteJSON(resp)
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(usersWs, ws)
			}
		}
		mutex.Unlock()
	}
}
