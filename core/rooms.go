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
	VoteStore    map[uuid.UUID]int
}

type Message struct {
	PlayerId uuid.UUID `json:"player_id"`
	Action   string    `json:"action"`
	Value    string    `json:"value"`
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
		VoteStore:    make(map[uuid.UUID]int),
	}

}

func (R *Room) Start() {

	fmt.Println("Room main loop Started")
	for {

		select {
		case m := <-R.InboundChan:
			switch m.Action {
			case "START":

			case "HINT":
				go R.handleHintMessages(m)
			case "VOTE":
				go R.handleVoteMessages(m)
			}
		case m := <-R.OutboundChan:
			go R.WritetoPlayers(m)
		}

	}

}

// Forward the hint message
func (R *Room) handleHintMessages(m *Message) {

	hinter := R.GetPlayer(m.PlayerId)
	fmt.Printf("Player '%s' gave the hint '%s'\n", hinter.Name, m.Value)

	R.OutboundChan <- m
}

// Store the vote messages
func (R *Room) handleVoteMessages(m *Message) {

	Voter := R.GetPlayer(m.PlayerId).Name

	// No vote
	if m.Value == "" {
		fmt.Printf("Player '%s' didnt Vote\n", Voter)
		return
	}

	playerId, err := uuid.Parse(m.Value)
	if err != nil {
		fmt.Printf("Error in voting player with Id '%s': %s\n", m.Value, err.Error())
	}

	R.Lock.Lock()
	R.VoteStore[playerId] += 1
	R.Lock.Unlock()

	Voted := R.GetPlayer(playerId).Name

	fmt.Printf("Player '%s' voted '%s'\n", Voter, Voted)

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
