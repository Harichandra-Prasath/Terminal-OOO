package core

import (
	"sync"

	"github.com/google/uuid"
)

type RoomConfig struct {
	MaxPlayers int
	MinPlayers int
	RoundTime  int
	Topic      string
}

type Room struct {
	Lock    sync.RWMutex
	ID      uuid.UUID
	Host    *Player
	Players []*Player
	Cfg     *RoomConfig
	Status  string
}

func CreateNewRoom(Cfg *RoomConfig, host *Player) *Room {

	return &Room{
		ID:      uuid.New(),
		Host:    host,
		Players: []*Player{host},
		Cfg:     Cfg,
		Status:  "YET_TO_START",
	}

}

func (R *Room) AddPlayer(player *Player) {

	R.Lock.Lock()
	defer R.Lock.Unlock()

	R.Players = append(R.Players, player)
}
