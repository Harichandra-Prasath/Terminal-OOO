package core

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id     uuid.UUID
	Name   string
	WsConn *websocket.Conn
	Alive  bool
	Liar   bool
}

func CreateNewPlayer(name string) *Player {
	return &Player{
		Id:    uuid.New(),
		Name:  name,
		Alive: true,
		Liar:  false,
	}
}
