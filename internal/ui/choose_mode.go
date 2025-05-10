package ui

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type chooseModel struct {
	list  list.Model
	Value *int
}

func (m chooseModel) Init() tea.Cmd {
	return nil
}

func (m chooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "enter" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	*m.Value = m.list.Index()
	return m, cmd
}

func (m chooseModel) View() string {
	return docStyle.Render(m.list.View())
}

func RunChooseMode(value *int) {
	items := []list.Item{
		item{title: "Host game", desc: "You will be distributing the game"},
		item{title: "Join game", desc: "Find other players"},
		item{title: "View game", desc: "Become a fan"},
	}

	m := chooseModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Choose game mode"
	m.Value = value

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		logrus.Println("Error running program:", err)
		os.Exit(1)
	}
}
