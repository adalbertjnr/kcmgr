package bubble

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type keyMap struct {
	Delete          key.Binding
	Enter           key.Binding
	NamespacePrompt key.Binding
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
	NamespacePrompt: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "set ns"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "sel the ctx"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "del the ctx"),
	),
}

func setupKeyBindings(list *list.Model) {
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Delete, keys.Enter, keys.NamespacePrompt}
	}

	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Delete, keys.Enter, keys.NamespacePrompt}
	}

}
