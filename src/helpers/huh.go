package helpers

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func HuhTheme() huh.Theme {
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		styles := huh.ThemeBase(isDark)
		styles.Focused.Base = lipgloss.NewStyle().SetString(lipgloss.NewStyle().Render(""))
		styles.Focused.Card = styles.Focused.Base
		return styles
	})
}
