package main

import (
	"github.com/TrashPony/Trader2/traderBot"
	"github.com/TrashPony/Trader2/traderInfo"
	"github.com/TrashPony/Trader2/webSocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	account := traderInfo.GetAccount()
	if account != nil {
		go traderBot.Run(account)
	}

	router := mux.NewRouter()
	router.HandleFunc("/ws", webSocket.HandleConnections)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	go webSocket.Sender()

	log.Println("http server started on :8083")
	err := http.ListenAndServe(":8083", router)
	if err != nil {
		log.Panic(err)
	}
}
