package core

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id     uuid.UUID
	Name   string
	WsConn *websocket.Conn
}

func CreateNewPlayer(name string) *Player {
	return &Player{
		Id:   uuid.New(),
		Name: name,
	}
}
