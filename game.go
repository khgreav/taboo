package main

import (
	"fmt"
	"log"
)

type Game struct {
	// Connected players
	players map[*Player]bool
	// Channel for new player connections
	connect chan *Player
	// Channel for player disconnections
	disconnect chan *Player
}

func createGame() *Game {
	return &Game{
		players:    make(map[*Player]bool),
		connect:    make(chan *Player),
		disconnect: make(chan *Player),
	}
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.connect:
			log.Println("New player connected from:", player.conn.RemoteAddr().String())
			g.players[player] = true
			player.conn.WriteMessage(1, []byte(
				fmt.Sprintf(
					"Welcome player %s, you are on team %d",
					player.conn.RemoteAddr().String(),
					player.team,
				)),
			)
		case player := <-g.disconnect:
			_, exists := g.players[player]
			if exists {
				delete(g.players, player)
				log.Println("Player disconnected from:", player.conn.RemoteAddr().String())
			}
		}
	}
}
