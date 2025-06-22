package bubble

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/adalbertjnr/kcmgr/internal/client"
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/ui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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

	case spinner.TickMsg:
		if m.namespaceState() && m.LoadingNamespaces {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}

	case namespacesOutput:
		if m.namespaceState() && m.LoadingNamespaces {

			if msg.Err != nil {
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
		}

	case tea.KeyMsg:
		if m.namespaceState() && m.LoadingNamespaces {
			switch msg.String() {
			case "esc", "ctrl+c", "q":
				m.State = normalState
				m.LoadingNamespaces = false
				return m, nil
			}
			return m, nil
		}

		if !m.LoadingNamespaces && m.namespaceState() {
			var cmd tea.Cmd
			m.Namespaces, cmd = m.Namespaces.Update(msg)

			switch msg.String() {
			case "enter":
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
				return m, cmd
			}

			return m, cmd
		}

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

type namespacesOutput struct {
	Namespaces []models.Namespace
	Err        error
}

func (m Model) fetchNamespacesCmd(contextName string) tea.Cmd {
	return func() tea.Msg {
		namespaces, err := client.GetNamespacesByContext(m.KubeConfig, contextName)
		if err != nil {
			slog.Error("update", "message", "fetch namespaces cmd", "error", err)
		}
		slog.Info("fetch namespaces", "namespaces", namespaces)
		return namespacesOutput{Namespaces: namespaces, Err: err}
	}
}
