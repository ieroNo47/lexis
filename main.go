// lexis app
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	Border(lipgloss.RoundedBorder())

type word []rune
type grid []word

type model struct {
	grid grid
}

// write empty versions of init update and view functions required for our bubbletea model
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	// create a view of the grid
	rows := make([]string, len(m.grid))
	for _, w := range m.grid {
		row := []string{}
		for _, r := range w {
			row = append(row, style.Render(string(r)))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Bottom, row...))
	}
	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return view
}

func main() {
	// create a new bubbletea program with our model
	grid := make([]word, 5)
	for i := range grid {
		grid[i] = []rune{' ', ' ', ' ', ' ', ' '}
	}
	grid[0] = []rune{'l', 'e', 'x', 'i', 's'}
	p := tea.NewProgram(model{grid: grid}, tea.WithAltScreen())

	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
