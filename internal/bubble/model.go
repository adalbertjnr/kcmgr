package bubble

import (
	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
)

type Model struct {
	Current                string
	List                   list.Model
	Width                  int
	Heigth                 int
	State                  confirmationState
	PendingDeleteContext   *models.Context
	FocusedButton          int
	DetailedView           string
	ContextMessage         string
	Namespaces             list.Model
	SelectedNamespaceIndex int
	TargetContext          string
	Spinner                spinner.Model
	KubeConfig             string
	LoadingNamespaces      bool
	NamespaceFetchError    bool
}

func New(contextsTitle, namespacesTitle, kubeconfig string, currentContext string, contexts []list.Item) Model {
	contextsList := list.New(contexts, list.NewDefaultDelegate(), 0, 0)
	contextsList.Title = contextsTitle

	setupKeyBindings(&contextsList)
	highlightCurrentContext(currentContext, &contextsList)

	namespaceList := list.New([]list.Item{}, list.NewDefaultDelegate(), 75, 25)
	namespaceList.Title = namespacesTitle
	namespaceList.SetFilteringEnabled(true)
	namespaceList.SetShowFilter(true)

	sp := spinner.New()
	// sp.Spinner = spinner.Pulse
	// sp.Spinner = spinner.Monkey
	// sp.Spinner = spinner.Dot
	// sp.Spinner = spinner.Globe
	// sp.Spinner = spinner.Hamburger
	// sp.Spinner = spinner.Jump
	// sp.Spinner = spinner.Line
	// sp.Spinner = spinner.Meter
	// sp.Spinner = spinner.MiniDot
	// sp.Spinner = spinner.Moon
	sp.Spinner = spinner.Points
	return Model{
		List:       contextsList,
		Current:    currentContext,
		State:      normalState,
		Namespaces: namespaceList,
		Spinner:    sp,
		KubeConfig: kubeconfig,
	}
}

type confirmationState int

const (
	normalState confirmationState = iota
	deleteState
	namespaceLoadingState
	namespaceSelectState
)

type verb string

const (
	deleteContext    verb = "delete"
	switchContext    verb = "switch"
	defaultNamespace verb = "defaultNamespace"
)

func highlightCurrentContext(currentContext string, l *list.Model) {
	for i, listItem := range l.Items() {
		if ctx, ok := listItem.(*models.Context); ok && ctx.Name == currentContext {
			l.Select(i)
		}
	}
}
