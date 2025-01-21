package common

type GameInterface interface {
	Field() [][]int
	LiderBoard() []SnakeScore
	Width() int
	Height() int
	MoveSnake(id, dir int)
}

type SnakeScore struct {
	UserId int
	Score  int
	Alive  bool
}

type Player struct {
	ID        int
	Name      string
	IpAddress string
	Port      int
	Score     int
}
