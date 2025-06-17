package core

import "github.com/google/uuid"

type RoomConfig struct {
	MaxPlayers int
	MinPlayers int
	RoundTime  int
	Topic      string
}

type Room struct {
	ID      uuid.UUID
	Host    *Player
	Players []*Player
	Cfg     *RoomConfig
}

func CreateNewRoom(Cfg *RoomConfig, host *Player) *Room {

	return &Room{
		ID:      uuid.New(),
		Host:    host,
		Players: []*Player{host},
		Cfg:     Cfg,
	}

}
