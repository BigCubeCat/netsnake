package ui

import (
	"fmt"
	"snake/internal/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m uiModel) renderUserId(basic string, id int) string {
	var style = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color(utils.IntToHexColor(int(^uint(id) >> 1)))).
		Background(lipgloss.Color(utils.IntToHexColor(id)))
	return style.Render(basic)
}

func (m uiModel) renderAlive(isAlive bool) string {
	if isAlive {
		return lipgloss.NewStyle().Render("🐍")
	}
	return "💀"
}

func (m uiModel) viewLeaderBoard() string {
	rows := [][]string{}
	board := m.Game.LiderBoard()
	for i, user := range board {
		rows = append(rows, []string{
			m.renderUserId(fmt.Sprintf("[%d]", i+1), user.UserId),
			m.renderUserId(fmt.Sprintf("%d", user.UserId), user.UserId),
			m.renderUserId(fmt.Sprintf("%d", user.Score), user.UserId),
			m.renderAlive(user.Alive),
		})
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	usualStyle := lipgloss.NewStyle()

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("white"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			default:
				return usualStyle
			}
		}).
		Headers("CHART", "ID", "SCORE", "ALIVE").
		Rows(rows...)
	return t.Render()
}
