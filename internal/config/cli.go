package config

import (
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"os"
)

type CliConfig struct {
	MulticastAddress string
	UnicastAddress   string

	LogLevel string
	IsHost   bool
}

func ArgParse() CliConfig {
	// TODO: доделать
	var (
		config        CliConfig
		levelSelector *string
		multicast     *string
		isServer      *bool
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
	isServer = parser.Flag("f", "force", &argparse.Options{
		Required: false, Help: "if used, host game",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		logrus.Fatal("cannot parse arguments: ", err.Error())
		return config
	}

	config.MulticastAddress = *multicast
	config.LogLevel = *levelSelector
	config.IsHost = *isServer

	return config
}
