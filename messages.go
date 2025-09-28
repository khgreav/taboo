package main

import "fmt"

type MessageType string

const (
	ConnectMsg      MessageType = "connect"
	ConnectAckMsg   MessageType = "connect_ack"
	ChangeNameMsg   MessageType = "change_name"
	NameChangedMsg  MessageType = "name_changed"
	PlayerJoinedMsg MessageType = "player_joined"
	PlayerLeftMsg   MessageType = "player_left"
	PlayerListMsg   MessageType = "player_list"
	PlayerReadyMsg  MessageType = "player_ready"
	WordListMsg     MessageType = "word_list"
	SkipWordMsg     MessageType = "skip_word"
	WordSkippedMsg  MessageType = "word_skipped"
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

type ConnectMessage struct {
	TypeProperty
	PlayerId *string `json:"playerId"`
	Name     string  `json:"name"`
}

type ConnectAckMessage struct {
	TypeProperty
	PlayerIdProperty
	Name string `json:"name"`
}

type ChangeNameMessage struct {
	TypeProperty
	PlayerIdProperty
	Name string `json:"name"`
}

type NameChangedMessage struct {
	TypeProperty
	PlayerIdProperty
	Name string `json:"name"`
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

type PlayerListMessage struct {
	TypeProperty
	Players []PlayerInfo `json:"players"`
}

type PlayerReadyMessage struct {
	TypeProperty
	PlayerIdProperty
	IsReady bool `json:"isReady"`
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

func ConstructMessageContainer(messageType MessageType) (MessageBase, error) {
	switch messageType {
	case ConnectMsg:
		return &ConnectMessage{}, nil
	case ChangeNameMsg:
		return &ChangeNameMessage{}, nil
	case PlayerReadyMsg:
		return &PlayerReadyMessage{}, nil
	case SkipWordMsg:
		return &SkipWordMessage{}, nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}
