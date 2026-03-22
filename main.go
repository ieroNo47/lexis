// lexis app
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	game game
	log  *log.Logger
	help help.Model
	keys keyMap
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
		resultBarStyleNormal = resultBarStyleNormal.Width(msg.Width - oHorizontal)
		resultBarStyleWin = resultBarStyleWin.Width(msg.Width - oHorizontal)
		resultBarStyleLoss = resultBarStyleLoss.Width(msg.Width - oHorizontal)
		helpBarStyle = helpBarStyle.Width(msg.Width - oHorizontal)
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
	}
	// if log level is debug, print the current string in the active row
	if m.log.GetLevel() == log.DebugLevel {
		rowIndex, colIndex, rowString := m.game.debugState()
		m.log.Debug("", "row", rowIndex, "col", colIndex, "string", rowString)
	}
	return m, nil
}

func (m model) View() tea.View {
	// todo: Move logit to update header/footer to Update()
	// header
	header := headerStyle.Render("lexis")
	var resultS string
	var resultRow string
	rowIndex, colIndex, _ := m.game.debugState()
	if m.game.isWon() {
		resultS = fmt.Sprintf("You won! %d/%d attempts", rowIndex+1, len(m.game.grid.words))
		resultRow = resultBarStyleWin.Render(resultS)
	} else if m.game.isLost() {
		resultS = fmt.Sprintf("Better luck next time! The answer was: %s", m.game.Answer())
		resultRow = resultBarStyleLoss.Render(resultS)
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
	view := lipgloss.JoinVertical(lipgloss.Center, header, m.game.grid.render(), m.game.keyboard.render(), resultRow, helpRow)
	v := tea.NewView(containerStyle.Render(view))
	v.WindowTitle = "lexis"
	v.AltScreen = true
	return v
}

func main() {
	if err := os.Remove("debug.log"); err != nil && !os.IsNotExist(err) {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
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
	answer := []rune("vogue")
	p := tea.NewProgram(model{
		game: newGame(answer, logger),
		log:  logger,
		help: newHelp(),
		keys: keys,
	})

	logger.Info("==== Starting lexis ====")
	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
