package game

type SnakeDirection uint8

const DIR_UP = SnakeDirection(1)
const DIR_DOWN = SnakeDirection(3)
const DIR_LEFT = SnakeDirection(4)
const DIR_RIGHT = SnakeDirection(12)

/*
Возращает true, если нужно поворачивать змейку по новому направлению
*/
func NeedRotate(prev, next SnakeDirection) bool {
	return prev&next == 0
}

func NewDirection(prev, next SnakeDirection) SnakeDirection {
	if NeedRotate(prev, next) {
		return next
	}
	return prev
}
