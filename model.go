// lexis app
package main

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
)

const (
	stateLoading = iota
	stateRowNotFull
	stateInvalidWord
	statePlaying
)

type model struct {
	game    game
	log     *log.Logger
	help    help.Model
	keys    keyMap
	state   int
	spinner spinner.Model
}

// initCompleteMsg is a message that is sent when the game initialization is complete and the answer is ready
type initCompleteMsg bool

// rowNotFullMsg is a message that is sent when the user tries to submit a guess but the current row is not full
type rowNotFullMsg bool

// invalidWordMsg is a message that is sent when the user tries to submit a guess but the word is not valid
type invalidWordMsg bool

// validWordMsg is a message that is sent when the user submits a valid guess
type validWordMsg bool

// Init initializes the model and starts the game by getting the answer from the answer provider
func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
	}

	answerCmd := func() tea.Msg {
		m.log.Debug("[Init]")
		m.game.initProvider()
		return initCompleteMsg(true)
	}
	cmds = append(cmds, answerCmd)
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	m.log.Debug("[Update]", "msg", spew.Sdump(msg))
	switch msg := msg.(type) {
	// === INIT ===
	case initCompleteMsg:
		m.log.Debug("Initialization complete")
		m.game.start()
		m.state = statePlaying
		return m, nil
	// === WINDOW RESIZE ===
	case tea.WindowSizeMsg:
		oHorizontal := updateStyles(msg)
		m.help.SetWidth(msg.Width - oHorizontal)
		m.log.Debug("Window resized", "width", msg.Width, "height", msg.Height)
	case tea.KeyPressMsg:
		switch {
		// === QUIT ===
		case key.Matches(msg, m.keys.Quit):
			m.log.Info("==== Bye! ====")
			return m, tea.Quit
		// === LETTERS ===
		case key.Matches(msg, m.keys.Letter):
			if m.state != stateLoading {
				m.state = statePlaying
				m.game.processLetter(msg.Text)
			} else {
				m.log.Info("Cannot enter letters while loading")
			}
		// === DELETE ===
		case key.Matches(msg, m.keys.Delete):
			if m.state != stateLoading {
				m.state = statePlaying
				m.game.processDelete()
			} else {
				m.log.Info("Cannot delete while loading")
			}
		// === RESTART ===
		case key.Matches(msg, m.keys.Restart):
			m.log.Info("==== Restarting game ====")
			m.game.reset()
		// === SUBMIT ===
		case key.Matches(msg, m.keys.Submit):
			cmds := []tea.Cmd{m.spinner.Tick}
			m.state = stateLoading
			submitCmd := func() tea.Msg {
				err := m.game.rowReady()
				if err != nil {
					if errors.Is(err, ErrRowNotFull) {
						return rowNotFullMsg(true)
					} else if errors.Is(err, ErrInvalidWord) {
						return invalidWordMsg(true)
					}
				}
				return validWordMsg(true)
			}
			cmds = append(cmds, submitCmd)
			return m, tea.Batch(cmds...)
		}
	// === SPINNER TICK ===
	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	// === SUBMIT RESULTS ===
	case rowNotFullMsg:
		m.state = stateRowNotFull
		return m, nil
	case invalidWordMsg:
		m.state = stateInvalidWord
		return m, nil
	case validWordMsg:
		m.state = statePlaying
		m.game.Submit()
	default:
		// if log level is debug, print the current string in the active row
		if m.log.GetLevel() == log.DebugLevel {
			rowIndex, colIndex, rowString := m.game.debugState()
			m.log.Debug("", "row", rowIndex, "col", colIndex, "string", rowString)
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	// todo: Move logic to update header/footer to Update()
	// header
	header := headerStyle.Render("lexis")
	var resultS string
	var resultRow string
	rowIndex, colIndex, _ := m.game.debugState()
	if m.game.isWon() {
		resultS = fmt.Sprintf("You won in %d/%d attempts!", rowIndex+1, len(m.game.grid.words))
		resultRow = resultBarStyleWin.Render(resultS)
	} else if m.game.isLost() {
		resultS = fmt.Sprintf("Better luck next time! The answer was: %s", m.game.Answer())
		resultRow = resultBarStyleLoss.Render(resultS)
	} else if m.state == stateLoading {
		resultRow = resultBarStyleLoading.Render(m.spinner.View())
	} else if m.state == stateRowNotFull {
		resultRow = resultBarStyleError.Render("Row is not full!")
	} else if m.state == stateInvalidWord {
		resultRow = resultBarStyleError.Render("Invalid word!")

	} else {
		// debug row
		resultS = fmt.Sprintf("Row: %d, Col: %d, RL: %d, L: %c, A: %s",
			rowIndex,
			colIndex,
			len(m.game.grid.words[rowIndex])-1,
			m.game.grid.words[rowIndex][colIndex].r,
			m.game.Answer())
		resultRow = resultBarStyleNormal.Render(resultS)
	}
	helpRow := helpBarStyle.Render(m.help.View(m.keys))
	view := lipgloss.JoinVertical(lipgloss.Center,
		header,
		m.game.grid.render(),
		m.game.keyboard.render(),
		resultRow,
		helpRow)
	v := tea.NewView(containerStyle.Render(view))
	v.WindowTitle = "lexis"
	v.AltScreen = true
	return v
}

// newModel creates a new model with the given logger and initializes the answer provider and spinner
func newModel(logger *log.Logger) model {
	provider := newRandomAnswerProvider()
	s := spinner.New()
	s.Spinner = spinner.Points
	return model{
		game:    newGame(provider, logger),
		log:     logger,
		help:    newHelp(),
		keys:    keys,
		state:   stateLoading,
		spinner: s,
	}
}
