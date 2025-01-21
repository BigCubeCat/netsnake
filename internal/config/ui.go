package config

const MASTER_MODE = 0
const NORMAL_MODE = 1
const VIEWER_MODE = 2

type UiConfig struct {
	Mode        int
	GameAddress string
}
