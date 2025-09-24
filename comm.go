package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

func sendUnicastMessage(player *Player, msg MessageBase) error {
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

func broadcastMessage(players map[string]*Player, msg MessageBase, excluded *string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal %s message: %w", msg.GetType(), err)
	}

	for _, player := range players {
		if excluded != nil && player.id == *excluded {
			continue
		}
		err = player.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			slog.Warn(
				"Failed to send name changed message",
				slog.String("player_id", player.id),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}
