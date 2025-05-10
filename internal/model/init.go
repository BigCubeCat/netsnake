package model

import (
	"math/rand"
	"time"
)

// NewGame создает новую игру.
func NewGame(width, height, foodStatic int) *Game {
	randSource := rand.NewSource(time.Now().UnixNano())
	game := &Game{
		State: &GameState{
			Width:    width,
			Height:   height,
			Snakes:   make([]*Snake, 0),
			Food:     make([]Point, 0),
			StateID:  0,
			GameOver: false,
		},
		FoodStatic: foodStatic,
		Rand:       rand.New(randSource),
	}

	game.GenerateFood()
	return game
}

// GenerateFood генерирует еду на поле.
func (g *Game) GenerateFood() {
	foodCount := g.FoodStatic + g.CountAliveSnakes()
	diff := 0
	if len(g.State.Food) < foodCount {
		diff = foodCount - len(g.State.Food)
	}
	for i := 0; i < diff; i++ {
		g.State.Food = append(g.State.Food, g.RandomEmptyPoint())
	}
}

// RandomEmptyPoint возвращает случайную пустую точку на поле.
func (g *Game) RandomEmptyPoint() Point {
	for {
		x := g.Rand.Intn(g.State.Width)
		y := g.Rand.Intn(g.State.Height)
		p := Point{X: x, Y: y}

		if !g.isOccupied(p) {
			return p
		}
	}
}

// isOccupied проверяет, занята ли точка змейкой или едой.
func (g *Game) isOccupied(p Point) bool {
	for _, snake := range g.State.Snakes {
		for _, bodyPart := range snake.Body {
			if bodyPart == p {
				return true
			}
		}
	}

	for _, food := range g.State.Food {
		if food == p {
			return true
		}
	}

	return false
}

func (g *Game) CountAliveSnakes() int {
	count := 0
	for _, snake := range g.State.Snakes {
		if snake.IsAlive {
			count++
		}
	}
	return count
}
