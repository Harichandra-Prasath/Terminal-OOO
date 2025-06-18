package core

import (
	"fmt"
	"net/http"
)

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
		room.AddPlayer(host)
		fmt.Printf("New Room added to server with ID '%s' by host '%s'\n", room.ID.String(), host.Name)
		writeGoodResponse(http.StatusCreated, w, "Room Created Successfully", map[string]any{
			"room_id": room.ID,
			"host_id": host.Id.String(),
		})
	}
}

func (S *Server) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomRequest

	if parseData(r, w, &req) {
		newPlayer := CreateNewPlayer(req.PlayerName)
		room := S.GetRoom(req.RoomId)
		if room == nil {
			writeBadResponse(http.StatusNotFound, w, "Room not found on the server")
		}

		if room.Status == "STARTED" {
			writeBadResponse(http.StatusBadRequest, w, "Room is already started")
		}

		if len(room.Players) >= room.Cfg.MaxPlayers {
			writeBadResponse(http.StatusBadRequest, w, "Max Players Reached")
		}

		room.AddPlayer(newPlayer)
		fmt.Printf("Player '%s' added to the Room '%s'\n", req.PlayerName, room.ID.String())
		writeGoodResponse(http.StatusCreated, w, "Player added to the room", map[string]any{
			"player_id": newPlayer.Id.String(),
		})
	}

}

func (S *Server) StartGameHandler(w http.ResponseWriter, r *http.Request) {

}
