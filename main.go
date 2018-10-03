package main

import (
	"./webSocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"./traderInfo"
	"./traderBot"
)

func main() {

	account := traderInfo.GetAccount()
	if account != nil {

		traderBot.Run(account)

		router := mux.NewRouter()
		router.HandleFunc("/ws", webSocket.HandleConnections)
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./src/static/")))

		go webSocket.Sender()

		log.Println("http server started on :8081")
		err := http.ListenAndServe(":8081", router)
		if err != nil {
			log.Panic(err)
		}
	}
}
