package bubble

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/adalbertjnr/kcmgr/internal/client"
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case normalState:
		return m.UpdateNormalState(msg)
	case deleteState:
		return m.updateDeleteState(msg)
	case namespaceSelectState:
		return m.updateNamespaceState(msg)
	default:
		return m, nil
	}

}

func (m Model) handleContextAction(verb verb, actionFunc func(string) error) (tea.Model, tea.Cmd) {
	if ctx, ok := m.List.SelectedItem().(*models.Context); ok {
		if err := actionFunc(ctx.Name); err != nil {
			log.Printf("Failed to %s context: %v", verb, err)
			return m, tea.Quit
		}
	}

	switch verb {
	case deleteContext:
		contextItems, err := kubectl.KubernetesContexts()
		if err != nil {
			slog.Error("update", "message", "failed refreshing contexts", "error", err)
			return m, nil
		}

		m.List.SetItems(contextItems)
		m.State = normalState
		m.PendingDeleteContext = nil
		return m, nil
	case switchContext:
		m.ContextMessage = ui.SuccessMessage.Render(
			fmt.Sprintf("️️Context switched to: %s", m.List.SelectedItem().(*models.Context).Name),
		)
		return m, tea.Quit
	}

	return m, tea.Quit
}

type NamespacesOutput struct {
	Namespaces []models.Namespace
	Err        error
}

func (m Model) fetchNamespacesCmd(contextName string) tea.Cmd {
	cache, ok := m.NamespaceCache[contextName]

	return func() tea.Msg {
		if ok {
			slog.Info("fetch namespaces", "cache", "hit", "context", contextName)
			return cache
		}

		namespaces, err := client.GetNamespacesByContext(m.KubeConfig, contextName)
		if err != nil {
			slog.Error("update", "message", "fetch namespaces cmd", "error", err)
		}

		messageOutput := NamespacesOutput{Namespaces: namespaces, Err: err}
		m.NamespaceCache[contextName] = messageOutput
		slog.Info("fetch namespaces", "namespaces", namespaces, "cache", "not hit")
		return messageOutput
	}
}
