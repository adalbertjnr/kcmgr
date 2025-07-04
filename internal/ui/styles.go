package ui

import "github.com/charmbracelet/lipgloss"

var (
	ConfirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Bold(true).Padding(1, 2)

	SuccessMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	DetailedViewPadding = lipgloss.NewStyle().
				PaddingTop(2).
				PaddingLeft(1)

	CheckBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("236"))
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

var (
	NamespaceSpiner = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)

	NamespacesLoaded = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(2)
)
