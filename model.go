// lexis app
package main

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	game           game
	log            *log.Logger
	help           help.Model
	keys           keyMap
	loading        bool
	answerProvider answerProvider
	spinner        spinner.Model
}

// initCompleteMsg is a message that is sent when the game initialization is complete and the answer is ready
type initCompleteMsg struct {
	answer string
}

// Init initializes the model and starts the game by getting the answer from the answer provider
func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
	}

	answerCmd := func() tea.Msg {
		m.log.Debug("[Init]")
		answer := m.answerProvider.getAnswer()
		return initCompleteMsg{answer}
	}
	cmds = append(cmds, answerCmd)
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	m.log.Debug("[Update]", "msg", spew.Sdump(msg))
	switch msg := msg.(type) {
	case initCompleteMsg:
		m.log.Debug("Initialization complete")
		m.game.init(msg.answer)
		m.loading = false
		return m, nil
	case tea.WindowSizeMsg:
		oHorizontal := updateStyles(msg)
		m.help.SetWidth(msg.Width - oHorizontal)
		m.log.Debug("Window resized", "width", msg.Width, "height", msg.Height)
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.log.Info("==== Bye! ====")
			return m, tea.Quit
		case key.Matches(msg, m.keys.Letter):
			m.game.processLetter(msg.Text)
		case key.Matches(msg, m.keys.Delete):
			m.game.processDelete()
		case key.Matches(msg, m.keys.Restart):
			m.log.Info("==== Restarting game ====")
			m.game.reset()
			return m, nil
		case key.Matches(msg, m.keys.Submit):
			m.game.processSubmit()
		}
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
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
	} else if m.loading {
		resultRow = resultBarStyleLoading.Render(m.spinner.View())
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
		game:           newGame(logger),
		log:            logger,
		help:           newHelp(),
		keys:           keys,
		loading:        true,
		answerProvider: provider,
		spinner:        s,
	}
}
