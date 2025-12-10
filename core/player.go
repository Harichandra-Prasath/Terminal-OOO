package core

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id     string
	Name   string
	WsConn *websocket.Conn
	Alive  bool
	Liar   bool
}

func CreateNewPlayer(name string) *Player {
	return &Player{
		Id:    uuid.New().String()[:7],
		Name:  name,
		Alive: true,
		Liar:  false,
	}
}
