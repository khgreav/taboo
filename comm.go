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

	ss, err := GetSchemaStorage()
	if err != nil {
		return fmt.Errorf("failed to get schema storage: %w", err)
	}

	err = ss.validate(msg.GetType(), data)
	if err != nil {
		return fmt.Errorf("failed to validate outgoing %s message: %w", msg.GetType(), err)
	}

	err = player.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("failed to send message %s to player %s: %w", msg.GetType(), player.id, err)
	}

	slog.Debug("Outgoing unicast message.", "type", msg.GetType(), "content", msg)

	return nil
}

func sendErrorMessage(player *Player, failedMessage MessageType, errorMessage string) error {
	msg := &ErrorResponseMessage{
		TypeProperty: TypeProperty{
			Type: ErrorResponseMsg,
		},
		FailedType: failedMessage,
		Error:      errorMessage,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal %s message: %w", ErrorResponseMsg, err)
	}

	ss, err := GetSchemaStorage()
	if err != nil {
		return fmt.Errorf("failed to get schema storage: %w", err)
	}

	err = ss.validate(ErrorResponseMsg, data)
	if err != nil {
		return fmt.Errorf("failed to validate outgoing %s message: %w", ErrorResponseMsg, err)
	}

	err = player.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("failed to send message %s to player %s: %w", ErrorResponseMsg, player.id, err)
	}

	slog.Debug("Outgoing unicast error message.", "content", msg)

	return nil
}

func broadcastMessage(players map[string]*Player, msg MessageBase, excluded *string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal %s message: %w", msg.GetType(), err)
	}

	ss, err := GetSchemaStorage()
	if err != nil {
		return fmt.Errorf("failed to get schema storage: %w", err)
	}

	err = ss.validate(msg.GetType(), data)
	if err != nil {
		return fmt.Errorf("failed to validate outgoing %s message: %w", msg.GetType(), err)
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

	slog.Debug("Outgoing broadcast message.", "type", msg.GetType(), "content", msg)

	return nil
}
