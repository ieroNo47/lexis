// lexis app
package main

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	grid     grid
	answer   []rune
	finished bool // whether the game is finished
	win      bool // whether the game is won
}

// write empty versions of init update and view functions required for our bubbletea model
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// oVertical := containerStyle.GetBorderTopSize() +
		// 	containerStyle.GetBorderBottomSize() +
		// 	containerStyle.GetMarginTop() +
		// 	containerStyle.GetMarginBottom()

		oHorizontal := containerStyle.GetBorderLeftSize() +
			containerStyle.GetBorderRightSize() +
			containerStyle.GetMarginLeft() +
			containerStyle.GetMarginRight()

		// size of the parent container adjusted to be the window size - the size of the borders and margins
		containerStyle = containerStyle.Width(msg.Width - oHorizontal)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		// If the game is finished, ignore all other key presses
		if m.finished {
			if msg.String() == "r" {
				m.grid.reset()
				m.finished = false
				m.win = false
			}
			return m, nil
		}

		switch msg.String() {
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			// if the key is a letter, add it to the grid
			m.grid.setLetter(msg.Runes[0])
		case "enter":
			// if row is full, evaluate the row
			// row is full if the last letter is not a space
			if m.grid.rowFull() {
				if isMatch(m.answer, m.grid.words[m.grid.rowIndex]) {
					m.win = true // mark the game as won
				}
				// temp slice to keep track of letters that are still to be matched
				tw := make(tempWord, len(m.answer))
				copy(tw, m.answer)
				// first pass: check for exact matches
				for i, l := range m.grid.words[m.grid.rowIndex] {
					// change style based on match
					if l.r == m.answer[i] {
						m.grid.updateState(m.grid.rowIndex, i, matched) // mark the letter as matched
						tw = tw.remove(l.r)                             // remove the letter from the temporary word
					}
				}
				// second pass: check for exists matches and not matches
				// having a separate pass for exists matches allows us to not mark a letter as exists if it was already matched
				for i, l := range m.grid.words[m.grid.rowIndex] {
					if tw.has(l.r) && l.state != matched {
						m.grid.words[m.grid.rowIndex][i].style = existsMatchStyle
						m.grid.words[m.grid.rowIndex][i].state = exists // mark the letter as exists
						tw = tw.remove(l.r)                             // remove the letter from the temporary word
					} else if l.state != matched {
						m.grid.words[m.grid.rowIndex][i].style = notMatchStyle
						m.grid.words[m.grid.rowIndex][i].state = notMatched // mark the letter as not matched
					}
				}
				// move to the next row if we're one row before the end

				if m.win {
					// if the game is won, mark it as finished
					m.finished = true
				} else {
					rowChanged := m.grid.goToNextRow()
					if !rowChanged {
						m.finished = true // mark the game as finished if there are no more rows
					}
				}
			}
		case "backspace":
			// if the key is backspace, delete the last letter
			m.grid.deleteLetter()
		}

	}

	return m, nil
}

func (m model) View() string {
	// status row
	var statusRow string
	if m.finished {
		if m.win {
			statusRow = lipgloss.NewStyle().Foreground(lipgloss.Color("#7d7")).Render("You won!", fmt.Sprintf("%d/%d attempts", m.grid.rowIndex+1, len(m.grid.words)))
		} else {
			statusRow = lipgloss.NewStyle().Foreground(lipgloss.Color("#cc0")).Render("Better luck next time!", fmt.Sprintf("The answer was: %s", string(m.answer)))
		}
		statusRow = lipgloss.NewStyle().Foreground(lipgloss.Color("#8af")).Render(statusRow + "\nGame Over! Press ctrl+c or Esc to exit, r to restart.")
	} else {
		// debug row
		statusRow = lipgloss.NewStyle().Render(fmt.Sprintf("Row: %d, Col: %d, RL: %d, L: %c, A: %s",
			m.grid.rowIndex,
			m.grid.colIndex,
			len(m.grid.words[m.grid.rowIndex])-1,
			m.grid.words[m.grid.rowIndex][m.grid.colIndex].r,
			string(m.answer)))

	}
	// rows = append(rows, statusRow)
	view := lipgloss.JoinVertical(lipgloss.Center, m.grid.render(), statusRow)
	return containerStyle.Render(view)
}

func main() {
	// create a new bubbletea program with our model
	grid := newGrid(6, 5)
	p := tea.NewProgram(model{grid: grid, answer: []rune("whole")}, tea.WithAltScreen())

	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func isMatch(answer []rune, guess word) bool {
	for i, l := range guess {
		if l.r != answer[i] {
			return false
		}
	}
	return true
}

// tempWord is a slice alias for []rune that provides methods to check for existence and remove letters.
// Used to keep track of letters that are still to be matched.
type tempWord []rune

func (tw tempWord) has(r rune) bool {
	return slices.Contains(tw, r)
}

func (tw tempWord) remove(r rune) tempWord {
	for i, tr := range tw {
		if tr == r {
			return slices.Delete(tw, i, i+1)
		}
	}
	return tw
}
