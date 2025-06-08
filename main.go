// lexis app
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles for each state of a letter in the grid
var defaultStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	Border(lipgloss.RoundedBorder())

var activeStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#8af")).
	Inherit(defaultStyle)

var exactMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#7d7")).
	Inherit(defaultStyle)

var existsMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#cc0")).
	Inherit(defaultStyle)

var notMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#444")).
	Inherit(defaultStyle)

// letter represents a single letter and its style
type letter struct {
	r     rune
	style lipgloss.Style
}

// word represents a row of letters
type word []letter

// grid represents the grid where each row is a word
type grid []word

type model struct {
	grid   grid
	ri     int // row index
	ci     int // column index
	answer []rune
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
					m.grid[m.ri][m.ci].r = msg.Runes[0]
					// move to the next column only if we're one position before the end
					if m.ci < len(m.grid[m.ri])-1 {
						m.ci++
					}
				}
			}
		case "enter":
			// if the key is enter, move to the next row if we're one row before the end and the row is full
			if m.ri < len(m.grid)-1 && (m.ci == len(m.grid[m.ri])-1 && m.grid[m.ri][m.ci].r != ' ') {
				for i, l := range m.grid[m.ri] {
					// change style based on match
					if l.r == m.answer[i] {
						m.grid[m.ri][i].style = exactMatchStyle
					} else {
						// TODO: correctly handle letters that exist twice or more times in the answer
						for j, al := range m.answer {
							if j != i && l.r == al {
								m.grid[m.ri][i].style = existsMatchStyle
								break
							}
							m.grid[m.ri][i].style = notMatchStyle
						}

					}
				}
				m.ri++
				m.ci = 0 // reset column index
			}
		case "backspace":
			// if the key is backspace, remove the last character in the current row
			if m.ci > 0 {
				if m.grid[m.ri][m.ci].r == ' ' {
					m.ci--
				}
			}
			m.grid[m.ri][m.ci].r = ' '
		}

	}
	// reset the style of all letters in the row
	for i := range m.grid[m.ri] {
		m.grid[m.ri][i].style = defaultStyle
	}
	// set the style of the current letter to active
	m.grid[m.ri][m.ci].style = activeStyle
	return m, nil
}

func (m model) View() string {
	// create a view of the grid
	rows := make([]string, len(m.grid))
	for _, w := range m.grid {
		row := []string{}
		for _, l := range w {
			row = append(row, l.style.Render(string(l.r)))
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
		grid[i] = make([]letter, 5)
		for j := range grid[i] {
			grid[i][j] = letter{r: ' ', style: defaultStyle}
		}
	}
	p := tea.NewProgram(model{grid: grid, answer: []rune("lexis")}, tea.WithAltScreen())

	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
