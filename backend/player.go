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
	// Session token
	sessionToken string
	// Player name
	name string
	// Player ready status
	isReady bool
	// player team
	team Team
	// Player connected status
	connected bool
}

func (p *Player) SetConnection(conn *websocket.Conn) {
	p.conn = conn
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
		TypeProperty: TypeProperty{
			Type: PlayerJoinedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: p.id,
		},
		Name: p.name,
	}
}

func (p Player) CreatePlayerLeftMessage() *PlayerLeftMessage {
	return &PlayerLeftMessage{
		TypeProperty: TypeProperty{
			Type: PlayerLeftMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: p.id,
		},
	}
}

func (p Player) CreatePlayerDisconnectedMessage() *PlayerDisconnectedMessage {
	return &PlayerDisconnectedMessage{
		TypeProperty: TypeProperty{
			Type: PlayerDisconnectedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: p.id,
		},
	}
}

func (p Player) CreatePlayerReconnectedMessage() *PlayerReconnectedMessage {
	return &PlayerReconnectedMessage{
		TypeProperty: TypeProperty{
			Type: PlayerReconnectedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: p.id,
		},
	}
}
