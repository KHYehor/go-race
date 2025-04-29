package race

type CarStatus int

const (
	Moving CarStatus = iota
	Obstacle
	Finished
	Error
)

type CarForRace struct {
	ID    string `json:"id"`
	Speed int    `json:"speed"`
	Road  int    `json:"road"`
	Luck  int    `json:"luck"`
}

type CarMessage struct {
	CarId       string    `json:"carId"`
	Status      CarStatus `json:"status"`
	Description string    `json:"description"`
}
