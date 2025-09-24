package main

import "fmt"

type MessageType string

const (
	ConnectMsg    	MessageType = "connect"
	ConnectAckMsg 	MessageType = "connect_ack"
	ChangeNameMsg 	MessageType = "change_name"
	NameChangedMsg 	MessageType = "name_changed"
	PlayerJoinedMsg MessageType = "player_joined"
	PlayerListMsg  	MessageType = "player_list"
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
}

type ConnectAckMessage struct {
	TypeProperty
	PlayerIdProperty
}

type ChangeNameMessage struct {
	TypeProperty
	PlayerIdProperty
	Name     string `json:"name"`
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

type PlayerListMessage struct {
	TypeProperty
	Players []PlayerInfo `json:"players"`
}

func ConstructMessageContainer(messageType MessageType) (MessageBase, error) {
	switch messageType {
	case ConnectMsg:
		return &ConnectMessage{}, nil
	case ChangeNameMsg:
		return &ChangeNameMessage{}, nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}
