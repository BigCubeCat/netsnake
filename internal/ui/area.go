package ui

import (
	"snake/internal/model"
	"snake/internal/utils"

	// "snake/internal/utils"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func NewUi(game *model.Game, userId int) uiModel {
	m := uiModel{
		Game:   game,
		UserId: userId,
		width:  (*game).Width(),
		height: (*game).Height(),
	}
	return m
}

type userUiAction struct{}

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
	field := m.Game.Field()
	snakeBoard := []string{}
	for i := 0; i < len(field); i++ {
		row := []string{}
		for j := 0; j < len(field[i]); j++ {
			cellValue := field[i][j]
			// logrus.Debug(fmt.Sprintf("field[%d][%d]=%d", i, j, cellValue))
			if cellValue == 1 {
				row = append(row, lipgloss.NewStyle().Render("🍎"))
			} else if cellValue == 0 {
				row = append(row, lipgloss.NewStyle().Render("  "))
			} else if cellValue == m.UserId {
				row = append(row,
					lipgloss.NewStyle().Background(
						lipgloss.Color(
							utils.IntToHexColor(cellValue),
						),
					).Render("🎲"),
				)
			} else {
				row = append(row,
					lipgloss.NewStyle().Background(
						lipgloss.Color(
							utils.IntToHexColor(cellValue),
						),
					).Render("  "),
				)
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
