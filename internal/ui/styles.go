package ui

import "github.com/charmbracelet/lipgloss"

var (
	SuccessMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	DetailedViewPadding = lipgloss.NewStyle().PaddingTop(3)
)

var (
	Button = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		MarginTop(2)

	ButtonFocused = Button.Foreground(lipgloss.Color("30")).
			Background(lipgloss.Color("9")).
			Bold(true)
)
