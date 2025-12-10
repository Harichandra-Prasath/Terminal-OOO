package core

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (S *Server) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if parseData(r, w, &req) {

		// Create a config
		host := CreateNewPlayer(req.PlayerName)

		room := CreateNewRoom(&RoomConfig{
			MaxPlayers: req.MaxSize,
			MinPlayers: req.MinSize,
			RoundTime:  req.RoundTime,
		}, host)

		S.AddRoom(room)
		fmt.Printf("New Room added to server with ID '%s' by host '%s'\n", room.ID, host.Name)
		go room.Start()
		writeGoodResponse(http.StatusCreated, w, "Room Created Successfully", map[string]any{
			"room_id": room.ID,
			"host_id": host.Id,
		})
	}
}

func (S *Server) InitialiseHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	roomId := r.PathValue("roomId")
	playerId := queryParams.Get("playerId")

	room := S.GetRoom(roomId)
	if room == nil {
		writeBadResponse(http.StatusNotFound, w, "Room not found on the server")
		return
	}

	player := room.GetPlayer(playerId)
	if player == nil {
		writeBadResponse(http.StatusNotFound, w, "Player not on the room")
		return
	}

	// Upgrade the connection to WS
	NewConn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error in Upgrading the connection: " + err.Error())
		return
	}

	player.WsConn = NewConn
	fmt.Printf("Player '%s' initialised his WS Connection\n", player.Name)
	go room.ListenfromPlayer(player)
}

func (S *Server) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomRequest

	if parseData(r, w, &req) {
		newPlayer := CreateNewPlayer(req.PlayerName)
		room := S.GetRoom(req.RoomId)
		if room == nil {
			writeBadResponse(http.StatusNotFound, w, "Room not found on the server")
			return
		}

		if room.Status == "STARTED" {
			writeBadResponse(http.StatusBadRequest, w, "Room is already started")
			return
		}

		if len(room.Players) >= room.Cfg.MaxPlayers {
			writeBadResponse(http.StatusBadRequest, w, "Max Players Reached")
			return
		}

		room.AddPlayer(newPlayer)
		fmt.Printf("Player '%s' added to the Room '%s'\n", req.PlayerName, room.ID)
		writeGoodResponse(http.StatusCreated, w, "Player added to the room", map[string]any{
			"player_id": newPlayer.Id,
		})
	}

}

func (S *Server) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	var req StartRoomRequest

	if parseData(r, w, &req) {
		room := S.GetRoom(req.RoomId)
		if room == nil {
			writeBadResponse(http.StatusNotFound, w, "Room not found on the server")
			return
		}

		if room.Host.Id != req.HostId {
			writeBadResponse(http.StatusUnauthorized, w, "Acess denied")
			return
		}

		if len(room.Players) < room.Cfg.MinPlayers {
			writeBadResponse(http.StatusBadRequest, w, "Not enough players to start")
			return
		}

		if room.Status == "STARTED" {
			writeBadResponse(http.StatusBadRequest, w, "room already started")
			return
		}

		// Update the room status
		room.Status = "STARTED"
		fmt.Printf("Room with ID '%s' Started\n", req.HostId)

		// send the action messages
		room.InboundChan <- &Message{Action: "START"}
		writeGoodResponse(http.StatusCreated, w, "Room Started", map[string]any{})
	}
}
