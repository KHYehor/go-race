package server

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	ctx        context.Context
	httpServer *http.Server
	wsHandler  *WS
}

func (s *Server) setRoutes() {
	mux := http.NewServeMux()

	// /cars handler
	mux.HandleFunc("/cars", s.wsHandler.CarsConnHandler)

	s.httpServer.Handler = mux
}

func New() *Server {
	return &Server{
		ctx:       context.Background(),
		wsHandler: NewWS(context.Background()),
		httpServer: &http.Server{
			Addr: ":8080",
		},
	}
}

func (s *Server) Start() {
	s.setRoutes()
	log.Printf("Listening...")
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Printf("Error happened in ListenAndServe: %v", err)
	}
}
