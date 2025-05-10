package model_test

import (
	// "fmt"
	"github.com/bigcubecat/netsnake/internal/model"
	"testing"
)

// Testmodel.NewGame checks if the game is initialized correctly.
func TestNewGame(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	if g.State.Width != 10 {
		t.Errorf("Expected width 10, got %d", g.State.Width)
	}
	if g.State.Height != 10 {
		t.Errorf("Expected height 10, got %d", g.State.Height)
	}
	if len(g.State.Snakes) != 0 {
		t.Errorf("Expected no snakes, got %d", len(g.State.Snakes))
	}
	if g.State.StateID != 0 {
		t.Errorf("Expected StateID 0, got %d", g.State.StateID)
	}
	if g.State.GameOver {
		t.Errorf("Expected game not over, got GameOver=true")
	}
}

// TestAddSnake checks if a snake is added correctly.
func TestAddSnake(t *testing.T) {
	g := model.NewGame(20, 20, 1)
	success := g.AddSnake(1, model.Master)
	if !success {
		t.Error("Failed to add snake when there should be space")
	}
	if len(g.State.Snakes) != 1 {
		t.Errorf("Expected 1 snake, got %d", len(g.State.Snakes))
	}
	snake := g.State.Snakes[0]
	if len(snake.Body) != 2 {
		t.Errorf("Expected snake body length 2, got %d", len(snake.Body))
	}
	if snake.Score != 0 {
		t.Errorf("Expected score 0, got %d", snake.Score)
	}
	if !snake.IsAlive {
		t.Error("Expected snake to be alive")
	}
}

// TestMoveSnakes checks snake movement and eating food.
func TestMoveSnakes(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	snake := g.State.Snakes[0]
	originalHead := snake.Body[0]
	g.MoveSnakes()
	newHead := snake.Body[0]
	if newHead == originalHead {
		t.Error("Snake did not move")
	}
	if len(snake.Body) != 2 {
		t.Errorf("Expected body length 2, got %d", len(snake.Body))
	}
	// Simulate eating food
	g.State.Food = append(g.State.Food, model.Point{newHead.X + 1, newHead.Y})
	g.MoveSnakes()
	if len(snake.Body) != 3 {
		t.Errorf("Expected body length 3 after eating food, got %d", len(snake.Body))
	}
}

// TestChangeDirection checks if the snake's direction changes correctly.
func TestChangeDirection(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	snake := g.State.Snakes[0]
	originalDir := snake.Direction
	g.ChangeDirection(1, model.Up)
	if snake.Direction == originalDir {
		t.Error("Direction did not change when it should")
	}
	// Try to change to opposite direction
	g.ChangeDirection(1, model.Down)
	if snake.Direction == model.Down {
		t.Error("Direction changed to opposite, which should not happen")
	}
}

// TestIsGameOver checks if the game over condition is detected correctly.
func TestIsGameOver(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	// No snakes, should be game over
	if !g.IsGameOver() {
		t.Error("Expected game to be over with no snakes")
	}
	g.AddSnake(1, model.Master)
	if g.IsGameOver() {
		t.Error("Expected game not to be over with alive snakes")
	}
	snake := g.State.Snakes[0]
	snake.IsAlive = false
	if !g.IsGameOver() {
		t.Error("Expected game to be over with no alive snakes")
	}
}

// TestToroidalField checks the toroidal field behavior.
func TestToroidalField(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	snake := g.State.Snakes[0]
	snake.Direction = model.Right
	// Move snake right past the edge
	snake.Body[0].X = g.State.Width - 1
	g.MoveSnakes()
	newHead := snake.Body[0]
	if newHead.X != 0 {
		t.Errorf("Expected head to wrap to X=0, got X=%d", newHead.X)
	}
}

// TestFoodGeneration ensures food is generated correctly.
func TestFoodGeneration(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.GenerateFood()
	expectedFood := 1 + len(g.State.Snakes)
	if len(g.State.Food) != expectedFood {
		t.Errorf("Expected %d food items, got %d", expectedFood, len(g.State.Food))
	}
	// Add a snake and check food regeneration
	g.AddSnake(1, model.Master)
	g.GenerateFood()
	expectedFood = 1 + len(g.State.Snakes)
	if len(g.State.Food) != expectedFood {
		t.Errorf(
			"Expected %d food items after adding snake, got %d",
			expectedFood,
			len(g.State.Food),
		)
	}
}

// TestHandleCollision checks if collision handling converts body parts to food.
func TestHandleCollision(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	snake := g.State.Snakes[0]
	g.HandleCollision(snake)
	// Check if some food is added
	if len(g.State.Food) < 1 {
		t.Error("Expected some food to be added after collision")
	}
}

func TestField(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(2, model.Master)
	g.MoveSnakes()
	snake := g.State.Snakes[0]
	field := g.Field()

	foodCount := 0
	for _, row := range field {
		for _, cell := range row {
			if cell == 1 {
				foodCount++
			}
		}
	}
	expectedFood := g.FoodStatic + len(g.State.Snakes)
	if foodCount != expectedFood {
		t.Errorf("Expected %d food items, got %d", expectedFood, foodCount)
	}

	for _, part := range snake.Body {
		x := part.X % g.State.Width
		y := part.Y % g.State.Height
		if field[y][x] != snake.ID {
			t.Errorf("Expected cell (%d,%d) to be %d, got %d", x, y, snake.ID, field[y][x])
		}
	}
}

func TestFieldToroidal(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	snake := g.State.Snakes[0]
	snake.Body[0] = model.Point{-1, -1}
	field := g.Field()
	x := (-1 + g.State.Width) % g.State.Width
	y := (-1 + g.State.Height) % g.State.Height
	if field[y][x] != snake.ID {
		t.Errorf("Expected cell (%d,%d) to be %d, got %d",
			x,
			y,
			snake.ID, field[y][x],
		)
	}
}

func TestFieldCollision(t *testing.T) {
	g := model.NewGame(10, 10, 1)
	g.AddSnake(1, model.Master)
	g.AddSnake(2, model.Master)
	g.State.Snakes[1].Body[0] = g.State.Snakes[0].Body[0]
	field := g.Field()
	cellValue := field[g.State.Snakes[0].Body[0].Y][g.State.Snakes[0].Body[0].X]
	if cellValue != g.State.Snakes[1].ID {
		t.Errorf("Expected cell to be %d, got %d", g.State.Snakes[1].ID, cellValue)
	}
}
