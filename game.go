package main

import (
	"context"
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
	// Player IDs per team
	teamPlayers map[Team][]string
	// Team scores
	teamScores map[Team]int
	// Channel for incoming player messages
	messages chan MessageBase
	// IDs of words to be used in game
	wordIds []uint
	// Index of the start of next batch
	wordOffset uint
	// Index of the currently guessed word
	currentWordIdx uint
	// Current round number
	roundNumber uint
	// Current round information
	currentRound *Round
	// Round cancel context
	roundCtx context.Context
	// Round cancel function
	roundCancel context.CancelFunc
}

func CreateGame() *Game {
	return &Game{
		gameState:      InLobby,
		playerMtx:      sync.RWMutex{},
		players:        make(map[string]*Player, 4),
		teamPlayers:    make(map[Team][]string),
		teamScores:     make(map[Team]int),
		messages:       make(chan MessageBase),
		wordIds:        wordStorage.GetShuffledIds(),
		wordOffset:     0,
		currentWordIdx: 0,
		roundNumber:    0,
		currentRound:   nil,
		roundCtx:       nil,
		roundCancel:    nil,
	}
}

func (g *Game) reset() {
	if g.roundCancel != nil {
		g.roundCancel()
		g.roundCancel = nil
	}
	g.gameState = InLobby
	g.teamPlayers[Red] = []string{}
	g.teamPlayers[Blue] = []string{}
	g.teamScores[Red] = 0
	g.teamScores[Blue] = 0
	g.wordIds = wordStorage.GetShuffledIds()
	g.wordOffset = 0
	g.currentWordIdx = 0
	g.roundNumber = 0
	g.currentRound = nil
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
				g.prepareRound()
			}
		case *StartRoundMessage:
			g.startRound(message.PlayerId)
		case *SkipWordMessage:
			err = g.skipCurrentWord(message.PlayerId)
		case *GuessWordMessage:
			err = g.guessWord(message.PlayerId)
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

	// remove player from team player list
	for i, v := range g.teamPlayers[player.team] {
		if v == playerId {
			g.teamPlayers[player.team] = append(g.teamPlayers[player.team][:i], g.teamPlayers[player.team][i+1:]...)
			break
		}
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
		g.playerMtx.Unlock()
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

	if len(g.teamPlayers[team]) >= MaxTeamMembers {
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
		g.teamPlayers[team] = append(g.teamPlayers[team], playerId)
		slog.Debug("Player assigned to team", "player_id", playerId, "team", team, "teamplayers", g.teamPlayers)
	} else {
		for i, v := range g.teamPlayers[oldTeam] {
			if v == playerId {
				g.teamPlayers[oldTeam] = append(g.teamPlayers[oldTeam][:i], g.teamPlayers[oldTeam][i+1:]...)
				break
			}
		}
		if team != Unassigned {
			g.teamPlayers[team] = append(g.teamPlayers[team], playerId)
		}
		slog.Debug("Player changed team", "player_id", playerId, "old_team", oldTeam, "new_team", team, "teamplayers", g.teamPlayers)
	}
	player.SetReady(false)

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
		sendErrorMessage(g.players[playerId], ChangeTeamMsg, "Game is already running, cannot change ready status.")
		return false, fmt.Errorf("game is already running, cannot change ready status")
	}

	// find player
	player, exists := g.players[playerId]
	if !exists {
		return false, fmt.Errorf("player with ID %s not found", playerId)
	}

	if player.team == Unassigned {
		sendErrorMessage(player, ChangeTeamMsg, "Cannot set ready state if player has not chosen a team.")
		return false, fmt.Errorf("cannot set ready state if player has not chosen a team")
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

func (g *Game) prepareRound() {
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	// pick words and broadcast to players
	words, err := wordStorage.GetWordsByIds(g.wordIds[g.wordOffset : g.wordOffset+10])
	if err != nil {
		// TODO: handle this situation better, game can start, but words are nowhere to be found
		slog.Error("Failed to get words for round", "err", err)
		return
	}
	g.wordOffset += 10
	// select team and players for round
	team := selectTeam(g.roundNumber)
	if len(g.teamPlayers[team]) < 2 {
		slog.Error("Not enough players to start the round")
		return
	}

	hintGiverId, guesserId := g.selectTeamPlayers(team, g.roundNumber)
	// create round object
	g.currentRound = &Round{
		Team:        team,
		GuesserId:   guesserId,
		HintGiverId: hintGiverId,
		Duration:    RoundDuration,
		Words:       words,
	}

	// broadcast round prepare message
	players := g.GetPlayersCopyUnlocked()
	startRoundMsg := g.currentRound.CreateRoundSetupMessage()
	err = broadcastMessage(players, startRoundMsg, nil)
	if err != nil {
		slog.Error("Failed to broadcast round setup message", "err", err)
		return
	}

	g.gameState = InProgress
}

func (g *Game) startRound(playerId string) {
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	if g.gameState != InProgress {
		sendErrorMessage(g.players[playerId], StartRoundMsg, "Cannot start round, game not in progress state.")
		slog.Warn("Cannot start round, game not in progress state.")
		return
	}

	if playerId != g.currentRound.HintGiverId {
		sendErrorMessage(g.players[playerId], StartRoundMsg, "Only the hint giver can start the round.")
		slog.Warn("Only the hint giver can start the round.", "player_id", playerId)
		return
	}

	players := g.GetPlayersCopyUnlocked()
	roundStartedMsg := g.currentRound.CreateRoundStartedMessage()
	err := broadcastMessage(players, roundStartedMsg, nil)
	if err != nil {
		slog.Error("Failed to broadcast round started message", "err", err)
		return
	}

	g.gameState = InRound

	g.roundCtx, g.roundCancel = context.WithCancel(context.Background())
	go func(ctx context.Context, duration int) {
		select {
		case <-time.After(time.Duration(duration) * time.Second):
			g.endRound()
		case <-ctx.Done():
			slog.Info("Round cancelled")
			return
		}
	}(g.roundCtx, g.currentRound.Duration)
}

func (g *Game) endRound() {
	g.playerMtx.Lock()

	if g.gameState != InRound {
		slog.Warn("Cannot end round, game not in round state")
		g.playerMtx.Unlock()
		return
	}

	g.gameState = InProgress
	g.roundNumber++
	players := g.GetPlayersCopyUnlocked()
	endRoundMsg := g.currentRound.CreateRoundEndedMessage()
	err := broadcastMessage(players, endRoundMsg, nil)
	if err != nil {
		slog.Error("Failed to broadcast round ended message", "err", err)
		return
	}
	g.playerMtx.Unlock()
	g.prepareRound()
}

func (g *Game) guessWord(playerId string) error {
	// lock before accessing players
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	player := g.players[playerId]

	if g.gameState != InRound {
		sendErrorMessage(player, SkipWordMsg, "Round is not running, cannot guess word.")
		return fmt.Errorf("round is not running, cannot guess word")
	}

	if g.currentRound.HintGiverId != playerId {
		sendErrorMessage(player, SkipWordMsg, "Only the hint giver can mark a guess.")
		return fmt.Errorf("only the hin giver can mark a guess")
	}

	g.teamScores[player.team]++
	g.currentWordIdx++

	players := g.GetPlayersCopyUnlocked()
	guessedMsg := &WordGuessedMessage{
		TypeProperty: TypeProperty{
			Type: WordGuessedMsg,
		},
		PlayerIdProperty: PlayerIdProperty{
			PlayerId: playerId,
		},
		RedScore:  g.teamScores[Red],
		BlueScore: g.teamScores[Blue],
	}
	broadcastMessage(players, guessedMsg, nil)

	if (g.wordOffset - g.currentWordIdx) <= 5 {
		// pick words and broadcast to players
		words, err := wordStorage.GetWordsByIds(g.wordIds[g.wordOffset : g.wordOffset+10])
		if err != nil {
			slog.Error("Failed to get new batch of words", "err", err)
			return fmt.Errorf("failed to get new batch of words: %w", err)
		}
		g.wordOffset += 10

		wordListMsg := &WordListMessage{
			TypeProperty: TypeProperty{
				Type: WordListMsg,
			},
			Words: words,
		}

		err = broadcastMessage(players, wordListMsg, nil)
		if err != nil {
			slog.Error("Failed to broadcast word list message", "err", err)
			return fmt.Errorf("failed to broadcast word list message: %w", err)
		}
	}

	return nil
}

func (g *Game) skipCurrentWord(playerId string) error {
	// lock before accessing players
	g.playerMtx.Lock()
	defer g.playerMtx.Unlock()

	if g.gameState != InRound {
		sendErrorMessage(g.players[playerId], SkipWordMsg, "Round is not running, cannot skip word.")
		return fmt.Errorf("round is not running, cannot skip word")
	}

	if g.currentRound.HintGiverId != playerId {
		sendErrorMessage(g.players[playerId], SkipWordMsg, "Only the hint giver can skip.")
		return fmt.Errorf("only the hint giver can skip words")
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

	if (g.wordOffset - g.currentWordIdx) <= 5 {
		// pick words and broadcast to players
		words, err := wordStorage.GetWordsByIds(g.wordIds[g.wordOffset : g.wordOffset+10])
		if err != nil {
			slog.Error("Failed to get new batch of words", "err", err)
			return fmt.Errorf("failed to get new batch of words: %w", err)
		}
		g.wordOffset += 10

		wordListMsg := &WordListMessage{
			TypeProperty: TypeProperty{
				Type: WordListMsg,
			},
			Words: words,
		}

		err = broadcastMessage(players, wordListMsg, nil)
		if err != nil {
			slog.Error("Failed to broadcast word list message", "err", err)
			return fmt.Errorf("failed to broadcast word list message: %w", err)
		}
	}

	return nil
}

func selectTeam(round uint) Team {
	team := Red
	if round%2 == 1 {
		team = Blue
	}
	return team
}

func (g *Game) selectTeamPlayers(team Team, round uint) (string, string) {
	if (round/2)%2 == 0 {
		return g.teamPlayers[team][0], g.teamPlayers[team][1]
	}
	return g.teamPlayers[team][1], g.teamPlayers[team][0]
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
