package logging

import (
	"fmt"

	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/model"
	"github.com/sirupsen/logrus"
)

func SetupLogger(argparseConfig config.CliConfig) {
	if argparseConfig.LogLevel == "INFO" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if argparseConfig.LogLevel == "WARN" {
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func LogData(g *model.Game) {
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
