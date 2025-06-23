package bubble

import (
	"fmt"

	"github.com/adalbertjnr/kcmgr/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.deleteState() {
		confirmationText := fmt.Sprintf(
			"Are you sure you want to delete context '%s'?\n\n"+
				"Press [enter] or [y] to confirm\n\n"+
				"or\n\n"+
				"Press [q], [esc] or [ctrl+c] to cancel", m.PendingDeleteContext.Name,
		)

		confirmButton := button("Confirm", m.FocusedButton == 1)
		cancelButton := button("Cancel", m.FocusedButton == 0)

		buttons := lipgloss.JoinHorizontal(
			lipgloss.Top,
			confirmButton,
			cancelButton,
		)

		box := ui.ConfirmStyle.Render(confirmationText, buttons)

		return lipgloss.Place(
			m.Width,
			m.Heigth,
			lipgloss.Center,
			lipgloss.Center,
			box,
		)
	}

	if m.namespaceState() {
		if m.LoadingNamespaces {
			msg := fmt.Sprintf("%s Loading namespaces for %s...", m.Spinner.View(), m.TargetContext)
			spinnerView := ui.NamespaceSpiner.Render(msg)
			return lipgloss.Place(m.Width, m.Heigth, lipgloss.Center, lipgloss.Center, spinnerView)
		} else {
			namespacesView := ui.NamespacesLoaded.Render(m.Namespaces.View())
			return lipgloss.Place(m.Width, m.Heigth, lipgloss.Center, lipgloss.Center, namespacesView)
		}
	}

	width := (m.Width / 2) - 2
	heigth := m.Heigth - 2

	listView := lipgloss.NewStyle().
		Width(width).
		Height(m.Heigth).
		Render(m.List.View())

	detailedView := lipgloss.NewStyle().
		Width(width).
		Height(heigth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFD7")).
		Render(m.DetailedView)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, listView, detailedView)
}

func (m Model) deleteState() bool {
	return m.State == deleteState
}
func (m Model) namespaceState() bool {
	return m.State == namespaceSelectState
}
