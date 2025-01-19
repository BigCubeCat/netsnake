package ui

import (
	"snake/internal/game"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

const (
	tickRate = 120 * time.Millisecond
)

type model struct {
	UserId game.ItemType
	gc     *game.GameController
	width  int
	height int
}

func NewUi(gc *game.GameController, userId game.ItemType) model {
	logrus.Debug("geometry = ", gc.Field.Geometry.Row, gc.Field.Geometry.Col)
	m := model{
		UserId: userId,
		gc:     gc,
		width:  gc.Field.Geometry.Row,
		height: gc.Field.Geometry.Col,
	}
	return m
}

type userUiAction struct{}

func (m model) Init() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return userUiAction{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		logrus.Debug("window resize message")
	case tea.KeyMsg:
		keyCode := msg.String()
		logrus.Error("KEYCODE: ", keyCode)
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k", "w":
			m.gc.UserChangeDirection(m.UserId, game.DIR_UP)
			return m, nil
		case "down", "j", "s":
			m.gc.UserChangeDirection(m.UserId, game.DIR_DOWN)
			return m, nil
		case "left", "h", "a":
			m.gc.UserChangeDirection(m.UserId, game.DIR_LEFT)
			return m, nil
		case "right", "l", "d":
			m.gc.UserChangeDirection(m.UserId, game.DIR_RIGHT)
			return m, nil
		}
	case userUiAction:
		return m, tea.Tick(tickRate, func(t time.Time) tea.Msg {
			return userUiAction{}
		})
	}
	return m, cmd
}

func (m model) View() string {
	snakeBoard := []string{}
	for i := 0; i < m.height; i++ {
		row := []string{}
		for j := 0; j < m.width; j++ {
			cellValue := m.gc.Field.At(i, j)
			if cellValue == game.FOOD_FIELD {
				row = append(row, lipgloss.NewStyle().Render("🍎"))
			} else if cellValue == game.EMPTY_FIELD {
				row = append(row, "  ")
			} else if cellValue == m.UserId {
				row = append(row, lipgloss.NewStyle().Background(lipgloss.Color("#AFFFAF")).Render("◖◗"))
			} else {
				row = append(row, lipgloss.NewStyle().Background(lipgloss.Color("#FFAFAF")).Render("◖◗"))
			}
		}
		snakeBoard = append(
			snakeBoard,
			lipgloss.
				JoinHorizontal(
					lipgloss.Top,
					row...,
				),
		)
	}

	return lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				snakeBoard...,
			),
		)
}
