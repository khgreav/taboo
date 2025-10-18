package main

import "sync"

type GameManager struct {
	games        map[string]*Game
	playerInGame map[string]string // playerId -> gameId
	mtx          sync.RWMutex
}

func (gm *GameManager) GetGame(gameId string) (*Game, bool) {
	gm.mtx.RLock()
	defer gm.mtx.RUnlock()

	game, exists := gm.games[gameId]
	return game, exists
}

func (gm *GameManager) CreateGame() *Game {
	gm.mtx.Lock()
	defer gm.mtx.Unlock()

	game := CreateGame()
	gm.games[game.id] = game
	return game
}

func (gm *GameManager) RemoveGame(gameId string) {
	gm.mtx.Lock()
	defer gm.mtx.Unlock()
	delete(gm.games, gameId)
}
