// grid.go provides the grid data structure and logic for managing letter states and styles
package main

import "github.com/charmbracelet/lipgloss"

// consts for states of a letter, matched, exists, not matched
// currently helps when checking the letter state over different iterations
const (
	matched    = iota // letter is in the correct position
	exists            // letter is in the word but not in the correct position
	notMatched        // letter is not in the word
	notChecked        // letter has not been checked yet
)

// map to get string representation of each state
var states = map[int]string{
	matched:    "matched",
	exists:     "exists",
	notMatched: "notMatched",
	notChecked: "notChecked",
}

// match a state to a style
var stateStyles = map[int]lipgloss.Style{
	matched:    exactMatchStyle,
	exists:     existsMatchStyle,
	notMatched: notMatchStyle,
	notChecked: defaultStyle,
}

// letter represents a single letter and its style
type letter struct {
	r     rune
	style lipgloss.Style
	state int
}

// word represents a row of letters
type word []letter

// grid represents the grid where each row is a word
type grid struct {
	words    []word
	colIndex int
	rowIndex int
}

// Initializes a new grid with the specified number of rows and columns
func newGrid(rows, cols int) grid {
	grid := grid{
		words:    make([]word, rows),
		colIndex: 0,
		rowIndex: 0,
	}
	for i := range grid.words {
		grid.words[i] = make([]letter, cols)
		for j := range grid.words[i] {
			grid.words[i][j] = letter{r: ' ', style: stateStyles[notChecked], state: notChecked}
		}
	}
	return grid
}

// setLetter sets a letter in the grid at the next available position
func (g *grid) setLetter(r rune) {
	if g.rowIndex < len(g.words) && g.colIndex < len(g.words[g.rowIndex]) {
		g.words[g.rowIndex][g.colIndex].r = r
		// move to the next column if we're not at the end of the row
		if g.colIndex < len(g.words[g.rowIndex])-1 {
			g.colIndex++ // move to the next column
			g.updateActiveCell()
		}
	}
}

// deleteLetter deletes the last letter in the current row
func (g *grid) deleteLetter() {
	// we move back only if we're not at the first column and the current letter is not empty
	if g.colIndex > 0 {
		if g.words[g.rowIndex][g.colIndex].r == ' ' {
			g.colIndex--
			g.updateActiveCell()
		}
	}
	g.words[g.rowIndex][g.colIndex].r = ' ' // delete the letter
}

// rowFull checks if the current row is full (i.e., all letters are filled)
func (g *grid) rowFull() bool {
	return g.colIndex == len(g.words[g.rowIndex])-1 && g.words[g.rowIndex][g.colIndex].r != ' '
}

// goToNextRow moves to the next row in the grid if possible
func (g *grid) goToNextRow() bool {
	if g.rowIndex < len(g.words)-1 {
		g.rowIndex++   // move to the next row
		g.colIndex = 0 // reset column index
		g.updateActiveCell()
		return true
	} else {
		return false // no more rows
	}
}

// updateStyle updates the style of a letter at a specific position in the grid
func (g *grid) updateStyle(row, col int, style lipgloss.Style) {
	if row < len(g.words) && col < len(g.words[row]) {
		g.words[row][col].style = style
	}
}

// updateActiveCell updates the style of the active cell in the current row
func (g *grid) updateActiveCell() {
	for i := range g.words[g.rowIndex] {
		if i == g.colIndex {
			g.updateStyle(g.rowIndex, i, activeStyle) // set the active style for the current cell
		} else {
			g.updateStyle(g.rowIndex, i, defaultStyle) // reset style for other cells in the active row
		}
	}
}

// updateState updates the state of a letter at a specific position in the grid
func (g *grid) updateState(row, col int, state int) {
	if row < len(g.words) && col < len(g.words[row]) {
		style, exists := stateStyles[state]
		if exists {
			g.words[row][col].style = style // update the style based on the state
			g.words[row][col].state = state // update the state of the letter
		}
	}
}

// reset resets the grid to its initial state
func (g *grid) reset() {
	for i := range g.words {
		for j := range g.words[i] {
			g.words[i][j] = letter{r: ' ', style: stateStyles[notChecked], state: notChecked}
		}
	}
	g.colIndex = 0
	g.rowIndex = 0
	g.updateActiveCell() // reset the active cell style
}

// render renders the grid as a string, with each letter styled according to its state
func (g *grid) render() string {
	rows := make([]string, 0, len(g.words))
	for _, w := range g.words {
		letters := []string{}
		for _, l := range w {
			// render each letter with its style
			letters = append(letters, l.style.Render(string(l.r)))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, letters...))
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}
