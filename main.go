package main

import (
	"log"
	"net/http"
)

var addr = "localhost:8081"

func main() {
	log.SetFlags(0)
	game := createGame()
	go game.run() // TODO: MULTIPLE GAME ROOMS
	http.Handle("/", http.FileServer(http.Dir("frontend/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		playerConnHandler(game, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
