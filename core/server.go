package core

import (
	"fmt"
	"net/http"
	"sync"
)

type ServerCfg struct {
	Port int
}

type Server struct {
	Lock   sync.RWMutex
	Rooms  map[string]*Room
	server *http.Server
	Cfg    ServerCfg
}

func (S *Server) registerRoutes(handler http.Handler) {
	m := handler.(*http.ServeMux)
	m.HandleFunc("POST /create-room", S.CreateRoomHandler)
	m.HandleFunc("POST /join-room", S.JoinRoomHandler)
	m.HandleFunc("POST /start-room", S.StartGameHandler)
	m.HandleFunc("GET /init/{roomId}", S.InitialiseHandler)
}

func (S *Server) AddRoom(room *Room) {
	S.Lock.Lock()
	defer S.Lock.Unlock()

	S.Rooms[room.ID] = room

}

func (S *Server) GetRoom(id string) *Room {
	S.Lock.RLock()
	defer S.Lock.RUnlock()

	return S.Rooms[id]
}

func NewServer(Cfg *ServerCfg) *Server {

	return &Server{
		Rooms:  make(map[string]*Room),
		Cfg:    *Cfg,
		server: &http.Server{Handler: http.NewServeMux()},
	}
}

func (S *Server) Start() error {

	S.registerRoutes(S.server.Handler)
	S.server.Addr = fmt.Sprintf(":%d", S.Cfg.Port)
	fmt.Println("Server Starting on Port:", S.Cfg.Port)

	return S.server.ListenAndServe()
}
