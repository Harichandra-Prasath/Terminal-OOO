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
			"id":   room.ID,
			"host": host.Name,
		})
	}
}

func (S *Server) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
}

func (S *Server) StartGameHandler(w http.ResponseWriter, r *http.Request) {
}
