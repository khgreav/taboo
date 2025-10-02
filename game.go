package main

import (
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type GameState int

const RoundDuration = 60

const (
	InLobby GameState = iota
	InProgress
	InRound
)

type Game struct {
	// Is the game currently running
	gameState GameState
	// Player mutex
	playerMtx sync.RWMutex
	// Connected players
	players map[string]*Player
	// Number of red team players
	redTeamCount int
	// Number of blue team players
	blueTeamCount int
	// Channel for incoming player messages
	messages chan MessageBase
	// IDs of words to be used in game
	wordIds []uint
	// Index of the start of next batch
	wordIdx uint
	// Index of the current word
	currentWordIdx uint
	// Current round information
	currentRound *Round
}

func CreateGame() *Game {
	return &Game{
		gameState:      InLobby,
		playerMtx:      sync.RWMutex{},
		players:        make(map[string]*Player, 4),
		redTeamCount:   0,
		blueTeamCount:  0,
		messages:       make(chan MessageBase),
		wordIds:        wordStorage.GetShuffledIds(),
		wordIdx:        0,
		currentWordIdx: 0,
		currentRound:   nil,
	}
}

func (g *Game) reset() {
	g.gameState = InLobby
	g.redTeamCount = 0
	g.blueTeamCount = 0
	g.wordIds = wordStorage.GetShuffledIds()
	g.wordIdx = 0
	g.currentWordIdx = 0
	for k := range g.players {
		delete(g.players, k)
	}
}

func (g *Game) run() {
	var err error
	for {
		message := <-g.messages
		switch message := message.(type) {
		case *ChangeNameMessage:
			err = g.changePlayerName(message.PlayerId, message.Name)
		case *ChangeTeamMessage:
			err = g.changePlayerTeam(message.PlayerId, message.Team)
		case *PlayerReadyMessage:
			var allReady bool
			allReady, err = g.changePlayerReadyStatus(message.PlayerId, message.IsReady)
			if allReady {
				g.startRound()
			}
		case *SkipWordMessage:
			err = g.skipCurrentWord(message.PlayerId)
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
			id:      newId,
			conn:    conn,
			name:    name,
			isReady: false,
			team:    -1,
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
		Name: name,
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

	player, exists := g.players[playerId]
	if !exists {
		slog.Error("player not found", "player_id", playerId)
	}
	// delete player from map
	delete(g.players, playerId)

	if len(g.players) == 0 {
		slog.Info("All players have left. Resetting the game.")
		g.reset()
		g.playerMtx.Unlock()
		return
	}

	switch player.team {
	case Red:
		g.redTeamCount--
	case Blue:
		g.blueTeamCount--
	}

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

func (g *Game) changePlayerTeam(playerId string, team Team) error {
	// lock before accessing players
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	// find player
	player, exists := g.players[playerId]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerId)
	}

	if g.gameState != InLobby {
		return fmt.Errorf("game not in lobby state, cannot change team")
	}

	if (team == Red && g.redTeamCount >= 2) || (team == Blue && g.blueTeamCount >= 2) {
		errMsg := &ErrorResponseMessage{
			TypeProperty: TypeProperty{
				Type: ErrorResponseMsg,
			},
			FailedType: ChangeTeamMsg,
			Error:      "Team is full.",
		}
		sendUnicastMessage(player, errMsg)
		return nil
	}

	oldTeam := player.team
	// old team is the same as team to assign, ignore
	if oldTeam == team {
		return nil
	}

	// change team
	player.SetTeam(team)
	if oldTeam == Unassigned {
		if team == Red {
			g.redTeamCount++
		} else {
			g.blueTeamCount++
		}
	} else {
		if team == Red {
			g.redTeamCount++
			g.blueTeamCount--
		} else if team == Blue {
			g.redTeamCount--
			g.blueTeamCount++
		} else {
			g.redTeamCount--
			g.blueTeamCount--
		}
	}

	players := g.GetPlayersCopyUnlocked()
	teamChangedMsg := &TeamChangedMessage{
		TypeProperty: TypeProperty{
			Type: TeamChangedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
		Team: team,
	}
	broadcastMessage(players, teamChangedMsg, nil)
	return nil
}

func (g *Game) changePlayerReadyStatus(playerId string, isReady bool) (bool, error) {
	// lock before accessing players
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	if g.gameState != InLobby {
		return false, fmt.Errorf("game is already running, cannot change ready status")
	}

	// find player
	player, exists := g.players[playerId]
	if !exists {
		return false, fmt.Errorf("player with ID %s not found", playerId)
	}

	if player.team == Unassigned {
		errMsg := &ErrorResponseMessage{
			TypeProperty: TypeProperty{
				Type: ErrorResponseMsg,
			},
			FailedType: ChangeTeamMsg,
			Error:      "Cannot set ready state if player has not chosen a team.",
		}
		sendUnicastMessage(player, errMsg)
		return false, nil
	}

	// change ready status
	player.isReady = isReady
	// create ready status change message
	msg := &PlayerReadyMessage{
		TypeProperty: TypeProperty{
			Type: PlayerReadyMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
		IsReady: isReady,
	}
	// get copy of players to unlock early
	players := g.GetPlayersCopyUnlocked()
	// broadcast ready status change
	broadcastMessage(players, msg, nil)
	allReady := true
	for _, p := range players {
		if !p.isReady {
			allReady = false
			break
		}
	}
	return allReady, nil
}

func (g *Game) startRound() {
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()
	g.gameState = InRound

	// pick words and broadcast to players
	words, err := wordStorage.GetWordsByIds(g.wordIds[g.wordIdx : g.wordIdx+10])
	if err != nil {
		// TODO: handle this situation better, game can start, but words are nowhere to be found
		slog.Error("Failed to get words for round", "err", err)
		return
	}
	g.wordIdx += 10
	// create or update round
	if g.currentRound == nil {
		// first round
		team := Red
		guesserId, hintGiverId := g.selectTeamPlayers(team)
		if guesserId == "" || hintGiverId == "" {
			slog.Error("Not enough players to start round")
			return
		}
		g.currentRound = &Round{
			Team:        team,
			GuesserId:   guesserId,
			HintGiverId: hintGiverId,
			Duration:    RoundDuration,
			Words:       words,
		}
	} else {
		// subsequent rounds, switch teams
		team := g.currentRound.Team.GetOppositeTeam()
		guesserId, hintGiverId := g.selectTeamPlayers(team)
		if guesserId == "" || hintGiverId == "" {
			slog.Error("Not enough players to start round")
			return
		}
		g.currentRound.Team = team
		g.currentRound.GuesserId = guesserId
		g.currentRound.HintGiverId = hintGiverId
		g.currentRound.Words = words
	}
	// generate round start time
	g.currentRound.StartTime = time.Now().UnixMilli() + 5000
	// broadcast round start message
	players := g.GetPlayersCopyUnlocked()
	startRoundMsg := g.currentRound.CreateStartRoundMessage()
	broadcastMessage(players, startRoundMsg, nil)

	// TODO: start end round goroutine with duration
}

func (g *Game) skipCurrentWord(playerId string) error {
	// lock before accessing players
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	if g.gameState != InRound {
		return fmt.Errorf("round is not running, cannot skip word")
	}
	g.currentWordIdx++
	players := g.GetPlayersCopyUnlocked()
	skippedMsg := &WordSkippedMessage{
		TypeProperty: TypeProperty{
			Type: WordSkippedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
	}
	broadcastMessage(players, skippedMsg, nil)
	return nil
}

func (g *Game) selectTeamPlayers(team Team) (string, string) {
	players := make([]string, 0)
	for _, p := range g.players {
		if p.team == team {
			players = append(players, p.id)
		}
	}
	if len(players) < 2 {
		return "", ""
	}
	return players[0], players[1]
}

func (g *Game) CreatePlayerList() []PlayerInfo {
	players := g.GetPlayersCopy()
	playerList := make([]PlayerInfo, 0, len(players))
	for _, p := range players {
		playerList = append(playerList, PlayerInfo{
			Id:      p.id,
			Name:    p.name,
			Team:    p.team,
			IsReady: p.isReady,
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
