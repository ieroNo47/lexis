package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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
	BorderForeground(lipgloss.Color("#88c0d0")).
	Inherit(defaultStyle)

var exactMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#a3be8c")).
	Foreground(lipgloss.Color("#a3be8c")).
	Inherit(defaultStyle)

var existsMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#ebcb8b")).
	Foreground(lipgloss.Color("#ebcb8b")).
	Inherit(defaultStyle)

var notMatchStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Margin(0).
	BorderForeground(lipgloss.Color("#4c566a")).
	Foreground(lipgloss.Color("#4c566a")).
	Inherit(defaultStyle)

// top bar style
var headerStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#3b4252")).
	Foreground(lipgloss.Color("#88c0d0"))

// result bar styles
var resultBarStyleNormal = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#3b4252")).
	Foreground(lipgloss.Color("#88c0d0"))

var resultBarStyleWin = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#a3be8c")).
	Foreground(lipgloss.Color("#2e3440"))

var resultBarStyleLoss = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#ebcb8b")).
	Foreground(lipgloss.Color("#2e3440"))

var resultBarStyleLoading = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#b48ead")).
	Foreground(lipgloss.Color("#2e3440"))

// status bar
var helpBarStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#3b4252")).
	Foreground(lipgloss.Color("#88c0d0"))

var helpTextStyle = lipgloss.NewStyle()

func updateStyles(msg tea.WindowSizeMsg) int {
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
	resultBarStyleLoading = resultBarStyleLoading.Width(msg.Width - oHorizontal)
	helpBarStyle = helpBarStyle.Width(msg.Width - oHorizontal)
	return oHorizontal
}
