package game

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type GameField struct {
	Geometry  GameCoords
	chunkSize uint
	Area      [][]ItemType
}

func NewField(width, height uint, chunk uint) *GameField {
	field := new(GameField)
	field.Geometry = GameCoords{Row: int(height), Col: int(width)}

	field.Area = make([][]ItemType, field.Geometry.Row)
	for i := range field.Geometry.Row {
		field.Area[i] = make([]ItemType, field.Geometry.Col)
		for j := range field.Geometry.Col {
			field.Area[i][j] = EMPTY_FIELD
		}
	}

	field.chunkSize = chunk
	return field
}

/*
Устанавливает значение клетки в поле
*/
func (field *GameField) UpdateItem(row, col int, item ItemType) {
	field.Area[(field.Geometry.Row+row)%field.Geometry.Row][(field.Geometry.Col+col)%field.Geometry.Col] = item
}

/*
Заменяет змейку пользователя на еду с вероятностью
*/
func (field *GameField) ReplaceWithProb(userId ItemType, prob uint) {
	for i := range field.Geometry.Row {
		for j := range field.Geometry.Col {
			if field.Area[i][j] == userId {
				if randomWithPercent(prob) {
					field.Area[i][j] = FOOD_FIELD
				} else {
					field.Area[i][j] = EMPTY_FIELD
				}
			}
		}
	}
}

/*
Возращает значение в поле
*/
func (field GameField) At(row, col int) ItemType {
	return field.Area[row%field.Geometry.Row][col%field.Geometry.Col]
}

/*
Вернет true, если существует пустое значение
*/
func (field GameField) HasEmptyCell() bool {
	for i := range field.Geometry.Row {
		for j := range field.Geometry.Col {
			if field.Area[i][j] == EMPTY_FIELD {
				return true
			}
		}
	}
	return false
}

/*
Возращает случайную пустую клетку
*/
func (field GameField) RandomEmptyCoord() GameCoords {
	row := RandRange(0, int(field.Geometry.Row))
	col := RandRange(0, int(field.Geometry.Row))
	if field.At(row, col) == EMPTY_FIELD {
		return GameCoords{row, col}
	}
	if field.HasEmptyCell() {
		for {
			row = RandRange(0, int(field.Geometry.Row))
			col = RandRange(0, int(field.Geometry.Row))
			if field.At(row, col) == EMPTY_FIELD {
				return GameCoords{row, col}
			}
		}
	}
	logrus.Error("FIELD IS FULL")
	return GameCoords{0, 0}
}

/*
Создает змейку
*/
func (field *GameField) CreateSnake(userId ItemType, chunkSize uint) (*Snake, error) {
	var (
		head GameCoords
		tail GameCoords
	)

	sq := findSquares(field.Area, int(field.chunkSize))
	if len(sq) == 0 {
		return nil, errors.New("room is full")
	}
	index := RandRange(0, len(sq))
	squareCoords := sq[index]

	head.Row = squareCoords.Row + RandRange(0, int(field.chunkSize))
	head.Col = squareCoords.Col + RandRange(0, int(field.chunkSize))
	tail.Row = head.Row
	tail.Col = head.Col

	dirs := [...]SnakeDirection{DIR_DOWN, DIR_LEFT, DIR_RIGHT, DIR_UP}
	direction := dirs[RandRange(0, len(dirs))]
	if direction == DIR_DOWN {
		tail.Row = (field.Geometry.Row + 1) % field.Geometry.Row
	} else if direction == DIR_UP {
		tail.Row = (field.Geometry.Row + field.Geometry.Row - 1) % field.Geometry.Row
	} else if direction == DIR_LEFT {
		tail.Col = (field.Geometry.Col + field.Geometry.Col - 1) % field.Geometry.Col
	} else {
		tail.Col = (field.Geometry.Col + 1) % field.Geometry.Col
	}

	field.UpdateItem(head.Row, head.Col, userId)
	field.UpdateItem(tail.Row, tail.Col, userId)
	snake := Snake{
		Direction: direction,
		UserId:    userId,
		Body:      []GameCoords{head, tail},
		Size:      2,
		IsAlive:   true,
	}
	return &snake, nil
}

func (field GameField) moveCoord(coords GameCoords, dir SnakeDirection) GameCoords {
	if dir == DIR_DOWN {
		return GameCoords{(coords.Row + 1) % field.Geometry.Row, coords.Col}
	}
	if dir == DIR_UP {
		return GameCoords{(field.Geometry.Row + coords.Row - 1) % field.Geometry.Row, coords.Col}
	}
	if dir == DIR_RIGHT {
		return GameCoords{coords.Row, (coords.Col + 1) % field.Geometry.Col}
	}
	if dir == DIR_LEFT {
		return GameCoords{coords.Row, (field.Geometry.Col + coords.Col - 1) % field.Geometry.Col}
	}
	logrus.Error("invalid direction:", dir, " for game coords: ", coords)
	return coords
}

func (field *GameField) MoveSnake(snake Snake) (Snake, ItemType) {
	newHeadCoord := field.moveCoord(snake.Body[0], snake.Direction)
	cellContent := field.At(newHeadCoord.Row, newHeadCoord.Col)
	newSnake := snake
	if cellContent > FOOD_FIELD {
		newSnake.IsAlive = false
		return newSnake, cellContent
	}
	field.UpdateItem(newHeadCoord.Row, newHeadCoord.Col, snake.UserId)
	newSnake.Body = append([]GameCoords{newHeadCoord}, newSnake.Body...)
	if cellContent == FOOD_FIELD {
		newSnake.Size++
	} else {
		lastIndex := len(newSnake.Body) - 1
		node := newSnake.Body[lastIndex]
		field.UpdateItem(node.Row, node.Col, EMPTY_FIELD)
		newSnake.Body = newSnake.Body[:lastIndex]
	}
	return newSnake, 0
}

func (field *GameField) SpawnFood() {
	coord := field.RandomEmptyCoord()
	field.UpdateItem(coord.Row, coord.Col, FOOD_FIELD)
}

func (field GameField) countFood() uint {
	result := uint(0)
	for i := 0; i < field.Geometry.Row; i++ {
		for j := 0; j < field.Geometry.Col; j++ {
			if field.At(i, j) == FOOD_FIELD {
				result++
			}
		}
	}
	return result
}
