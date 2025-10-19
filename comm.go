package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

func SendUnicastMessage(player *Player, msg MessageBase) error {
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

func SendDirectErrorMessage(conn *websocket.Conn, errorMsg ErrorResponseMessage) {
	var err error

	data, err := json.Marshal(errorMsg)
	if err != nil {
		slog.Error(
			"Failed to marshal error message",
			slog.String("error", err.Error()),
		)
		return
	}

	ss, err := GetSchemaStorage()
	if err != nil {
		slog.Error(
			"Failed to get schema from storage.",
			slog.String("error", err.Error()),
		)
		return
	}

	err = ss.validate(ErrorResponseMsg, data)
	if err != nil {
		slog.Error(
			"Failed to validate outgoing error message",
			slog.String("error", err.Error()),
		)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		slog.Error(
			"Failed to send error message.",
			slog.String("client", conn.RemoteAddr().String()),
			slog.String("error", err.Error()),
		)
		return
	}

	slog.Debug(
		"Outgoing unicast error message.",
		slog.String("content", string(data)),
	)
}

func SendErrorMessage(player *Player, errorMsg ErrorResponseMessage) {
	var err error

	data, err := json.Marshal(errorMsg)
	if err != nil {
		slog.Error(
			"Failed to marshal error message",
			slog.String("error", err.Error()),
		)
		return
	}

	ss, err := GetSchemaStorage()
	if err != nil {
		slog.Error(
			"Failed to get schema from storage.",
			slog.String("error", err.Error()),
		)
		return
	}

	err = ss.validate(ErrorResponseMsg, data)
	if err != nil {
		slog.Error(
			"Failed to validate outgoing error message",
			slog.String("error", err.Error()),
		)
		return
	}

	err = player.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		slog.Error(
			"Failed to send error message.",
			slog.String("player_id", player.id),
			slog.String("error", err.Error()),
		)
		return
	}

	slog.Debug(
		"Outgoing unicast error message.",
		slog.String("content", string(data)),
	)
}

func BroadcastMessage(players map[string]*Player, msg MessageBase, excluded *string) error {
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
				"Failed to send message",
				slog.String("player_id", player.id),
				slog.String("message_type", string(msg.GetType())),
				slog.String("error", err.Error()),
			)
		}
	}

	slog.Debug("Outgoing broadcast message.", "type", msg.GetType(), "content", msg)

	return nil
}
