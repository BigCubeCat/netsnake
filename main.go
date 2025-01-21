package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/logging"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/network"
	"github.com/bigcubecat/netsnake/internal/ui"
	"github.com/bigcubecat/netsnake/internal/utils"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sirupsen/logrus"
)

func main() {
	var game *model.Game
	var peer *network.Peer
	argparseConfig := config.ArgParse()
	logging.SetupLogger(argparseConfig)
	envConfig := config.EnvParse()
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()
	logrus.SetOutput(f)

	uiConfig := config.UiConfig{}
	conf := config.Config{
		CliConfig: argparseConfig,
		EnvConfig: envConfig,
		UiConfig:  uiConfig,
	}
	fmt.Println(conf)

	ui.RunChooseMode(&uiConfig.Mode)

	// создаем объект peer, который работает с сетью и игрой
	peer = network.NewPeer(game, &conf)
	peer.StartGorutines()

	if uiConfig.Mode == 0 {
		// создаем игру, если мы - MASTER
		game = model.NewGame(
			int(envConfig.FieldWidth),
			int(envConfig.FieldHeight),
			int(envConfig.FoodStatic),
		)
	} else {
		// найти игры
	}

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
				game.MoveSnakes()
				game.MoveSnake(otherId, utils.RandRange(0, 4))
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
