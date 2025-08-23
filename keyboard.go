// keyboard.go provides the keyboard data structure and logic for managing letter states and styles
package main

import (
	"github.com/charmbracelet/lipgloss"
)

type position struct {
	row    int // row index
	column int // column index
}

type keyboardLetter struct {
	position position
	letter   letter
}

type keyboard struct {
	letters map[rune]keyboardLetter // map of letters to their positions and styles
	layout  [3][]keyboardLetter     // keyboard layout with 3 rows
}

// qwertyuiop
// asdfghjkl
// zxcvbnm
var letters = map[rune]keyboardLetter{
	// Row 0: qwertyuiop
	'q': {position: position{row: 0, column: 0}, letter: letter{r: 'q', style: stateStyles[notChecked], state: notChecked}},
	'w': {position: position{row: 0, column: 1}, letter: letter{r: 'w', style: stateStyles[notChecked], state: notChecked}},
	'e': {position: position{row: 0, column: 2}, letter: letter{r: 'e', style: stateStyles[notChecked], state: notChecked}},
	'r': {position: position{row: 0, column: 3}, letter: letter{r: 'r', style: stateStyles[notChecked], state: notChecked}},
	't': {position: position{row: 0, column: 4}, letter: letter{r: 't', style: stateStyles[notChecked], state: notChecked}},
	'y': {position: position{row: 0, column: 5}, letter: letter{r: 'y', style: stateStyles[notChecked], state: notChecked}},
	'u': {position: position{row: 0, column: 6}, letter: letter{r: 'u', style: stateStyles[notChecked], state: notChecked}},
	'i': {position: position{row: 0, column: 7}, letter: letter{r: 'i', style: stateStyles[notChecked], state: notChecked}},
	'o': {position: position{row: 0, column: 8}, letter: letter{r: 'o', style: stateStyles[notChecked], state: notChecked}},
	'p': {position: position{row: 0, column: 9}, letter: letter{r: 'p', style: stateStyles[notChecked], state: notChecked}},

	// Row 1: asdfghjkl
	'a': {position: position{row: 1, column: 0}, letter: letter{r: 'a', style: stateStyles[notChecked], state: notChecked}},
	's': {position: position{row: 1, column: 1}, letter: letter{r: 's', style: stateStyles[notChecked], state: notChecked}},
	'd': {position: position{row: 1, column: 2}, letter: letter{r: 'd', style: stateStyles[notChecked], state: notChecked}},
	'f': {position: position{row: 1, column: 3}, letter: letter{r: 'f', style: stateStyles[notChecked], state: notChecked}},
	'g': {position: position{row: 1, column: 4}, letter: letter{r: 'g', style: stateStyles[notChecked], state: notChecked}},
	'h': {position: position{row: 1, column: 5}, letter: letter{r: 'h', style: stateStyles[notChecked], state: notChecked}},
	'j': {position: position{row: 1, column: 6}, letter: letter{r: 'j', style: stateStyles[notChecked], state: notChecked}},
	'k': {position: position{row: 1, column: 7}, letter: letter{r: 'k', style: stateStyles[notChecked], state: notChecked}},
	'l': {position: position{row: 1, column: 8}, letter: letter{r: 'l', style: stateStyles[notChecked], state: notChecked}},

	// Row 2: zxcvbnm
	'z': {position: position{row: 2, column: 0}, letter: letter{r: 'z', style: stateStyles[notChecked], state: notChecked}},
	'x': {position: position{row: 2, column: 1}, letter: letter{r: 'x', style: stateStyles[notChecked], state: notChecked}},
	'c': {position: position{row: 2, column: 2}, letter: letter{r: 'c', style: stateStyles[notChecked], state: notChecked}},
	'v': {position: position{row: 2, column: 3}, letter: letter{r: 'v', style: stateStyles[notChecked], state: notChecked}},
	'b': {position: position{row: 2, column: 4}, letter: letter{r: 'b', style: stateStyles[notChecked], state: notChecked}},
	'n': {position: position{row: 2, column: 5}, letter: letter{r: 'n', style: stateStyles[notChecked], state: notChecked}},
	'm': {position: position{row: 2, column: 6}, letter: letter{r: 'm', style: stateStyles[notChecked], state: notChecked}},
}

func newKeyboard() keyboard {
	// Initialize the keyboard with the letters and their positions
	layout := [3][]keyboardLetter{}
	layout[0] = make([]keyboardLetter, 10)
	layout[1] = make([]keyboardLetter, 9)
	layout[2] = make([]keyboardLetter, 7)
	for _, kl := range letters {
		layout[kl.position.row][kl.position.column] = kl
	}
	return keyboard{
		letters: letters,
		layout:  layout,
	}
}

func (k *keyboard) updateLetterState(r rune, state int) {
	if kl, exists := k.letters[r]; exists {
		row := kl.position.row
		col := kl.position.column
		style, exists := stateStyles[state]
		if exists {
			k.layout[row][col].letter.style = style // update the style based on the state
			k.layout[row][col].letter.state = state // update the state of the letter
		}
	}
}

func (k *keyboard) getLetterState(r rune) int {
	if kl, exists := k.letters[r]; exists {
		row := kl.position.row
		col := kl.position.column
		return k.layout[row][col].letter.state
	}
	return -1
}

func (k *keyboard) reset() {
	for r, row := range k.layout {
		for c := range row {
			k.layout[r][c].letter.style = stateStyles[notChecked] // reset style to not checked
			k.layout[r][c].letter.state = notChecked              // reset state to not checked
		}
	}
}

func (k *keyboard) render() string {
	rows := make([]string, len(k.layout))
	for i, row := range k.layout {
		letters := make([]string, len(row))
		for j, kl := range row {
			letters[j] = kl.letter.style.Render(string(kl.letter.r))
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, letters...)
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}
