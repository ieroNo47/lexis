// lexis app
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	Border(lipgloss.RoundedBorder())

type word []rune
type grid []word

type model struct {
	grid grid
	ri   int // row index
	ci   int // column index
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
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			// if the key is a letter, add it to the grid at the current position
			if m.ri < len(m.grid) && m.ci < len(m.grid[m.ri]) {
				if len(msg.Runes) > 0 {
					m.grid[m.ri][m.ci] = msg.Runes[0]
					// move to the next column only if we're one position before the end
					if m.ci < len(m.grid[m.ri])-1 {
						m.ci++
					}
				}
			}
		case "enter":
			// if the key is enter, move to the next row if we're one row before the end
			if m.ri < len(m.grid)-1 {
				m.ri++
				m.ci = 0 // reset column index
			}
		case "backspace":
			// if the key is backspace, remove the last character in the current row
			if m.ci > 0 {
				if m.grid[m.ri][m.ci] == ' ' {
					m.ci--
				}
			}
			m.grid[m.ri][m.ci] = ' '
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
			row = append(row, defaultStyle.Render(string(r)))
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
	// grid[0] = []rune{'l', 'e', 'x', 'i', 's'}
	p := tea.NewProgram(model{grid: grid}, tea.WithAltScreen())

	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
