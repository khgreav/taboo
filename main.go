package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

var addr = "localhost:8081"

func playerConnHandler(game *Game, w http.ResponseWriter, r *http.Request) {
	playerId := r.Header.Get("X-Player-ID")
	if playerId == "" {
		playerId = generatePlayerId()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Error accepting client connection:", err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		slog.Info("Player %s disconnected: %s (%d)", playerId, text, code)
		game.disconnect <- playerId
		return nil
	})

	game.connect <- &ConnectRequest{
		id:   playerId,
		conn: conn,
	}

	for {
		mtype, data, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message: ", "err", err)
			break
		}

		slog.Debug("Incoming message", "type", mtype, "content", data)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	game := createGame()
	go game.run() // TODO: MULTIPLE GAME ROOMS
	http.Handle("/", http.FileServer(http.Dir("frontend/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		playerConnHandler(game, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
