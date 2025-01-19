package model

// MoveSnakes обновляет положение всех змеек на поле.
func (g *Game) MoveSnakes() {
	for _, snake := range g.State.Snakes {
		if !snake.IsAlive {
			continue
		}

		// Вычисляем новую позицию головы
		newHead := g.calculateNewHead(snake)

		// Проверяем столкновения
		if g.isCollision(newHead) {
			snake.IsAlive = false
			g.HandleCollision(snake)
			continue
		}

		// Проверяем, съела ли змейка еду
		if g.isFood(newHead) {
			snake.Score++
			g.removeFood(newHead)
			snake.Body = append([]Point{newHead}, snake.Body...)
		} else {
			snake.Body = append([]Point{newHead}, snake.Body[:len(snake.Body)-1]...)
		}
	}

	g.GenerateFood()
	g.State.StateID++
}

// calculateNewHead вычисляет новую позицию головы змейки.
func (g *Game) calculateNewHead(snake *Snake) Point {
	head := snake.Body[0]
	switch snake.Direction {
	case Up:
		return Point{X: head.X, Y: (head.Y - 1 + g.State.Height) % g.State.Height}
	case Down:
		return Point{X: head.X, Y: (head.Y + 1) % g.State.Height}
	case Left:
		return Point{X: (head.X - 1 + g.State.Width) % g.State.Width, Y: head.Y}
	case Right:
		return Point{X: (head.X + 1) % g.State.Width, Y: head.Y}
	}
	return head
}

// isCollision проверяет, столкнулась ли змейка с чем-либо.
func (g *Game) isCollision(p Point) bool {
	for _, snake := range g.State.Snakes {
		for _, bodyPart := range snake.Body {
			if bodyPart == p {
				return true
			}
		}
	}
	return false
}

// isFood проверяет, находится ли еда в указанной точке.
func (g *Game) isFood(p Point) bool {
	for _, food := range g.State.Food {
		if food == p {
			return true
		}
	}
	return false
}

// removeFood удаляет еду из указанной точки.
func (g *Game) removeFood(p Point) {
	for i, food := range g.State.Food {
		if food == p {
			g.State.Food = append(g.State.Food[:i], g.State.Food[i+1:]...)
			break
		}
	}
}

// HandleCollision обрабатывает столкновение змейки.
func (g *Game) HandleCollision(snake *Snake) {
	// Превращаем часть тела змейки в еду с вероятностью 0.5
	for _, bodyPart := range snake.Body {
		if g.Rand.Float64() < 0.5 {
			g.State.Food = append(g.State.Food, bodyPart)
		}
	}
}
