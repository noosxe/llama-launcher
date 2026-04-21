package tui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name        string
	Border      lipgloss.Color
	Header      lipgloss.Color
	Selected    lipgloss.Color
	SelectedBg  lipgloss.Color
	Dim         lipgloss.Color
	StatusBarBg lipgloss.Color
	Text        lipgloss.Color
	Background  lipgloss.Color
	Success     lipgloss.Color
	Error       lipgloss.Color
}

var Themes = []Theme{
	CatppuccinMocha,
	CatppuccinMacchiato,
	CatppuccinFrappe,
	CatppuccinLatte,
	DefaultTheme,
}

var (
	CatppuccinMocha = Theme{
		Name:        "Catppuccin Mocha",
		Border:      lipgloss.Color("#585b70"),
		Header:      lipgloss.Color("#cba6f7"),
		Selected:    lipgloss.Color("#1e1e2e"),
		SelectedBg:  lipgloss.Color("#89b4fa"),
		Dim:         lipgloss.Color("#a6adc8"),
		StatusBarBg: lipgloss.Color("#313244"),
		Text:        lipgloss.Color("#cdd6f4"),
		Background:  lipgloss.Color("#1e1e2e"),
		Success:     lipgloss.Color("#a6e3a1"),
		Error:       lipgloss.Color("#f38ba8"),
	}

	CatppuccinMacchiato = Theme{
		Name:        "Catppuccin Macchiato",
		Border:      lipgloss.Color("#5b6078"),
		Header:      lipgloss.Color("#c6a0f6"),
		Selected:    lipgloss.Color("#24273a"),
		SelectedBg:  lipgloss.Color("#8aadf4"),
		Dim:         lipgloss.Color("#a5adcb"),
		StatusBarBg: lipgloss.Color("#363a4f"),
		Text:        lipgloss.Color("#cad3f5"),
		Background:  lipgloss.Color("#24273a"),
		Success:     lipgloss.Color("#a6da95"),
		Error:       lipgloss.Color("#ed8796"),
	}

	CatppuccinFrappe = Theme{
		Name:        "Catppuccin Frappe",
		Border:      lipgloss.Color("#626880"),
		Header:      lipgloss.Color("#ca9ee6"),
		Selected:    lipgloss.Color("#303446"),
		SelectedBg:  lipgloss.Color("#8caaee"),
		Dim:         lipgloss.Color("#a5adce"),
		StatusBarBg: lipgloss.Color("#414559"),
		Text:        lipgloss.Color("#c6d0f5"),
		Background:  lipgloss.Color("#303446"),
		Success:     lipgloss.Color("#a6d189"),
		Error:       lipgloss.Color("#e78284"),
	}

	CatppuccinLatte = Theme{
		Name:        "Catppuccin Latte",
		Border:      lipgloss.Color("#bcc0cc"),
		Header:      lipgloss.Color("#8839ef"),
		Selected:    lipgloss.Color("#eff1f5"),
		SelectedBg:  lipgloss.Color("#1e66f5"),
		Dim:         lipgloss.Color("#5c5f77"),
		StatusBarBg: lipgloss.Color("#e6e9ef"),
		Text:        lipgloss.Color("#4c4f69"),
		Background:  lipgloss.Color("#eff1f5"),
		Success:     lipgloss.Color("#40a02b"),
		Error:       lipgloss.Color("#d20f39"),
	}

	DefaultTheme = Theme{
		Name:        "Default Dark",
		Border:      lipgloss.Color("238"),
		Header:      lipgloss.Color("205"),
		Selected:    lipgloss.Color("255"),
		SelectedBg:  lipgloss.Color("63"),
		Dim:         lipgloss.Color("245"),
		StatusBarBg: lipgloss.Color("235"),
		Text:        lipgloss.Color("252"),
		Background:  lipgloss.Color("232"),
		Success:     lipgloss.Color("42"),
		Error:       lipgloss.Color("196"),
	}
)

type Styles struct {
	SidebarTitle   lipgloss.Style
	Item           lipgloss.Style
	SelectedItem   lipgloss.Style
	Status         lipgloss.Style
	Header         lipgloss.Style
	OverlayBox     lipgloss.Style
	OverlayTitle   lipgloss.Style
	OverlayItem    lipgloss.Style
	OverlaySelItem lipgloss.Style
	Divider        lipgloss.Style
    LogText        lipgloss.Style
}

func MakeStyles(t Theme) Styles {
	return Styles{
		SidebarTitle: lipgloss.NewStyle().
			Foreground(t.Header).
			Bold(true).
			Padding(0, 1),

		Item: lipgloss.NewStyle().
			Foreground(t.Dim).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Foreground(t.Selected).
			Background(t.SelectedBg).
			Bold(true).
			Padding(0, 1),

		Status: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.StatusBarBg).
			Padding(0, 1),

		Header: lipgloss.NewStyle().
			Foreground(t.Header).
			Bold(true),

		OverlayBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Background(t.Background).
			Padding(1, 2),

		OverlayTitle: lipgloss.NewStyle().
			Foreground(t.Header).
			Bold(true).
			MarginBottom(1),

		OverlayItem: lipgloss.NewStyle().
			Foreground(t.Text).
			Padding(0, 1),

		OverlaySelItem: lipgloss.NewStyle().
			Foreground(t.Selected).
			Background(t.SelectedBg).
			Padding(0, 1),

		Divider: lipgloss.NewStyle().
			Foreground(t.Border),

        LogText: lipgloss.NewStyle().
            Foreground(t.Dim).
            PaddingLeft(1),
	}
}
