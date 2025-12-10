package core

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

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
	ID           string
	Host         *Player
	Players      []*Player
	Cfg          *RoomConfig
	Status       string
	InboundChan  chan *Message
	OutboundChan chan *Message
	VoteStore    map[string]int
}

type Message struct {
	PlayerId string `json:"player_id"`
	Action   string `json:"action"`
	Value    string `json:"value"`
}

func CreateNewRoom(Cfg *RoomConfig, host *Player) *Room {

	return &Room{
		ID:           uuid.New().String()[:7],
		Host:         host,
		Players:      []*Player{host},
		Cfg:          Cfg,
		Status:       "YET_TO_START",
		InboundChan:  make(chan *Message),
		OutboundChan: make(chan *Message),
		VoteStore:    make(map[string]int),
	}

}

func (R *Room) Start() {

	fmt.Println("Room main loop Started")
	for {

		select {
		case m := <-R.InboundChan:
			switch m.Action {
			case "START":
				go R.initGame()
			case "ROUND":
				go R.startRound()
			case "HINT":
				go R.handleHintMessages(m)
			case "VOTE":
				go R.handleVoteMessages(m)
			}
		case m := <-R.OutboundChan:

			// Start the Voting
			if m.Action == "VOTE_START" {
				go R.startVotingSession()
			}
			go R.WritetoPlayers(m)
		}

	}

}

func (R *Room) initGame() {

	// Pick a player who should be liar
	playerIndex := rand.IntN(len(R.Players))
	Liar := R.Players[playerIndex]
	Liar.Liar = true

	// Dummy animal let it be Tiger
	go R.WritetoPlayers(&Message{
		Action: "INIT",
		Value:  "Tiger",
	}, playerIndex)

	go R.WritetoPlayer(&Message{
		Action: "INIT",
		Value:  "Liar",
	}, Liar)

}

func (R *Room) startVotingSession() {
	ticker := time.NewTicker(10 + 1*time.Second)
	for {
		select {
		case <-ticker.C:
			R.OutboundChan <- &Message{
				Action: "VOTE_END",
			}

		}

	}
}

// Start the game timer
func (R *Room) startRound() {
	ticker := time.NewTicker(time.Duration(R.Cfg.RoundTime+1) * time.Second)
	for {
		select {
		case <-ticker.C:
			R.OutboundChan <- &Message{
				Action: "VOTE_START",
			}
			return

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

	playerId := m.Value

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

func (R *Room) WritetoPlayers(message *Message, exceptions ...int) {

	Store := make(map[int]struct{})

	for _, exceptionIndex := range exceptions {
		Store[exceptionIndex] = struct{}{}
	}

	for index, player := range R.Players {
		if _, ok := Store[index]; ok {
			continue
		}
		err := R.WritetoPlayer(message, player)
		if err != nil {
			fmt.Printf("Error writing to the player '%s': '%s'\n", player.Name, err.Error())
		}
	}

}

func (R *Room) WritetoPlayer(message *Message, player *Player) error {

	conn := player.WsConn
	err := conn.WriteJSON(&message)
	if err != nil {
		return fmt.Errorf("writing to the player: '%s'", err.Error())
	}
	return nil
}

func (R *Room) AddPlayer(player *Player) {

	R.Lock.Lock()
	defer R.Lock.Unlock()

	R.Players = append(R.Players, player)
}

// Use Linear scan as the size will be small
func (R *Room) GetPlayer(id string) *Player {
	R.Lock.RLock()
	defer R.Lock.RUnlock()

	for _, player := range R.Players {
		if id == player.Id {
			return player
		}
	}

	return nil
}
