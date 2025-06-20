package bubble

import "github.com/adalbertjnr/kcmgr/internal/ui"

func button(label string, focused bool) string {
	base := ui.Button

	if focused {
		base = ui.ButtonFocused
	}

	return base.Render(label)
}
