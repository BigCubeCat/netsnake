package main

import (
	"os"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/logging"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/bigcubecat/netsnake/internal/network"
	"github.com/bigcubecat/netsnake/internal/ui"
	"github.com/sirupsen/logrus"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var peer *network.Peer
	argparseConfig := config.ArgParse()
	logging.SetupLogger(argparseConfig)
	envConfig := config.EnvParse()

	f, err := os.OpenFile(argparseConfig.GameName+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Printf("error opening file: %v", err)
	}
	defer f.Close()
	logrus.SetOutput(f)

	uiConfig := config.UiConfig{}
	ui.RunChooseMode(&uiConfig.Mode)
	conf := config.Config{
		CliConfig: argparseConfig,
		EnvConfig: envConfig,
		UiConfig:  uiConfig,
	}
	logrus.Println(conf)

	// создаем объект peer, который работает с сетью и игрой
	peer = network.NewPeer(
		model.NewGame(
			int(envConfig.FieldWidth),
			int(envConfig.FieldHeight),
			int(envConfig.FoodStatic),
		),
		&conf,
	)
	peer.StartGorutines()
	logrus.Println("after start go")

	if uiConfig.Mode == 0 {
		// создаем игру, если мы - MASTER
	} else {
		// найти игры
	}
	logrus.Println(peer.GameInstance)

	m := ui.NewUi(peer, peer.ID)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		logrus.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
