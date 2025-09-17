package main

import "fmt"

type MessageType string

const (
	ConnectMsg    MessageType = "connect"
	ConnectAckMsg MessageType = "connect_ack"
	SetNameMsg    MessageType = "set_name"
)

type MessageBase interface {
	GetType() MessageType
}

type TypeMessage struct {
	Type MessageType `json:"type"`
}

func (message TypeMessage) GetType() MessageType {
	return message.Type
}

type ConnectMessage struct {
	TypeMessage
	PlayerId *string `json:"playerId"`
}

type ConnectAckMessage struct {
	TypeMessage
	PlayerId string `json:"playerId"`
}

type SetNameMessage struct {
	TypeMessage
	PlayerId string `json:"playerId"`
	Name     string `json:"name"`
}

func ConstructMessageContainer(messageType MessageType) (MessageBase, error) {
	switch messageType {
	case ConnectMsg:
		return &ConnectMessage{}, nil
	case ConnectAckMsg:
		return &ConnectAckMessage{}, nil
	case SetNameMsg:
		return &SetNameMessage{}, nil
	default:
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}
}
