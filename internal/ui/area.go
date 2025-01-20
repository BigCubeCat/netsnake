package ui

import (
	"github.com/bigcubecat/netsnake/internal/utils"

	"github.com/charmbracelet/lipgloss"
)

func (m uiModel) viewGameField() string {
	field := m.Game.Field()
	snakeBoard := []string{}
	for i := 0; i < len(field); i++ {
		row := []string{}
		for j := 0; j < len(field[i]); j++ {
			cellValue := field[i][j]
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
					).Render("🎱"),
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
