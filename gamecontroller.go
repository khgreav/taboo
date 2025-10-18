package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type GameController struct {
	gameManager *GameManager
}

type GameCreated struct {
	GameId string `json:"gameId"`
}

func (gc *GameController) CreateGame(w http.ResponseWriter, r *http.Request) {
	game := gc.gameManager.CreateGame()
	rsp := GameCreated{GameId: game.id}
	byteRsp, err := json.Marshal(rsp)
	if err != nil {
		gc.gameManager.RemoveGame(game.id)
		slog.Error("Failed to marshal game created response", "err", err)
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(byteRsp)
}

func (gc *GameController) JoinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

}
