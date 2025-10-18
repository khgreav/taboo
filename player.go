package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // DEV ONLY
	},
}

type PlayerInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Team      Team   `json:"team"`
	IsReady   bool   `json:"isReady"`
	Connected bool   `json:"connected"`
}

type Player struct {
	// Player ID
	id string
	// client websocket connection
	conn *websocket.Conn
	// Player name
	name string
	// Player ready status
	isReady bool
	// player team
	team Team
	// is connected
	connected bool
	// session token
	sessionToken string
}

func (p *Player) SetName(name string) {
	p.name = name
}

func (p *Player) SetReady(ready bool) {
	p.isReady = ready
}

func (p *Player) SetTeam(team Team) {
	p.team = team
}

func (p *Player) SetConnected(connected bool) {
	p.connected = connected
}

func (p Player) CreatePlayerJoinedMessage() *PlayerJoinedMessage {
	return &PlayerJoinedMessage{
		TypeProperty:     TypeProperty{Type: PlayerJoinedMsg},
		PlayerIdProperty: PlayerIdProperty{PlayerId: p.id},
		Name:             p.name,
	}
}

func (p Player) CreatePlayerDisconnectedMessage() *PlayerDisconnectedMessage {
	return &PlayerDisconnectedMessage{
		TypeProperty:     TypeProperty{Type: PlayerDisconnectedMsg},
		PlayerIdProperty: PlayerIdProperty{PlayerId: p.id},
	}
}

func (p Player) CreatePlayerReconnectedMessage() *PlayerReconnectedMessage {
	return &PlayerReconnectedMessage{
		TypeProperty:     TypeProperty{Type: PlayerReconnectedMsg},
		PlayerIdProperty: PlayerIdProperty{PlayerId: p.id},
	}
}
