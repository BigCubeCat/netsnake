package ui

import (
	"time"

	"github.com/bigcubecat/netsnake/internal/network"
	"github.com/sirupsen/logrus"

	tea "github.com/charmbracelet/bubbletea"
)

func NewUi(game *network.Peer, userId int) uiModel {
	logrus.Println("NewUi")
	m := uiModel{
		GamePeer: game,
		UserId:   userId,
		width:    game.GameInstance.Width(),
		height:   game.GameInstance.Height(),
	}
	return m
}

func (m uiModel) Init() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return userUiAction{}
	})
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k", "w":
			m.GamePeer.MoveSnake(0)
			return m, nil
		case "down", "j", "s":
			m.GamePeer.MoveSnake(1)
			return m, nil
		case "left", "h", "a":
			m.GamePeer.MoveSnake(2)
			return m, nil
		case "right", "l", "d":
			m.GamePeer.MoveSnake(3)
			return m, nil
		}
	case userUiAction:
		return m, tea.Tick(tickRate, func(t time.Time) tea.Msg {
			return userUiAction{}
		})
	}
	return m, cmd
}

func (m uiModel) View() string {
	return m.viewGameField() + "\n" + m.viewLeaderBoard()
}
