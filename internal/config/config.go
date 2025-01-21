package config

// объединенный конфик
type Config struct {
	// настройки из args
	cliConfig CliConfig
	// настройки из переменных среды
	envConfig EnvConfig
	// настройик из UI
	uiConfig UiConfig
}
