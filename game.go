package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type Game struct {
	playerMtx sync.RWMutex
	// Connected players
	players map[string]*Player
	// Channel for incoming player messages
	messages chan MessageBase
}

func createGame() *Game {
	return &Game{
		players:  make(map[string]*Player, 4),
		messages: make(chan MessageBase),
	}
}

func (g *Game) run() {
	for {
		message := <-g.messages
		switch message := message.(type) {
		default:
			slog.Warn("Unknown message type", "type", message.GetType())
		}
	}
}

func (g *Game) AddPlayer(conn *websocket.Conn, playerId *string) string {
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	if playerId == nil {
		newId := generatePlayerId()
		// new player without ID
		newPlayer := &Player{
			id:   newId,
			conn: conn,
			team: -1,
		}
		g.players[newId] = newPlayer
		msg := ConnectAckMessage{
			TypeMessage: TypeMessage{
				Type: ConnectAckMsg,
			},
			PlayerId: newId,
		}
		sendMessage(newPlayer, msg)
		return newId
	} else {
		return ""
	}
}

func sendMessage(player *Player, msg MessageBase) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal %s message: %w", msg.GetType(), err)
	}

	err = player.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("failed to send message %s to player %s: %w", msg.GetType(), player.id, err)
	}

	return nil
}
