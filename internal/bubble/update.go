package bubble

import (
	"fmt"
	"log"

	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/ui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()

		m.Width = msg.Width
		m.Heigth = msg.Height

		m.List.SetSize(msg.Width-h, msg.Height-v)

		if m.DetailedView == "" {
			if ctx, ok := m.List.SelectedItem().(*kubectl.Context); ok {
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
		if m.deleteState() {
			switch msg.String() {
			case "left", "h", "tab":
				m.FocusedButton = 1
				return m, nil
			case "right", "l":
				m.FocusedButton = 0
				return m, nil
			case "enter", "y", "Y":
				if m.FocusedButton == 1 {
					return m.handleContextAction(
						deleteContext,
						kubectl.DeleteKubernetesContext,
					)
				}
				m.State = normalState
				m.PendingDeleteContext = nil
			case "esc", "ctrl+c", "q":
				m.State = normalState
				m.PendingDeleteContext = nil
				return m, nil
			}
			return m, nil
		}
		switch {
		case key.Matches(msg, keys.Delete):
			if ctx, ok := m.List.SelectedItem().(*kubectl.Context); ok {
				m.State = deleteState
				m.PendingDeleteContext = ctx
				return m, nil
			}
		case key.Matches(msg, keys.Enter):
			return m.handleContextAction(
				switchContext,
				kubectl.SetKubernetesContext,
			)
		}
	}

	var cmd tea.Cmd
	previousIndex := m.List.Index()
	m.List, cmd = m.List.Update(msg)

	if m.List.Index() != previousIndex {
		if ctx, ok := m.List.SelectedItem().(*kubectl.Context); ok {
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

func (m Model) handleContextAction(verb verb, actionFunc func(string) error) (tea.Model, tea.Cmd) {
	if ctx, ok := m.List.SelectedItem().(*kubectl.Context); ok {
		if err := actionFunc(ctx.Name); err != nil {
			log.Printf("Failed to %s context: %v", verb, err)
			return m, tea.Quit
		}
	}

	switch verb {
	case deleteContext:
		contextItems, err := kubectl.KubernetesContexts()
		if err != nil {
			log.Printf("Failed refreshing contexts: %v", err)
			return m, nil
		}

		m.List.SetItems(contextItems)
		m.State = normalState
		m.PendingDeleteContext = nil
		return m, nil
	case switchContext:
		m.ContextMessage = ui.SuccessMessage.Render(fmt.Sprintf("️️⚙️ %sContext switched to: %s", " ", m.List.SelectedItem().(*kubectl.Context).Name))
		return m, tea.Quit
	}

	return m, tea.Quit
}

func (m Model) deleteState() bool {
	return m.State == deleteState
}
