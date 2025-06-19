package core

import (
	"fmt"
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
	Lock         sync.RWMutex
	ID           uuid.UUID
	Host         *Player
	Players      []*Player
	Cfg          *RoomConfig
	Status       string
	InboundChan  chan *Message
	OutboundChan chan *Message
}

type Message struct {
	PlayerName string
	Action     string
	Value      string
}

func CreateNewRoom(Cfg *RoomConfig, host *Player) *Room {

	return &Room{
		ID:           uuid.New(),
		Host:         host,
		Players:      []*Player{host},
		Cfg:          Cfg,
		Status:       "YET_TO_START",
		InboundChan:  make(chan *Message),
		OutboundChan: make(chan *Message),
	}

}

func (R *Room) Start() {

	fmt.Println("Room main loop Started")
	for {

		select {
		case m := <-R.InboundChan:
			if m.Action == "START" {
				// Start the Reading
				go R.ListenfromPlayers()
			}
		case m := <-R.OutboundChan:
			go R.WritetoPlayers(m)
		}

	}

}

func (R *Room) ListenfromPlayers() {

	for _, player := range R.Players {
		go R.ListenfromPlayer(player)
	}
}

func (R *Room) ListenfromPlayer(player *Player) {

	fmt.Printf("Started Listening from the Player '%s'\n", player.Name)

	for {
		var message Message
		err := player.WsConn.ReadJSON(&message)
		if err != nil {
			fmt.Printf("Error reading from player '%s'\n", player.Name)
			continue
		}
		R.InboundChan <- &message
	}

}

func (R *Room) WritetoPlayers(message *Message) {

	for _, player := range R.Players {
		conn := player.WsConn
		err := conn.WriteJSON(&message)
		if err != nil {
			fmt.Printf("Error writing to the player '%s': '%s'\n", player.Name, err.Error())
		}
	}

}

func (R *Room) AddPlayer(player *Player) {

	R.Lock.Lock()
	defer R.Lock.Unlock()

	R.Players = append(R.Players, player)
}

// Use Linear scan as the size will be small
func (R *Room) GetPlayer(id uuid.UUID) *Player {
	R.Lock.RLock()
	defer R.Lock.RUnlock()

	for _, player := range R.Players {
		if id == player.Id {
			return player
		}
	}

	return nil
}
