package bubble

import (
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var confirmStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("9")).Bold(true).Padding(1, 2)

type Model struct {
	Current              string
	List                 list.Model
	Width                int
	Heigth               int
	State                confirmationState
	PendingDeleteContext *kubectl.Context
	FocusedButton        int
	DetailedView         string
	ContextMessage       string
}

func New(appTitle string, currentContext string, contexts []list.Item) Model {
	list := list.New(contexts, list.NewDefaultDelegate(), 0, 0)
	list.Title = appTitle

	setupKeyBindings(&list)
	highlightCurrentContext(currentContext, &list)
	return Model{List: list, Current: currentContext, State: normalState}
}

type confirmationState int

const (
	normalState confirmationState = iota
	deleteState
)

type verb string

const (
	deleteContext verb = "delete"
	switchContext verb = "switch"
)

func highlightCurrentContext(currentContext string, l *list.Model) {
	for i, listItem := range l.Items() {
		if ctx, ok := listItem.(*kubectl.Context); ok && ctx.Name == currentContext {
			l.Select(i)
		}
	}
}
