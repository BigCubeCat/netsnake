package config

import (
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"os"
)

type CliConfig struct {
	PlayerName       string
	GameName         string
	MulticastAddress string
	UnicastAddress   string

	LogLevel string
}

func ArgParse() CliConfig {
	var (
		config        CliConfig
		levelSelector *string
		multicast     *string
		name          *string
		gameName      *string
	)
	parser := argparse.NewParser("Snake", "")
	levelSelector = parser.Selector(
		"l",
		"log-level",
		[]string{"INFO", "DEBUG", "WARN"},
		&argparse.Options{Required: false, Help: "logging level", Default: "INFO"},
	)
	multicast = parser.String(
		"m",
		"multicast",
		&argparse.Options{Required: true, Help: "mutlicast address"},
	)
	name = parser.String(
		"n",
		"name",
		&argparse.Options{Required: true, Help: "player name"},
	)
	gameName = parser.String(
		"g",
		"game",
		&argparse.Options{Required: false, Help: "name for your game, if you host"},
	)

	err := parser.Parse(os.Args)
	if err != nil {
		logrus.Fatal("cannot parse arguments: ", err.Error())
		return config
	}

	config.MulticastAddress = *multicast
	config.LogLevel = *levelSelector
	config.PlayerName = *name
	config.GameName = *gameName
	if config.GameName == "" {
		config.GameName = config.PlayerName + "_game"
	}

	return config
}
