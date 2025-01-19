package main

import (
	"fmt"
	"os"
	"snake/internal/config"
	"snake/internal/model"
	"snake/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	argparseConfig := config.ArgParse()
	fmt.Println(argparseConfig)
	envConfig := config.EnvParse()
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	defer f.Close()
	logrus.SetOutput(f)

	g := model.NewGame(
		int(envConfig.FieldWidth),
		int(envConfig.FieldHeight),
		int(envConfig.FoodStatic),
	)

	// Добавляем змейку
	g.AddSnake(1)
	g.AddSnake(2)

	dirs := make([]model.Direction, 4)
	dirs[0] = model.Up
	dirs[1] = model.Down
	dirs[2] = model.Left
	dirs[3] = model.Right

	// Основной игровой цикл
	for !g.IsGameOver() {
		g.MoveSnakes()
		fmt.Println("State ID:", g.State.StateID)
		for _, snake := range g.State.Snakes {
			fmt.Printf(
				"Snake %d: %v, Score: %d, Alive: %v\n",
				snake.ID,
				snake.Body,
				snake.Score,
				snake.IsAlive,
			)
		}
		fmt.Println("Food:", g.State.Food)
		fmt.Println("-------------------")
		time.Sleep(500 * time.Millisecond)
		g.ChangeDirection(1, dirs[utils.RandRange(0, 3)])
	}

	fmt.Println("Game Over!")
	/*
	   	if argparseConfig.LogLevel == "INFO" {
	   		logrus.SetLevel(logrus.InfoLevel)
	   	} else if argparseConfig.LogLevel == "WARN" {

	   		logrus.SetLevel(logrus.WarnLevel)
	   	} else {

	   		logrus.SetLevel(logrus.DebugLevel)
	   	}

	   myGame := game.NewGameController(envConfig)
	   userId, err := myGame.AddUser()

	   	if err != nil {
	   		logrus.Error(err.Error())
	   		return
	   	}

	   ticker := time.NewTicker(time.Millisecond * time.Duration(envConfig.Dt))
	   done := make(chan bool)

	   	go func() {
	   		for {
	   			select {
	   			case <-done:
	   				return
	   			case t := <-ticker.C:
	   				logrus.Println("Tick at", t)
	   				myGame.Step()
	   			}
	   		}
	   	}()

	   m := ui.NewUi(myGame, userId)

	   	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
	   		fmt.Println("Uh oh, we encountered an error:", err)
	   		os.Exit(1)
	   	}
	*/
}
