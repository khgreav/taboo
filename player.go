package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Player struct {
	// Player ID
	id string
	// client websocket connection
	conn *websocket.Conn
	// player team
	team int
}

func generatePlayerId() string {
	return uuid.NewString()
}
