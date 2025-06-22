package bubble

import (
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateDeleteState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		}
	}
	return m, nil
}
