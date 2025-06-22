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

		case NamespacesOutput:
			m.NamespaceFetchError = false
			m.Namespaces.SetSize(NAMESPACE_PANEL_WIDTH-40, NAMESPACE_PANEL_HEIGHT)
			slog.Info("Resizing", "window width", NAMESPACE_PANEL_WIDTH-40, "window heigth", NAMESPACE_PANEL_HEIGHT)
			if msg.Err != nil {
				m.NamespaceFetchError = true

				msg.Namespaces = []models.Namespace{
					{Name: msg.Err.Error()},
				}
			}

			items := make([]list.Item, len(msg.Namespaces))
			for i := range msg.Namespaces {
				items[i] = &models.Namespace{Name: msg.Namespaces[i].Name, Age: msg.Namespaces[i].Age}
			}

			m.Namespaces.SetItems(items)
			m.LoadingNamespaces = false
			m.Namespaces.SetFilteringEnabled(true)
			m.Namespaces.SetShowFilter(true)
			for _, i := range items {
				ns := i.(*models.Namespace)
				slog.Info("namespace item", "name", ns.Name, "filter", ns.FilterValue())
			}

			for i, item := range m.Namespaces.Items() {
				ns, ok := item.(*models.Namespace)
				if !ok {
					slog.Warn("item is not namespace", "index", i, "type", fmt.Sprintf("%T", item))
				} else {
					slog.Info("item filter value", "index", i, "filter", ns.FilterValue())
				}
			}
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

			slog.Info("selected namespace", "name", selectedNamespace.Name)

			if err := kubectl.SetDefaultNamespace(selectedNamespace.Name); err != nil {
				slog.Error("update", "message", "settings defualt namespace", "error", err)
				return m, tea.Quit
			}
			m.ContextMessage = ui.SuccessMessage.Render(
				fmt.Sprintf("Context switched to: %s\nNamespace switched to: %s", ctx, selectedNamespace.Name),
			)
			return m, tea.Quit
		case "esc", "ctrl+c":
			if m.Namespaces.FilterState() == list.Filtering {
				m.Namespaces.ResetFilter()
				return m, nil
			}
			m.State = normalState
			m.LoadingNamespaces = false
			return m, nil
		case "q":
			if m.Namespaces.FilterState() != list.Filtering {
				m.State = normalState
				m.LoadingNamespaces = false
				return m, nil
			}
		}
	}

	slog.Info("msg", "type", fmt.Sprintf("%T", msg))
	var cmd tea.Cmd
	m.Namespaces, cmd = m.Namespaces.Update(msg)
	vi := m.Namespaces.VisibleItems()
	fs := m.Namespaces.FilterState()
	fv := m.Namespaces.FilterValue()
	slog.Info("visible items", "count", len(vi))
	slog.Info("filter state", "state", fs.String())
	slog.Info("filter value", "value", fv)
	return m, cmd
}
