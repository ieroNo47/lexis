// lexis app
package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	grid     grid
	keyboard keyboard
	answer   []rune
	finished bool // whether the game is finished
	win      bool // whether the game is won
	log      *log.Logger
	help     help.Model
	keys     keyMap
}

// write empty versions of init update and view functions required for our bubbletea model
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	m.log.Debug("[Update]", "msg", spew.Sdump(msg))
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
		headerStyle = headerStyle.Width(msg.Width - oHorizontal)
		resultBarStyle = resultBarStyle.Width(msg.Width - oHorizontal)
		helpBarStyle = helpBarStyle.Width(msg.Width - oHorizontal)
		m.help.Width = msg.Width - oHorizontal
		m.log.Debug("Window resized", "width", msg.Width, "height", msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.log.Info("==== Bye! ====")
			return m, tea.Quit
		case key.Matches(msg, m.keys.Letter):
			// if the key is a letter, add it to the grid
			// https://pkg.go.dev/github.com/charmbracelet/bubbletea@v1.3.6#KeyMsg
			// Doc: Note that Key.Runes will always contain at least one character, so you can always safely call Key.Runes[0].
			if !m.finished {
				m.grid.setLetter(msg.Runes[0])
			}
		case key.Matches(msg, m.keys.Delete):
			// if the key is backspace, delete the last letter
			if !m.finished {
				m.grid.deleteLetter()
			}
		case key.Matches(msg, m.keys.Restart):
			m.log.Info("==== Restarting game ====")
			m.grid.reset()
			m.keyboard.reset()
			m.finished = false
			m.win = false
			resultBarStyle = resultBarStyle.Background(lipgloss.Color("#8af"))
			return m, nil
		case key.Matches(msg, m.keys.Submit):
			// if row is full, evaluate the row
			if m.grid.rowFull() {
				if isMatch(m.answer, m.grid.words[m.grid.rowIndex]) {
					m.log.Info("Match found")
					m.win = true // mark the game as won
				}
				// temp slice to keep track of letters that are still to be matched
				tw := make(tempWord, len(m.answer))
				copy(tw, m.answer)
				// first pass: check for exact matches
				m.log.Debug("== First Pass: Exact Matches ==")
				for i, l := range m.grid.words[m.grid.rowIndex] {
					// change style based on match
					if l.r == m.answer[i] {
						m.log.Debug("Update", "letter", string(l.r), "index", i, "state", states[matched])
						m.grid.updateState(m.grid.rowIndex, i, matched) // mark the letter as matched
						m.keyboard.updateLetterState(l.r, matched)      // update the keyboard state
						tw = tw.remove(l.r)                             // remove the letter from the temporary word
					}
				}
				// second pass: check for exists matches and not matches
				// having a separate pass for exists matches allows us to not mark a letter as exists if it was already matched
				var targetState = notMatched
				m.log.Debug("== Second Pass: Exists and Not Matches ==")
				for i, l := range m.grid.words[m.grid.rowIndex] {
					if tw.has(l.r) && l.state != matched {
						targetState = exists
						m.log.Debug("Row Update", "letter", string(l.r), "index", i, "state", states[exists])
						m.grid.updateState(m.grid.rowIndex, i, exists) // mark the letter as exists
						tw = tw.remove(l.r)                            // remove the letter from the temporary word
					} else if l.state != matched {
						m.log.Debug("Row Update", "letter", string(l.r), "index", i, "state", "notMatched")
						m.grid.updateState(m.grid.rowIndex, i, notMatched) // mark the letter as not matched
					}
					if m.keyboard.getLetterState(l.r) != matched { // only update the keyboard state if it is not already matched
						m.log.Debug("Keyboard Update", "letter", string(l.r), "to", states[targetState], "from", states[m.keyboard.getLetterState(l.r)])
						m.keyboard.updateLetterState(l.r, targetState)
					}
					targetState = notMatched // reset target state
				}
				// move to the next row if we're one row before the end

				if m.win {
					// if the game is won, mark it as finished
					m.log.Debug("Marking game as finished.", "reason", "win")
					m.finished = true
				} else {
					rowChanged := m.grid.goToNextRow()
					if !rowChanged {
						m.log.Debug("Marking game as finished.", "reason", "no more rows")
						m.finished = true // mark the game as finished if there are no more rows
					}
				}
			}
		}
	}
	// if log level is debug, print the current string in the active row
	if m.log.GetLevel() == log.DebugLevel {
		var currentRow string
		for _, l := range m.grid.words[m.grid.rowIndex] {
			currentRow += string(l.r)
		}
		m.log.Debug("", "row", m.grid.rowIndex, "col", m.grid.colIndex, "string", currentRow)
	}
	return m, nil
}

func (m model) View() string {
	// todo: Move logit to update header/footer to Update()
	// header
	header := headerStyle.Render("lexis")
	var resultS string
	// var statusS string
	if m.finished {
		if m.win {
			resultS = fmt.Sprintf("You won! %d/%d attempts", m.grid.rowIndex+1, len(m.grid.words))
			resultBarStyle = resultBarStyle.Background(lipgloss.Color("#7d7"))
		} else {
			resultS = fmt.Sprintf("Better luck next time! The answer was: %s", string(m.answer))
			resultBarStyle = resultBarStyle.Background(lipgloss.Color("#ffcd44ff"))
		}
		// statusS = "Press ctrl+c or Esc to exit, r to restart."
	} else {
		// debug row
		resultS = fmt.Sprintf("Row: %d, Col: %d, RL: %d, L: %c, A: %s",
			m.grid.rowIndex,
			m.grid.colIndex,
			len(m.grid.words[m.grid.rowIndex])-1,
			m.grid.words[m.grid.rowIndex][m.grid.colIndex].r,
			string(m.answer))
		// statusS = "Press Enter to submit, Backspace to delete."
	}
	resultRow := resultBarStyle.Render(resultS)
	helpRow := helpBarStyle.Render(m.help.View(m.keys))
	// rows = append(rows, helpRow)
	view := lipgloss.JoinVertical(lipgloss.Center, header, m.grid.render(), m.keyboard.render(), resultRow, helpRow)
	return containerStyle.Render(view)
}

func main() {
	// if err := os.Remove("debug.log"); err != nil && !os.IsNotExist(err) {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// create a logger that writes to a file
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	//nolint:errcheck
	defer f.Close()
	logger := log.NewWithOptions(f, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})
	// create a new bubbletea program with our model
	grid := newGrid(6, 5)
	grid.updateStyle(0, 0, activeStyle) // set the first cell as active
	keyboard := newKeyboard()
	p := tea.NewProgram(model{
		grid:     grid,
		keyboard: keyboard,
		answer:   []rune("minty"),
		log:      logger,
		help:     newHelp(),
		keys:     keys,
	}, tea.WithAltScreen())

	logger.Info("==== Starting lexis ====")
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
