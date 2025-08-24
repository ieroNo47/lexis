// help.go defines the help key bindings and styles
package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Letter  key.Binding
	Delete  key.Binding
	Submit  key.Binding
	Restart key.Binding
	Quit    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Letter, k.Delete, k.Submit, k.Restart, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Letter, k.Delete, k.Submit, k.Restart, k.Quit}}
}

var keys = keyMap{
	Letter: key.NewBinding(
		key.WithKeys("a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"),
		key.WithHelp("a-z", "Set Letter"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace", "delete"),
		key.WithHelp("backspace", "Delete Letter"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Submit"),
	),
	Restart: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "Restart"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "Quit"),
	),
}

func newHelp() help.Model {
	h := help.New()
	h.Styles.Ellipsis = helpTextStyle
	h.Styles.ShortKey = helpTextStyle
	h.Styles.ShortDesc = helpTextStyle
	h.Styles.ShortSeparator = helpTextStyle
	h.Styles.FullKey = helpTextStyle
	h.Styles.FullDesc = helpTextStyle
	h.Styles.FullSeparator = helpTextStyle
	return h
}
