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

type Team int

const (
	Unassigned Team = -1
	Red        Team = 0
	Blue       Team = 1
)

type PlayerInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Team    Team   `json:"team"`
	IsReady bool   `json:"isReady"`
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

func generatePlayerId() string {
	return uuid.NewString()
}
