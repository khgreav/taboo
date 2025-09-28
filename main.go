package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var addr = "localhost:8081"

func playerConnHandler(game *Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Error accepting client connection.", "err", err)
		return
	}

	ss, _ := GetSchemaStorage()

	go func() {

		var playerId string

		defer func() {
			if playerId != "" {
				game.RemovePlayer(playerId)
			}
			slog.Info("Client disconnected.", "playerId", playerId)
			conn.Close()
		}()

		for {
			mtype, data, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Failed to read message.", "err", err)
				break
			}
			slog.Debug("Incoming message.", "type", mtype, "content", data)

			msg, err := decodeIncomingMessage(ss, data)
			if err != nil {
				slog.Error("Failed to process incoming message.", "err", err)
				continue
			}

			if msg.GetType() == ConnectMsg {
				conMsg, ok := msg.(*ConnectMessage)
				if !ok {
					slog.Error("Failed to cast message to ConnectMessage")
					continue
				}
				playerId = game.AddPlayer(conn, conMsg.PlayerId, conMsg.Name)
				slog.Debug("Player ID stored", "playerId", playerId)
				continue
			} else {
				playerMsg, ok := msg.(PlayerMessage)
				if !ok {
					slog.Error("Failed to cast message to PlayerMessage")
					continue
				}
				if playerMsg.GetPlayerId() != playerId {
					slog.Error("Player ID mismatch", "expected", playerId, "got", playerMsg.GetPlayerId())
					continue
				}
			}

			game.messages <- msg
		}
	}()
}

func decodeIncomingMessage(ss *SchemaStorage, data []byte) (MessageBase, error) {
	var typeMsg TypeProperty
	if err := json.Unmarshal(data, &typeMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal incoming JSON message: %w", err)
	}

	if err := ss.validate(typeMsg.Type, data); err != nil {
		return nil, fmt.Errorf("failed to validate incoming %s message: %w", typeMsg.Type, err)
	}

	message, err := ConstructMessageContainer(typeMsg.Type)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, message); err != nil {
		return nil, fmt.Errorf("failed to parse incoming %s JSON message: %w", typeMsg.Type, err)
	}

	return message, nil
}

func initStorage() error {
	_, err := GetSchemaStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize schema storage")
	}
	_, err = GetWordStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize word storage")
	}
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	if err := initStorage(); err != nil {
		slog.Error("Failed to initialize game systems: %w", "err", err)
		os.Exit(1)
	}
	game := CreateGame()
	go game.run() // TODO: MULTIPLE GAME ROOMS
	http.Handle("/", http.FileServer(http.Dir("frontend/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		playerConnHandler(game, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
