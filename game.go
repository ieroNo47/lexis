// game.go manages the core game logic
package main

import (
	"slices"

	"github.com/charmbracelet/log"
)

const (
	// game states
	playing = iota
	won
	lost
)

type game struct {
	grid        grid
	answer      []rune
	state       int
	usedLetters []letter
	log         *log.Logger
}

func newGame(answer []rune, log *log.Logger) game {
	grid := newGrid(6, 5)
	grid.updateStyle(0, 0, activeStyle) // set the first cell as active
	return game{
		grid:   grid,
		answer: answer,
		log:    log,
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

func (g game) isOver() bool {
	return g.state == won || g.state == lost
}

func (g *game) Answer() string {
	return string(g.answer)
}

func (g *game) processLetter(text string) {
	if !g.isOver() && len(text) > 0 {
		g.grid.setLetter([]rune(text)[0])
	}
}

func (g *game) processDelete() {
	if !g.isOver() {
		g.grid.deleteLetter()
	}
}

func (g *game) processSubmit() {
	if !g.grid.rowFull() {
		return // do nothing if the row is not full
	}
	if g.isMatch() {
		g.log.Info("Match found")
		g.state = won // mark the game as won
		g.log.Debug("Marking game as finished.", "reason", "win")
		//
		return
	}
	// TODO: keyboard state updates go here
	if !g.grid.goToNextRow() {
		g.state = lost // mark the game as lost if there are no more rows
		g.log.Debug("Marking game as finished.", "reason", "no more rows")
	}
}

func (g *game) reset() {
	g.grid.reset()
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
