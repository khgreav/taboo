package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // DEV ONLY
	},
}

type PlayerInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Player struct {
	// Player ID
	id string
	// client websocket connection
	conn *websocket.Conn
	// Player name
	name string
	// player team
	team int
}

func (p *Player) SetName(name string) {
	p.name = name
}

func generatePlayerId() string {
	return uuid.NewString()
}
