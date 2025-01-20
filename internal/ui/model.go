package ui

import (
	"github.com/bigcubecat/netsnake/internal/model"

	"time"
)

const (
	tickRate = 120 * time.Millisecond
)

type uiModel struct {
	Game   *model.Game
	UserId int
	width  int
	height int
}

type userUiAction struct{}
