package bubble

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type keyMap struct {
	Delete key.Binding
	Enter  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Delete,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Delete},
	}
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select the context"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete a context"),
	),
}

func setupKeyBindings(list *list.Model) {
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Delete, keys.Enter}
	}

	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Delete, keys.Enter}
	}
}
