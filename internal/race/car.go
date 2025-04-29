package race

import (
	"context"
	"fmt"
	"time"
)

type Car struct {
	ctx   context.Context
	ch    chan<- *CarMessage
	id    string
	speed int
	road  int
	path  int
	luck  int
}

func NewCar(ctx context.Context, ch chan<- *CarMessage, cfr *CarForRace) *Car {
	return &Car{
		ctx:   ctx,
		ch:    ch,
		id:    cfr.ID,
		speed: cfr.Speed,
		road:  cfr.Road,
		luck:  cfr.Luck,
	}
}

func (c *Car) StartRace() {
	ticker := time.NewTicker(time.Duration(1000/c.speed) * time.Millisecond)
	defer ticker.Stop()

	// ticker with calculated speed
	for {
		select {
		// release goroutine
		case <-c.ctx.Done():
			fmt.Println("Car was removed from the race: id -", c.id)
			return
		case <-ticker.C:
			// the end of the road
			if c.path >= c.road {
				c.ch <- &CarMessage{
					CarId:  c.id,
					Status: Finished,
				}
				return
			}
			// the random obstacle
			if rng.Intn(11) >= c.luck {
				c.ch <- &CarMessage{
					CarId:  c.id,
					Status: Obstacle,
				}
				continue
			}
			// the movement
			c.ch <- &CarMessage{
				CarId:  c.id,
				Status: Moving,
			}
			c.path++
		}
	}
}
