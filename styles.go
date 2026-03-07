package main

import "charm.land/lipgloss/v2"

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

// status bar
var helpBarStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#3b4252")).
	Foreground(lipgloss.Color("#88c0d0"))

var helpTextStyle = lipgloss.NewStyle()
