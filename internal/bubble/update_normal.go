package bubble

import (
	"fmt"

	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/ui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) UpdateNormalState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()

		m.Width = msg.Width
		m.Heigth = msg.Height

		m.List.SetSize(msg.Width-h, msg.Height-v)

		if m.DetailedView == "" {
			if ctx, ok := m.List.SelectedItem().(*models.Context); ok {
				detailedContext, err := kubectl.GetRawContext(ctx.Context.Cluster)
				if err != nil {
					m.DetailedView = ui.DetailedViewPadding.Render(
						fmt.Sprintf("Error loading context details: %v", err),
					)
					return m, nil
				}

				m.DetailedView = ui.DetailedViewPadding.Render(detailedContext)
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Delete):
			if ctx, ok := m.List.SelectedItem().(*models.Context); ok {
				m.State = deleteState
				m.PendingDeleteContext = ctx
				return m, nil
			}
		case key.Matches(msg, keys.Enter):
			return m.handleContextAction(
				switchContext,
				kubectl.SetKubernetesContext,
			)
		case key.Matches(msg, keys.NamespacePrompt):
			if ctx, ok := m.List.SelectedItem().(*models.Context); ok {
				m.State = namespaceSelectState
				m.TargetContext = ctx.Name
				m.LoadingNamespaces = true
				m.Namespaces.SetItems([]list.Item{})
				return m, tea.Batch(m.Spinner.Tick, m.fetchNamespacesCmd(ctx.Name))
			}
		}
	}

	var cmd tea.Cmd
	previousIndex := m.List.Index()
	m.List, cmd = m.List.Update(msg)

	if m.List.Index() != previousIndex {
		if ctx, ok := m.List.SelectedItem().(*models.Context); ok {
			detailedContext, err := kubectl.GetRawContext(ctx.Context.Cluster)
			if err != nil {
				m.DetailedView = ui.DetailedViewPadding.Render(
					fmt.Sprintf("Error loading context details: %v", err),
				)
				return m, cmd
			}

			m.DetailedView = ui.DetailedViewPadding.Render(detailedContext)
		}
	}
	return m, cmd
}
