package ui

import (
	"github.com/bigcubecat/netsnake/internal/config"
	"github.com/bigcubecat/netsnake/internal/network"

	"time"
)

const (
	tickRate = 120 * time.Millisecond
)

type ScreenId int

const (
	GreatingScreen    ScreenId = iota
	JoinScreen                 // сцена со списком доступных игр
	EnterConfigScreen          // сцена где мастер настраивает имя игры и тп
	GameScreen                 // сцена с игрой
)

type uiModel struct {
	GamePeer *network.Peer
	Config   *config.Config
	UserId   int
	width    int
	height   int
}

type userUiAction struct{}
