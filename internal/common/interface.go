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
