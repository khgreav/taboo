package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Player struct {
	// client websocket connection
	conn *websocket.Conn
	// player team
	team int
}

func newPlayerConnHandler(game *Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error accepting client connection:", err)
		return
	}
	player := &Player{conn: conn, team: -1}
	game.connect <- player
}
