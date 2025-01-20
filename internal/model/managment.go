package model

func (g *Game) MoveSnake(id, dir int) {
	g.ChangeDirection(id, Direction(dir))
}

// ChangeDirection изменяет направление змейки.
func (g *Game) ChangeDirection(snakeID int, newDir Direction) {
	for _, snake := range g.State.Snakes {
		if snake.ID == snakeID && snake.IsAlive {
			// Проверяем, что новое направление не противоположно текущему
			if !g.isOppositeDirection(snake.Direction, newDir) {
				snake.Direction = newDir
			}
			break
		}
	}
}

// isOppositeDirection проверяет, является ли newDir противоположным currentDir.
func (g *Game) isOppositeDirection(currentDir, newDir Direction) bool {
	return (currentDir == Up && newDir == Down) ||
		(currentDir == Down && newDir == Up) ||
		(currentDir == Left && newDir == Right) ||
		(currentDir == Right && newDir == Left)
}

// IsGameOver проверяет, завершена ли игра.
func (g *Game) IsGameOver() bool {
	aliveSnakes := 0
	for _, snake := range g.State.Snakes {
		if snake.IsAlive {
			aliveSnakes++
		}
	}
	return aliveSnakes == 0
}
