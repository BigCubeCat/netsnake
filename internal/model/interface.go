package model

import "snake/internal/common"

func (g *Game) Field() [][]int {
	field := make([][]int, g.State.Height)
	for i := range field {
		field[i] = make([]int, g.State.Width)
	}
	for _, food := range g.State.Food {
		x := food.X % g.State.Width
		y := food.Y % g.State.Height
		if x < 0 {
			x += g.State.Width
		}
		if y < 0 {
			y += g.State.Height
		}
		field[y][x] = 1
	}
	for _, snake := range g.State.Snakes {
		// if snake.Role == Viewer {
		// 	continue
		// }
		for _, part := range snake.Body {
			x := part.X % g.State.Width
			y := part.Y % g.State.Height
			if x < 0 {
				x += g.State.Width
			}
			if y < 0 {
				y += g.State.Height
			}
			field[y][x] = snake.ID
		}
	}
	return field
}

func (g *Game) LiderBoard() []common.SnakeScore {
	var leaderboard []common.SnakeScore
	for _, snake := range g.State.Snakes {
		if snake.Role == Viewer {
			continue
		}
		leaderboard = append(leaderboard, common.SnakeScore{
			UserId: snake.ID,
			Score:  snake.Score,
			Alive:  snake.IsAlive,
		})
	}
	return leaderboard
}

func (g *Game) Width() int {
	return g.State.Width
}

func (g *Game) Height() int {
	return g.State.Height
}
