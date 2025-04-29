package race

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type CarTracker struct {
	ctx  context.Context
	ch   chan *CarMessage
	cars []*Car
}

func NewCarTracker(ctx context.Context) *CarTracker {
	return &CarTracker{
		ctx:  ctx,
		cars: make([]*Car, 0),
		ch:   make(chan *CarMessage, 1024),
	}
}

func (ct *CarTracker) validateIncomeMessage(message []byte) (*CarForRace, error) {
	var cfr *CarForRace
	err := json.Unmarshal(message, &cfr)
	return cfr, err
}

func (ct *CarTracker) StartCarRace(message []byte) {
	cfr, err := ct.validateIncomeMessage(message)
	if err != nil {
		log.Printf("error unmarshalling car data: : %v", err)
		ct.ch <- &CarMessage{
			Status:      Error,
			Description: err.Error(),
		}
		return
	}
	car := NewCar(ct.ctx, ct.ch, cfr)
	car.StartRace()
}

func (ct *CarTracker) ListenToCars(cb func(int, []byte) error) {
	for {
		select {
		case <-ct.ctx.Done():
			log.Println("Connection was aborted")
			return
		case msg := <-ct.ch:
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("internal error, can't parse into the []byte: %v", err)
				return
			}
			cb(websocket.TextMessage, data)
		}
	}
}
