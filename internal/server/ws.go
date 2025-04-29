package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"go-race/internal/race"
	"log"
	"net/http"
)

type WS struct {
	ctx      context.Context
	upgrader websocket.Upgrader
}

func NewWS(ctx context.Context) *WS {
	return &WS{
		ctx:      ctx,
		upgrader: websocket.Upgrader{},
	}
}

// CarsConnHandler Listen for cars connections
func (ws *WS) CarsConnHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(ws.ctx)
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade websocket error: %v", err)
		return
	}
	defer conn.Close()

	// Init car tracker and the listener
	ct := race.NewCarTracker(ctx)
	go ct.ListenToCars(conn.WriteMessage)

	ws.carsMessageHandler(cancel, conn, ct)
}

// carsMessageHandler Listen for incoming messages
func (ws *WS) carsMessageHandler(cancel context.CancelFunc, conn *websocket.Conn, ct *race.CarTracker) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			cancel()
			break
		}
		fmt.Printf("Received: %s\n", message)
		go ct.StartCarRace(message)
	}
}
