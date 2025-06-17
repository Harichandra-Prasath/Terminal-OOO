package core

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ServerCfg struct {
	Port int
}

type Server struct {
	Rooms  map[uuid.UUID]*Room
	server *http.Server
	Cfg    ServerCfg
}

func registerRoutes(m *http.ServeMux) {
}

func NewServer(Cfg *ServerCfg) *Server {

	mux := http.NewServeMux()
	registerRoutes(mux)
	return &Server{
		Rooms:  make(map[uuid.UUID]*Room),
		Cfg:    *Cfg,
		server: &http.Server{Handler: mux},
	}
}

func (S *Server) Start() error {

	S.server.Addr = fmt.Sprintf(":%d", S.Cfg.Port)
	fmt.Println("Server Starting on Port:", S.Cfg.Port)

	return S.server.ListenAndServe()
}
