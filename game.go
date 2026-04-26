// game.go manages the core game logic
package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/charmbracelet/log"
)

const (
	// game states
	loading = iota
	playing
	won
	lost
)

var (
	ErrInvalidWord = errors.New("invalid word")
	ErrRowNotFull  = errors.New("row not full")
)

type game struct {
	grid           grid
	keyboard       keyboard
	answerProvider answerProvider
	answer         []rune
	state          int
	tempWord       tempWord
	log            *log.Logger
}

func newGame(ap answerProvider, log *log.Logger) game {
	grid := newGrid(6, 5)
	grid.updateStyle(0, 0, activeStyle) // set the first cell as active
	return game{
		grid:           grid,
		keyboard:       newKeyboard(),
		answerProvider: ap,
		answer:         []rune{},
		state:          loading,
		log:            log,
	}
}

func (g game) isMatch() bool {
	guess := g.grid.words[g.grid.rowIndex]
	for i, l := range guess {
		if l.r != g.answer[i] {
			return false
		}
	}
	return true
}

func (g game) isLost() bool {
	return g.state == lost
}

func (g game) isWon() bool {
	return g.state == won
}

func (g game) inProgress() bool {
	return g.state == playing
}

func (g game) rowString() string {
	var rowString string
	for _, l := range g.grid.words[g.grid.rowIndex] {
		rowString += string(l.r)
	}
	return rowString
}

func (g game) Answer() string {
	return string(g.answer)
}

func (g *game) initProvider() {
	g.log.Debug("Initializing answer provider")
	g.answerProvider.init()
	g.log.Debug("Answer provider initialized")
}

func (g *game) start() {
	g.log.Debug("== Starting Game ==")
	answer := g.answerProvider.getAnswer()
	g.answer = []rune(answer)
	g.log.Debug("Answer", "answer", string(g.answer))
	g.state = playing
}

func (g *game) processLetter(text string) {
	if g.inProgress() && len(text) > 0 {
		g.grid.setLetter([]rune(text)[0])
	}
}

func (g *game) processDelete() {
	if g.inProgress() {
		g.grid.deleteLetter()
	}
}

// rowReady checks if the current row is ready to be submitted and returns an error if it is not
// may perform I/O operations such as validating the word with the answer provider, so it should be called in a separate goroutine to avoid blocking the UI
func (g *game) rowReady() error {
	if !g.grid.rowFull() {
		g.log.Info("Row is not full, cannot submit")
		return fmt.Errorf("cannot submit: %w", ErrRowNotFull)
	}
	if !g.answerProvider.validWord(g.rowString()) {
		g.log.Info("Invalid word submitted", "word", g.rowString())
		return fmt.Errorf("cannot submit: %w", ErrInvalidWord)
	}
	return nil
}

// Submit processes current row and updates letter states based on the answer.
// It also updates the game state to won or lost if applicable.
func (g *game) Submit() {
	if g.isMatch() {
		g.log.Info("Match found")
		g.state = won // mark the game as won
		g.log.Debug("Marking game as finished.", "reason", "win")
		for i, l := range g.grid.words[g.grid.rowIndex] {
			g.grid.updateState(g.grid.rowIndex, i, matched)
			g.keyboard.updateLetterState(l.r, matched) // update the keyboard state
		}
		return
	}
	tw := newTempWord(g.answer)
	// first pass: check for exact matches
	g.log.Debug("== First Pass: Exact Matches ==")
	for i, l := range g.grid.words[g.grid.rowIndex] {
		// change style based on match
		if l.r == g.answer[i] {
			g.log.Debug("Update", "letter", string(l.r), "index", i, "state", states[matched])
			g.grid.updateState(g.grid.rowIndex, i, matched) // mark the letter as matched
			g.keyboard.updateLetterState(l.r, matched)      // update the keyboard state
			tw = tw.remove(l.r)                             // remove the letter from the temporary word
		}
	}
	// 	// second pass: check for exists matches and not matches
	// 	// having a separate pass for exists matches allows us to not mark a letter as exists if it was already matched
	var targetState = notMatched
	g.log.Debug("== Second Pass: Exists and Not Matches ==")
	for i, l := range g.grid.words[g.grid.rowIndex] {
		if tw.has(l.r) && l.state != matched {
			targetState = exists
			g.log.Debug("Row Update", "letter", string(l.r), "index", i, "state", states[exists])
			g.grid.updateState(g.grid.rowIndex, i, exists) // mark the letter as exists
			tw = tw.remove(l.r)                            // remove the letter from the temporary word
		} else if l.state != matched {
			g.log.Debug("Row Update", "letter", string(l.r), "index", i, "state", "notMatched")
			g.grid.updateState(g.grid.rowIndex, i, notMatched) // mark the letter as not matched
		}
		if g.keyboard.getLetterState(l.r) != matched { // only update the keyboard state if it is not already matched
			g.log.Debug("Keyboard Update", "letter", string(l.r), "to", states[targetState], "from", states[g.keyboard.getLetterState(l.r)])
			g.keyboard.updateLetterState(l.r, targetState)
		}
		targetState = notMatched // reset target state
	}
	if g.grid.goToNextRow() {
		g.log.Info("Moving to next row")
		g.tempWord = make(tempWord, len(g.answer))
	} else {
		g.state = lost // mark the game as lost if there are no more rows
		g.log.Debug("Marking game as finished.", "reason", "no more rows")
		// TODO: return as custom error to handle in model?
	}
}

func (g *game) reset() {
	g.grid.reset()
	g.keyboard.reset()
	g.state = playing
}

func (g game) debugState() (int, int, string) {
	var currentRow string
	for _, l := range g.grid.words[g.grid.rowIndex] {
		currentRow += string(l.r)
	}
	return g.grid.rowIndex, g.grid.colIndex, currentRow
}

// tempWord is a slice alias for []rune that provides methods to check for existence and remove letters.
// Used to keep track of letters that are still to be matched.
type tempWord []rune

func newTempWord(answer []rune) tempWord {
	tw := make(tempWord, len(answer))
	copy(tw, answer)
	return tw
}

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
