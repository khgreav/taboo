package main

import "fmt"

type MessageType string

const (
	// general messages
	ErrorResponseMsg MessageType = "error_response"
	// player connections
	ConnectMsg            MessageType = "connect"
	ConnectAckMsg         MessageType = "connect_ack"
	ReconnectAckMsg       MessageType = "reconnect_ack"
	PlayerJoinedMsg       MessageType = "player_joined"
	PlayerLeftMsg         MessageType = "player_left"
	PlayerDisconnectedMsg MessageType = "player_disconnected"
	PlayerReconnectedMsg  MessageType = "player_reconnected"
	// lobby state
	PlayerListMsg       MessageType = "player_list"
	ChangeTeamMsg       MessageType = "change_team"
	TeamChangedMsg      MessageType = "team_changed"
	PlayerReadyMsg      MessageType = "player_ready"
	GameStateChangedMsg MessageType = "game_state_changed"
	// game rounds
	RoundSetupMsg   MessageType = "round_setup"
	StartRoundMsg   MessageType = "start_round"
	RoundStartedMsg MessageType = "round_started"
	RoundEndedMsg   MessageType = "round_ended"
	PauseRoundMsg   MessageType = "pause_round"
	RoundPausedMsg  MessageType = "round_paused"
	ResumeRoundMsg  MessageType = "resume_round"
	RoundResumedMsg MessageType = "round_resumed"
	GameEndedMsg    MessageType = "game_ended"
	ResetGameMsg    MessageType = "reset_game"
	GameResetMsg    MessageType = "game_reset"
	// round actions
	SkipWordMsg    MessageType = "skip_word"
	WordSkippedMsg MessageType = "word_skipped"
	GuessWordMsg   MessageType = "guess_word"
	WordGuessedMsg MessageType = "word_guessed"
	WordListMsg    MessageType = "word_list"
)

type MessageBase interface {
	GetType() MessageType
}

type TypeProperty struct {
	Type MessageType `json:"type"`
}

func (prop TypeProperty) GetType() MessageType {
	return prop.Type
}

type PlayerIdProperty struct {
	PlayerId string `json:"playerId"`
}

func (prop PlayerIdProperty) GetPlayerId() string {
	return prop.PlayerId
}

type PlayerMessage interface {
	MessageBase
	GetPlayerId() string
}

type ErrorResponseMessage struct {
	TypeProperty
	FailedType MessageType `json:"failedType"`
	Error      string      `json:"error"`
	ErrorCode  ErrorCode   `json:"errorCode"`
}

type ConnectMessage struct {
	TypeProperty
	PlayerId     *string `json:"playerId"`
	SessionToken *string `json:"sessionToken"`
	Name         string  `json:"name"`
}

type ConnectAckMessage struct {
	TypeProperty
	PlayerIdProperty
	SessionToken string `json:"sessionToken"`
	Name         string `json:"name"`
}

type ReconnectAckMessage struct {
	ConnectAckMessage
	Team              Team         `json:"team"`
	State             GameState    `json:"state"`
	RemainingDuration int          `json:"remainingDuration"`
	CurrentTeam       Team         `json:"currentTeam"`
	GuesserId         string       `json:"guesserId"`
	HintGiverId       string       `json:"hintGiverId"`
	RedScore          int          `json:"redScore"`
	BlueScore         int          `json:"blueScore"`
	Words             []*TabooWord `json:"words"`
}

type PlayerJoinedMessage struct {
	TypeProperty
	PlayerIdProperty
	Name string `json:"name"`
}

type PlayerLeftMessage struct {
	TypeProperty
	PlayerIdProperty
}

type PlayerDisconnectedMessage struct {
	TypeProperty
	PlayerIdProperty
}

type PlayerReconnectedMessage struct {
	TypeProperty
	PlayerIdProperty
}

type PlayerListMessage struct {
	TypeProperty
	Players []PlayerInfo `json:"players"`
}

type ChangeTeamMessage struct {
	TypeProperty
	PlayerIdProperty
	Team Team `json:"team"`
}

type TeamChangedMessage struct {
	TypeProperty
	PlayerIdProperty
	Team Team `json:"team"`
}

type PlayerReadyMessage struct {
	TypeProperty
	PlayerIdProperty
	IsReady bool `json:"isReady"`
}

type GameStateChangedMessage struct {
	TypeProperty
	State GameState `json:"state"`
}

type WordListMessage struct {
	TypeProperty
	Words []*TabooWord `json:"words"`
}

type SkipWordMessage struct {
	TypeProperty
	PlayerIdProperty
}

type WordSkippedMessage struct {
	TypeProperty
	PlayerIdProperty
}

type GuessWordMessage struct {
	TypeProperty
	PlayerIdProperty
}

type WordGuessedMessage struct {
	TypeProperty
	PlayerIdProperty
	RedScore  int `json:"redScore"`
	BlueScore int `json:"blueScore"`
}

type RoundSetupMessage struct {
	TypeProperty
	Team        Team         `json:"team"`
	GuesserId   string       `json:"guesserId"`
	HintGiverId string       `json:"hintGiverId"`
	Duration    int          `json:"duration"`
	Words       []*TabooWord `json:"words"`
}

type StartRoundMessage struct {
	TypeProperty
	PlayerIdProperty
}

type RoundStartedMessage struct {
	TypeProperty
	PlayerIdProperty
}

type RoundEndedMessage struct {
	TypeProperty
}

type RoundPausedMessage struct {
	TypeProperty
	RemainingDuration int `json:"remainingDuration"`
}

type ResumeRoundMessage struct {
	TypeProperty
	PlayerIdProperty
}

type RoundResumedMessage struct {
	TypeProperty
	PlayerIdProperty
}

type GameEndedMessage struct {
	TypeProperty
	RedScore  int `json:"redScore"`
	BlueScore int `json:"blueScore"`
}

type ResetGameMessage struct {
	TypeProperty
	PlayerIdProperty
}

type GameResetMessage struct {
	TypeProperty
	Players []PlayerInfo `json:"players"`
}

func ConstructMessageContainer(messageType MessageType) (MessageBase, error) {
	switch messageType {
	case ConnectMsg:
		return &ConnectMessage{}, nil
	case ChangeTeamMsg:
		return &ChangeTeamMessage{}, nil
	case PlayerReadyMsg:
		return &PlayerReadyMessage{}, nil
	case StartRoundMsg:
		return &StartRoundMessage{}, nil
	case SkipWordMsg:
		return &SkipWordMessage{}, nil
	case GuessWordMsg:
		return &GuessWordMessage{}, nil
	case ResumeRoundMsg:
		return &ResumeRoundMessage{}, nil
	case ResetGameMsg:
		return &ResetGameMessage{}, nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}
