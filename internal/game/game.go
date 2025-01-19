package game

import (
	"snake/internal/config"

	"github.com/sirupsen/logrus"
)

type GameController struct {
	Snakes []Snake
	Field  *GameField

	FoodStat uint
	FoodProb uint

	// Единица измерения времени, ms
	Dt        uint
	ChunkSize uint

	UserActions  map[ItemType]SnakeDirection
	UserFreeType ItemType
}

func NewGameController(config config.EnvConfig) *GameController {
	controller := new(GameController)
	controller.Field = NewField(config.FieldWidth, config.FieldHeight, config.ChunkSize)
	controller.FoodStat = config.FoodStatic
	controller.FoodProb = config.FoodProb
	controller.Dt = config.Dt
	controller.ChunkSize = config.ChunkSize

	controller.UserActions = make(map[ItemType]SnakeDirection)
	controller.UserFreeType = FOOD_FIELD

	return controller
}

func (controller *GameController) Step() {
	diedSnakes := make(map[int]bool, 0)
	if len(controller.Snakes) == 0 {
		logrus.Fatal("no snakes")
	}
	countFood := controller.Field.countFood()
	if countFood < controller.FoodStat+uint(len(controller.Snakes)) {
		controller.Field.SpawnFood()
	}
	for snakeIndex, snake := range controller.Snakes {
		if !snake.IsAlive {
			logrus.Info("died snake: ", snake.UserId)
			controller.Field.ReplaceWithProb(snake.UserId, controller.FoodProb)
			diedSnakes[snakeIndex] = true
		}
		currentMove := snake.Direction
		if controller.UserActions[snake.UserId] > 0 {
			currentMove = controller.UserActions[snake.UserId]
		}
		controller.Snakes[snakeIndex].Direction = NewDirection(snake.Direction, currentMove)
		newSnake, killerId := controller.Field.MoveSnake(controller.Snakes[snakeIndex])
		if killerId > 0 {
			diedSnakes[snakeIndex] = true
			continue
		}
		controller.Snakes[snakeIndex] = newSnake
	}
	var snakes []Snake
	for i, sn := range controller.Snakes {
		if !diedSnakes[i] {
			snakes = append(snakes, sn)
		}
	}
	controller.Snakes = snakes
	controller.UserActions = make(map[ItemType]SnakeDirection)
}

func (controller *GameController) UserChangeDirection(userId ItemType, action SnakeDirection) {
	logrus.Error("user action: ", userId, " action: ", action)
	controller.UserActions[userId] = action
}

func (controller *GameController) AddUser() (ItemType, error) {
	userId := controller.UserFreeType + 1
	snake, err := controller.Field.CreateSnake(userId, controller.ChunkSize)
	controller.Snakes = append(controller.Snakes, *snake)
	if err != nil {
		return userId, err
	}
	controller.UserActions[userId] = DIR_UP

	return userId, nil
}
