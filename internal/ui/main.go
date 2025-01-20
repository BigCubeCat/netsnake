package ui

import (
	"snake/internal/model"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func NewUi(game *model.Game, userId int) uiModel {
	m := uiModel{
		Game:   game,
		UserId: userId,
		width:  (*game).Width(),
		height: (*game).Height(),
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
			(*m.Game).MoveSnake(m.UserId, 0)
			return m, nil
		case "down", "j", "s":
			(*m.Game).MoveSnake(m.UserId, 1)
			return m, nil
		case "left", "h", "a":
			(*m.Game).MoveSnake(m.UserId, 2)
			return m, nil
		case "right", "l", "d":
			(*m.Game).MoveSnake(m.UserId, 3)
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
