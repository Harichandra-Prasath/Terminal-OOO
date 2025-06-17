package core

import "github.com/google/uuid"

type Player struct {
	Id   uuid.UUID
	Name string
}

func CreateNewPlayer(name string) *Player {
	return &Player{
		Id:   uuid.New(),
		Name: name,
	}
}
