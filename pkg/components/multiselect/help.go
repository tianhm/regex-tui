package multiselect

import "github.com/charmbracelet/bubbles/v2/key"

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("space", "x"),
		key.WithHelp("space/x", "select"),
	),
}

var upDown = key.NewBinding(
	key.WithKeys("up", "down", "k", "j"),
	key.WithHelp("↑/↓", "navigate"),
)

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{upDown, k.Toggle},
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{upDown, k.Toggle}
}
