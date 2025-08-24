package main

import "github.com/charmbracelet/lipgloss"

// parent container style
var containerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0, 0, 0, 0).
	Align(lipgloss.Center, lipgloss.Center)

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
	Foreground(lipgloss.Color("#7d7")).
	Inherit(defaultStyle)

var existsMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#cc0")).
	Foreground(lipgloss.Color("#cc0")).
	Inherit(defaultStyle)

var notMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#444")).
	Foreground(lipgloss.Color("#444")).
	Inherit(defaultStyle)

// top bar style
var headerStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#8af"))

// result bar style
var resultBarStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#8af"))

// status bar
var helpBarStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#8af"))

var helpTextStyle = lipgloss.NewStyle()
