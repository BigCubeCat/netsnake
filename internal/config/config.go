package config

// объединенный конфик
type Config struct {
	// настройки из args
	CliConfig CliConfig
	// настройки из переменных среды
	EnvConfig EnvConfig
	// настройик из UI
	UiConfig UiConfig
}
