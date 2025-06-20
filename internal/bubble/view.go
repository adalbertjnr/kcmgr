package bubble

import (
	"fmt"

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

		box := confirmStyle.Render(confirmationText, buttons)

		return lipgloss.Place(
			m.Width,
			m.Heigth,
			lipgloss.Center,
			lipgloss.Center,
			box,
		)
	}

	width := (m.Width / 2) - 2
	heigth := m.Heigth - 2

	listView := lipgloss.NewStyle().
		Width(width).
		Height(m.Heigth).
		BorderForeground(lipgloss.Color("8")).
		Render(m.List.View())

	detailedView := lipgloss.NewStyle().
		Width(width).
		Height(heigth).
		Border(lipgloss.ThickBorder()).
		MarginLeft(1).
		BorderForeground(lipgloss.Color("#00FFD7")).
		Render(m.DetailedView)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, listView, detailedView)
}
