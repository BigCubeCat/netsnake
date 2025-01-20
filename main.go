package main

import (
	"fmt"
	"os"
	"snake/internal/config"
	"snake/internal/model"
	"snake/internal/ui"
	"snake/internal/utils"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sirupsen/logrus"
)

func setup_logger(argparseConfig config.CliConfig) {
	if argparseConfig.LogLevel == "INFO" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if argparseConfig.LogLevel == "WARN" {
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func log_data(g *model.Game) {
	logrus.Debug(fmt.Sprintf("State ID:", g.State.StateID))
	for _, snake := range g.State.Snakes {
		logrus.Debug(fmt.Sprintf(
			"Snake %d: %v, Score: %d, Alive: %v\n",
			snake.ID,
			snake.Body,
			snake.Score,
			snake.IsAlive,
		))
	}
	logrus.Debug("Food:", g.State.Food)
	logrus.Debug("-------------------")
}

func main() {
	argparseConfig := config.ArgParse()
	setup_logger(argparseConfig)
	envConfig := config.EnvParse()
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()
	logrus.SetOutput(f)

	game := model.NewGame(
		int(envConfig.FieldWidth),
		int(envConfig.FieldHeight),
		int(envConfig.FoodStatic),
	)

	// Добавляем змейку
	userId := utils.RandRange(2, 16581375)
	otherId := utils.RandRange(2, 16581375)
	fmt.Println("userID = ", userId)
	game.AddSnake(userId, model.Master)
	game.AddSnake(otherId, model.Normal)

	game.MoveSnakes()

	ticker := time.NewTicker(time.Millisecond * time.Duration(envConfig.Dt))
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Основной игровой цикл
				game.MoveSnakes()
				game.MoveSnake(otherId, utils.RandRange(0, 4))
				log_data(game)
			}
		}
	}()
	m := ui.NewUi(game, userId)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}
}
