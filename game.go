package main

import (
	"log"
)

type Game struct {
	// Connected players
	players map[string]*Player
	// Channel for new player connections
	connect chan *ConnectRequest
	// Channel for player disconnections
	disconnect chan string
}

func createGame() *Game {
	return &Game{
		players:    make(map[string]*Player),
		connect:    make(chan *ConnectRequest),
		disconnect: make(chan string),
	}
}

func (g *Game) run() {
	for {
		select {
		case request := <-g.connect:
			log.Printf("New client %s connected from %s", request.id, request.conn.RemoteAddr().String())
			player, exists := g.players[request.id]
			if exists {
				if player.conn != nil {
					player.conn.Close()
				}
				player.conn = request.conn
				log.Printf("Player %s reconnected from %s", request.id, request.conn.RemoteAddr().String())
			} else {
				player = &Player{
					conn: request.conn,
					team: -1,
				}
				g.players[request.id] = player
				msg := map[string]string{
					"type": "new_client_id",
					"id":   request.id,
				}
				err := player.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Failed to send client ID to %s: %s", request.id, err)
				}
				log.Printf("Player %s connected from %s", request.id, request.conn.RemoteAddr().String())
			}
		case playerId := <-g.disconnect:
			player, exists := g.players[playerId]
			if exists {
				delete(g.players, playerId)
				log.Println("Player disconnected from:", player.conn.RemoteAddr().String())
			}
		}
	}
}
