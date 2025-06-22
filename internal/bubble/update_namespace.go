package bubble

import (
	"fmt"
	"log/slog"

	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/ui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateNamespaceState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.LoadingNamespaces {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			if m.namespaceState() && m.LoadingNamespaces {
				var cmd tea.Cmd
				m.Spinner, cmd = m.Spinner.Update(msg)
				return m, cmd
			}

		case namespacesOutput:
			m.NamespaceFetchError = false
			if msg.Err != nil {
				m.NamespaceFetchError = true

				msg.Namespaces = []models.Namespace{
					{Name: msg.Err.Error()},
				}
			}

			items := make([]list.Item, len(msg.Namespaces))
			for i, namespace := range msg.Namespaces {
				items[i] = &models.Namespace{Name: namespace.Name, Age: namespace.Age}
			}

			m.Namespaces.SetItems(items)
			m.LoadingNamespaces = false
			return m, nil

		case tea.KeyMsg:
			if m.namespaceState() && m.LoadingNamespaces {
				switch msg.String() {
				case "esc", "ctrl+c", "q":
					m.State = normalState
					m.LoadingNamespaces = false
				}
				return m, nil
			}
		}

		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		var cmd tea.Cmd
		m.Namespaces, cmd = m.Namespaces.Update(msg)

		switch msg.String() {
		case "enter":
			if m.NamespaceFetchError {
				return m, tea.Quit
			}
			selectedNamespace := m.Namespaces.SelectedItem().(*models.Namespace)
			ctx := m.List.SelectedItem().(*models.Context).Name
			if err := kubectl.SetKubernetesContext(ctx); err != nil {
				return m, tea.Quit
			}
			if err := kubectl.SetDefaultNamespace(selectedNamespace.Name); err != nil {
				slog.Error("update", "message", "settings defualt namespace", "error", err)
				return m, tea.Quit
			}
			m.ContextMessage = ui.SuccessMessage.Render(
				fmt.Sprintf("Context switched to: %s\nNamespace switched to: %s", ctx, selectedNamespace.Name),
			)
			return m, tea.Quit
		case "esc", "ctrl+c", "q":
			m.State = normalState
			m.LoadingNamespaces = false
			return m, nil
		}
		return m, cmd
	}
	return m, nil
}
