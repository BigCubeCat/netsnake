package game

import "fmt"

type GameCoords struct {
	Row int
	Col int
}

func (gc GameCoords) Stringify() string {
	return fmt.Sprintf("%d;%d", gc.Row, gc.Col)
}
