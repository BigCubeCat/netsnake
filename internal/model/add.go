package model

// AddSnake добавляет новую змейку на поле.
func (g *Game) AddSnake(id int) bool {
	// Ищем квадрат 5x5 для размещения змейки
	// yStart := utils.RandRange(0, g.State.Height-5)
	// xStart := utils.RandRange(0, g.State.Width-5)
	for y := 0; y < g.State.Height; y++ {
		for x := 0; x < g.State.Width; x++ {
			if g.canPlaceSnake(x, y) {
				head := Point{X: x + 2, Y: y + 2}
				tail := g.randomTailPosition(head)
				snake := &Snake{
					ID:        id,
					Body:      []Point{head, tail},
					Direction: g.oppositeDirection(head, tail),
					Score:     0,
					IsAlive:   true,
				}
				g.State.Snakes = append(g.State.Snakes, snake)
				return true
			}
		}
	}
	return false
}

// canPlaceSnake проверяет, можно ли разместить змейку в квадрате 5x5.
func (g *Game) canPlaceSnake(x, y int) bool {
	for dy := 0; dy < 5; dy++ {
		for dx := 0; dx < 5; dx++ {
			p := Point{X: (x + dx) % g.State.Width, Y: (y + dy) % g.State.Height}
			if g.isOccupied(p) {
				return false
			}
		}
	}
	return true
}

// randomTailPosition возвращает случайное положение хвоста змейки.
func (g *Game) randomTailPosition(head Point) Point {
	directions := []Point{
		{X: head.X - 1, Y: head.Y},
		{X: head.X + 1, Y: head.Y},
		{X: head.X, Y: head.Y - 1},
		{X: head.X, Y: head.Y + 1},
	}

	for _, dir := range directions {
		dir.X = (dir.X + g.State.Width) % g.State.Width
		dir.Y = (dir.Y + g.State.Height) % g.State.Height
		if !g.isOccupied(dir) {
			return dir
		}
	}
	return head
}

// oppositeDirection возвращает направление, противоположное направлению от tail к head.
func (g *Game) oppositeDirection(head, tail Point) Direction {
	if head.X == tail.X {
		if head.Y > tail.Y {
			return Down
		}
		return Up
	}
	if head.X > tail.X {
		return Right
	}
	return Left
}
