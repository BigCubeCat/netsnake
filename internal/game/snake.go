package game

type Snake struct {
	UserId    ItemType
	Direction SnakeDirection
	IsAlive   bool
	Body      []GameCoords
	Size      uint
}

func (s *Snake) Rotate(dir SnakeDirection) {
	if NeedRotate(s.Direction, dir) {
		s.Direction = dir
	}
}

type MoveResult uint8

// Суицид
const MOVE_RESULT_SUICIDE = MoveResult(0)

// Убиты
const MOVE_RESULT_KILLED = MoveResult(10)

// Передвинулись
const MOVE_RESULT_IDLE = MoveResult(1)

// Увеличелись
const MOVE_RESULT_INCR = MoveResult(2)
