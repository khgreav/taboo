package main

import (
	"fmt"
	"log/slog"
	"maps"
	"sync"

	"github.com/gorilla/websocket"
)

type Game struct {
	// Player mutex
	playerMtx 	sync.RWMutex
	// Connected players
	players 	map[string]*Player
	// Channel for incoming player messages
	messages 	chan MessageBase
}

func createGame() *Game {
	return &Game{
		players:  make(map[string]*Player, 4),
		messages: make(chan MessageBase),
	}
}

func (g *Game) run() {
	var err error
	for {
		message := <-g.messages
		switch message := message.(type) {
		case *ChangeNameMessage:
			err = g.changePlayerName(message.PlayerId, message.Name)
		default:
			slog.Warn("Unknown message type", "type", message.GetType())
		}
		if err != nil {
			slog.Error("Failed to process message", "type", message.GetType(), "err", err)
			err = nil
		}
	}
}

func (g *Game) AddPlayer(conn *websocket.Conn, playerId *string, name string) string {
	g.playerMtx.Lock()

	var player Player

	if playerId == nil {
		newId := generatePlayerId()
		// new player without ID
		player = Player{
			id:   newId,
			name: name,
			conn: conn,
			team: -1,
		}
		g.players[newId] = &player
	} else {
		return ""
	}

	// get copy of players to unlock early to not block other operations while sending messages
	players := g.GetPlayersCopyUnlocked()
	g.playerMtx.Unlock()

	// create a connect ack message for the new player
	msg := ConnectAckMessage{
		TypeProperty: TypeProperty{
			Type: ConnectAckMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: player.id,
		},
	}
	// send connect ack to the new player
	sendUnicastMessage(&player, msg)

	// create a player list message to send lobby state to player
	listMsg := &PlayerListMessage{
		TypeProperty: TypeProperty{
			Type: PlayerListMsg,
		},
		Players: g.CreatePlayerList(),
	}
	// send player list to the new player
	if err := sendUnicastMessage(&player, listMsg); err != nil {
		slog.Warn(
			"Failed to send player list message",
			slog.String("player_id", player.id),
			slog.String("error", err.Error()),
		)
	}
	// create a player joined message to notify other players
	joinedMsg := &PlayerJoinedMessage{
		TypeProperty: TypeProperty{
			Type: PlayerJoinedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: player.id,
		},
		Name: name,
	}
	// broadcast player joined message to all other players, excluding the new player
	broadcastMessage(players, joinedMsg, &player.id)

	return player.id
}

func (g *Game) RemovePlayer(playerId string) {
	g.playerMtx.Lock()
	// delete player from map
	delete(g.players, playerId)
	// get copy of players to unlock early
	players := g.GetPlayersCopyUnlocked()
	g.playerMtx.Unlock()

	// create player left message
	leftMsg := &PlayerLeftMessage{
		TypeProperty: TypeProperty{
			Type: PlayerLeftMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
	}
	// broadcast player left message to all other players
	broadcastMessage(players, leftMsg, nil)
}

func (g *Game) changePlayerName(playerId string, name string) error {
	// lock before accessing players
	g.playerMtx.Lock()
	// find player
	player, exists := g.players[playerId]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerId)
	}
	// change name
	player.SetName(name)
	// create name change message
	msg := &NameChangedMessage{
		TypeProperty: TypeProperty{
			Type: NameChangedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
		Name: name,
	}
	// get copy of players to unlock early
	players := g.GetPlayersCopyUnlocked()
	// unlock before sending messages
	g.playerMtx.Unlock()
	// broadcast name change
	broadcastMessage(players, msg, nil)
	return nil
}

func (g *Game) CreatePlayerList() []PlayerInfo {
	players := g.GetPlayersCopy()
	playerList := make([]PlayerInfo, 0, len(players))
	for _, p := range players {
		playerList = append(playerList, PlayerInfo{
			Id:   p.id,
			Name: p.name,
		})
	}
	return playerList
}

func (g *Game) GetPlayersCopy() map[string]*Player {
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()
	return g.GetPlayersCopyUnlocked()
}

func (g *Game) GetPlayersCopyUnlocked() map[string]*Player {
	playersCopy := make(map[string]*Player, len(g.players))
	maps.Copy(playersCopy, g.players)
	return playersCopy
}
