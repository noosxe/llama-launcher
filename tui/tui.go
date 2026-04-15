package tui

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"llama-launcher/config"
)

const sidebarWidth = 35

type ViewState int

const (
	MainView ViewState = iota
	MenuView
	SettingsView
	AlertView
	ConfirmView
)

type Container struct {
	Cmd       *exec.Cmd
	LogLines  []string
	IsRunning bool
}

type logMsg struct {
	modelName string
	line      string
}

type containerExitMsg struct {
	modelName string
	err       error
}

type model struct {
	config *config.Config
	theme  Theme
	styles Styles
	state  ViewState
	width  int
	height int

	// Main View state
	selectedModel int
	containers    map[string]*Container
	logsViewport  viewport.Model

	// Menu View state
	menuOptions     []string
	selectedMenuIdx int

	// Settings View state
	selectedThemeIdx int

	// Alert View state
	alertMsg string

	// Confirm View state
	pendingStop *config.Model

	// Live Stats
	stats MachineStats
}

func Run(cfgFile string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	theme := CatppuccinMocha
	vp := viewport.New(0, 0)

	m := model{
		config:           cfg,
		theme:            theme,
		styles:           MakeStyles(theme),
		state:            MainView,
		containers:       make(map[string]*Container),
		logsViewport:     vp,
		menuOptions:      []string{"Settings", "Exit"},
		selectedThemeIdx: 0,
		stats:            MachineStats{CPU: "...", RAM: "...", GPU: "...", VRAM: "..."},
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	program = p
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run tui: %w", err)
	}

	return nil
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{tea.WindowSize(), tea.EnterAltScreen, fetchStatsCmd, tickStatsCmd()}

	out, err := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err == nil {
		runningNames := strings.Split(strings.TrimSpace(string(out)), "\n")
		isRunning := make(map[string]bool)
		for _, rn := range runningNames {
			if rn != "" {
				isRunning[rn] = true
			}
		}

		for _, cfgModel := range m.config.Models {
			if isRunning[cfgModel.ContainerName] {
				name := cfgModel.Name
				container := &Container{
					IsRunning: true,
					LogLines:  []string{fmt.Sprintf("Attached to running container %s...", name)},
				}
				m.containers[name] = container
				cmds = append(cmds, m.attachLogs(cfgModel, container))
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.logsViewport.Width = m.width - sidebarWidth - 1
		m.logsViewport.Height = m.height - 1 // Leave room for status bar
		// Adjust viewport content to fit new width
		m.updateViewport()

	case tea.KeyMsg:
		switch m.state {
		case MainView:
			switch msg.String() {
			case "q", "ctrl+c":
				m.killAllContainers()
				return m, tea.Quit
			case "esc":
				m.state = MenuView
				m.selectedMenuIdx = 0
			case "up", "k":
				if m.selectedModel > 0 {
					m.selectedModel--
					m.updateViewport()
				}
			case "down", "j":
				if m.selectedModel < len(m.config.Models)-1 {
					m.selectedModel++
					m.updateViewport()
				}
			case "enter":
				cfgModel := m.config.Models[m.selectedModel]
				if c, exists := m.containers[cfgModel.Name]; exists && c.IsRunning {
					m.state = ConfirmView
					m.pendingStop = &cfgModel
				} else {
					cmds = append(cmds, m.toggleContainer(cfgModel))
				}
			default:
				m.logsViewport, cmd = m.logsViewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		case MenuView:
			switch msg.String() {
			case "esc":
				m.state = MainView
			case "q", "ctrl+c":
				m.killAllContainers()
				return m, tea.Quit
			case "up", "k":
				if m.selectedMenuIdx > 0 {
					m.selectedMenuIdx--
				}
			case "down", "j":
				if m.selectedMenuIdx < len(m.menuOptions)-1 {
					m.selectedMenuIdx++
				}
			case "enter":
				if m.menuOptions[m.selectedMenuIdx] == "Exit" {
					m.killAllContainers()
					return m, tea.Quit
				} else if m.menuOptions[m.selectedMenuIdx] == "Settings" {
					m.state = SettingsView
					// find current theme index
					for i, t := range Themes {
						if t.Name == m.theme.Name {
							m.selectedThemeIdx = i
							break
						}
					}
				}
			}

		case SettingsView:
			switch msg.String() {
			case "esc":
				m.state = MenuView
			case "q", "ctrl+c":
				m.killAllContainers()
				return m, tea.Quit
			case "up", "k":
				if m.selectedThemeIdx > 0 {
					m.selectedThemeIdx--
					m.applyTheme(Themes[m.selectedThemeIdx])
				}
			case "down", "j":
				if m.selectedThemeIdx < len(Themes)-1 {
					m.selectedThemeIdx++
					m.applyTheme(Themes[m.selectedThemeIdx])
				}
			case "enter":
				m.state = MainView
			}
			
		case AlertView:
			switch msg.String() {
			case "esc", "enter", "q", "ctrl+c":
				m.state = MainView
			}
			
		case ConfirmView:
			switch msg.String() {
			case "esc", "n", "N", "q", "ctrl+c":
				m.state = MainView
				m.pendingStop = nil
			case "enter", "y", "Y":
				if m.pendingStop != nil {
					cmds = append(cmds, m.stopContainer(*m.pendingStop))
				}
				m.state = MainView
				m.pendingStop = nil
			}
		}

	case tickStatsMsg:
		cmds = append(cmds, fetchStatsCmd, tickStatsCmd())

	case statsMsg:
		m.stats = MachineStats(msg)

	// Mouse scrolling for viewport
	case tea.MouseMsg:
		if m.state == MainView {
			m.logsViewport, cmd = m.logsViewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case logMsg:
		if container, exists := m.containers[msg.modelName]; exists {
			container.LogLines = append(container.LogLines, msg.line)
			// Truncate logs if they get too large (e.g. 5000 lines max)
			if len(container.LogLines) > 5000 {
				container.LogLines = container.LogLines[1000:]
			}
			
			// Only update viewport if this is the currently selected model
			if m.config.Models[m.selectedModel].Name == msg.modelName {
				atBottom := m.logsViewport.AtBottom()
				m.updateViewport()
				if atBottom {
					m.logsViewport.GotoBottom()
				}
			}
		}

	case containerExitMsg:
		lastLog := ""
		if c, exists := m.containers[msg.modelName]; exists {
			c.IsRunning = false
			if len(c.LogLines) > 0 {
				lastLog = c.LogLines[len(c.LogLines)-1]
			}
		}
		
		if msg.err != nil && !strings.Contains(msg.err.Error(), "signal: killed") {
			m.alertMsg = fmt.Sprintf("Model '%s' failed.\n\nError: %v\nDetail: %s", msg.modelName, msg.err, lastLog)
			m.state = AlertView
		}
		
		if m.config.Models[m.selectedModel].Name == msg.modelName {
			m.updateViewport()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) applyTheme(t Theme) {
	m.theme = t
	m.styles = MakeStyles(t)
	m.updateViewport() // Refresh styles on text
}

func (m *model) updateViewport() {
	if len(m.config.Models) == 0 {
		m.logsViewport.SetContent("No models configured.")
		return
	}
	modelName := m.config.Models[m.selectedModel].Name
	if container, exists := m.containers[modelName]; exists {
		var styled []string
		for _, line := range container.LogLines {
			styled = append(styled, m.styles.LogText.Render(line))
		}
		m.logsViewport.SetContent(strings.Join(styled, "\n"))
	} else {
		m.logsViewport.SetContent(m.styles.LogText.Render("Not running. Press [Enter] to start."))
	}
}

// stopContainer properly terminates the container process and daemon instance
func (m *model) stopContainer(cfgModel config.Model) tea.Cmd {
	name := cfgModel.Name

	if container, exists := m.containers[name]; exists && container.IsRunning {
		// Run docker rm -f in background to strictly enforce docker stop
		go exec.Command("docker", "rm", "-f", cfgModel.ContainerName).Run()

		if container.Cmd != nil && container.Cmd.Process != nil {
			container.Cmd.Process.Kill()
		}
	}
	return nil
}

// attachLogs attaches to a running container's logging stream
func (m *model) attachLogs(cfgModel config.Model, container *Container) tea.Cmd {
	name := cfgModel.Name
	return func() tea.Msg {
		logsCmd := BuildDockerLogsCmd(cfgModel.ContainerName)
		container.Cmd = logsCmd

		stdout, err := logsCmd.StdoutPipe()
		if err != nil {
			return containerExitMsg{modelName: name, err: err}
		}
		stderr, err := logsCmd.StderrPipe()
		if err != nil {
			return containerExitMsg{modelName: name, err: err}
		}

		if err := logsCmd.Start(); err != nil {
			return containerExitMsg{modelName: name, err: err}
		}

		var wg sync.WaitGroup
		wg.Add(2)

		readPipe := func(reader io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				line := scanner.Text()
				if program != nil {
					program.Send(logMsg{modelName: name, line: line})
				}
			}
		}
		go readPipe(stdout)
		go readPipe(stderr)

		err = logsCmd.Wait()
		wg.Wait()
		return containerExitMsg{modelName: name, err: err}
	}
}

// toggleContainer starts or stops a container
func (m *model) toggleContainer(configModel config.Model) tea.Cmd {
	name := configModel.Name

	// Stop it if it exists and is running (Fallback defensive check)
	if container, exists := m.containers[name]; exists && container.IsRunning {
		return m.stopContainer(configModel)
	}

	// Otherwise, start it
	cmd, err := BuildDockerCmd(m.config, &configModel)
	if err != nil {
		m.containers[name] = &Container{
			IsRunning: false,
			LogLines: []string{fmt.Sprintf("Failed to build command: %v", err)},
		}
		m.updateViewport()
		return nil
	}

	container := &Container{
		IsRunning: true,
		LogLines:  []string{fmt.Sprintf("Starting container for %s...", name)},
	}
	m.containers[name] = container
	m.updateViewport()

	return func() tea.Msg {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return containerExitMsg{modelName: name, err: fmt.Errorf("%v: %s", err, string(out))}
		}

		// Fire attachLogs command separately so it can dispatch to event loop immediately!
		// Wait, returning the function closure directly simulates returning the actual tea.Msg from attachLogs!
		attachFn := m.attachLogs(configModel, container)
		return attachFn()
	}
}

func (m *model) killAllContainers() {
	for _, cfgModel := range m.config.Models {
		if c, exists := m.containers[cfgModel.Name]; exists && c.IsRunning {
			if c.Cmd != nil && c.Cmd.Process != nil {
				c.Cmd.Process.Kill() // Sever the observer ONLY
			}
		}
	}
}

// Global reference for async log streaming
var program *tea.Program

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	contentHeight := m.height - 1 // Status bar
	mainInnerWidth := m.width - sidebarWidth - 1

	// Render Sidebar
	var sidebarLines []string
	sidebarLines = append(sidebarLines, m.styles.SidebarTitle.Render(" Models"))
	
	for i, cfgModel := range m.config.Models {
		name := cfgModel.Name
		
		status := "[ ]"
		if c, exists := m.containers[name]; exists && c.IsRunning {
			status = "[●]"
		}

		display := fmt.Sprintf(" %s %s", status, name)
		if len(display) > sidebarWidth-4 {
			display = display[:sidebarWidth-7] + "..."
		}

		if i == m.selectedModel {
			sidebarLines = append(sidebarLines, m.styles.SelectedItem.Render(display))
		} else {
			styled := m.styles.Item.Render(display)
			if status == "[●]" {
				// highlight running items even if unselected
				styled = lipgloss.NewStyle().Foreground(m.theme.Success).Padding(0, 1).Render(display)
			}
			sidebarLines = append(sidebarLines, styled)
		}
	}

	sidebarLines = append(sidebarLines, "")
	sidebarLines = append(sidebarLines, m.styles.Item.Render(" Keys: [Enter] Toggle"))
	sidebarLines = append(sidebarLines, m.styles.Item.Render("       [Esc] Options"))
	sidebarLines = append(sidebarLines, m.styles.Item.Render("       [Q] Quit"))

	// Build rows
	var rows []string
	divider := m.styles.Divider.Render("│")

	viewContent := strings.Split(m.logsViewport.View(), "\n")

	for i := 0; i < contentHeight; i++ {
		var left string
		if i < len(sidebarLines) {
			left = sidebarLines[i]
		}
		left = lipgloss.NewStyle().Background(m.theme.Background).Width(sidebarWidth).Render(left)

		var right string
		if i < len(viewContent) {
			right = viewContent[i]
		}
		right = lipgloss.NewStyle().Background(m.theme.Background).Width(mainInnerWidth).Render(right)

		rows = append(rows, left+divider+right)
	}

	mainBody := strings.Join(rows, "\n")

	// Status bar
	statusText := fmt.Sprintf(" Model: %s  |  CPU: %s  |  RAM: %s  |  GPU: %s  |  VRAM: %s",
		m.config.Models[m.selectedModel].Name,
		m.stats.CPU,
		m.stats.RAM,
		m.stats.GPU,
		m.stats.VRAM,
	)
	statusBar := m.styles.Status.Width(m.width).Render(statusText)

	ui := mainBody + "\n" + statusBar

	// Render Overlays
	if m.state == MenuView || m.state == SettingsView || m.state == AlertView || m.state == ConfirmView {
		ui = m.renderOverlay(ui)
	}

	return ui
}

func (m model) renderOverlay(background string) string {
	var overlay string

	if m.state == MenuView {
		var lines []string
		lines = append(lines, m.styles.OverlayTitle.Render("Menu"))
		for i, opt := range m.menuOptions {
			if i == m.selectedMenuIdx {
				lines = append(lines, m.styles.OverlaySelItem.Render("▸ "+opt))
			} else {
				lines = append(lines, m.styles.OverlayItem.Render("  "+opt))
			}
		}
		overlay = m.styles.OverlayBox.Render(strings.Join(lines, "\n"))
	} else if m.state == SettingsView {
		var lines []string
		lines = append(lines, m.styles.OverlayTitle.Render("Select Theme"))
		for i, t := range Themes {
			if i == m.selectedThemeIdx {
				lines = append(lines, m.styles.OverlaySelItem.Render("▸ "+t.Name))
			} else {
				lines = append(lines, m.styles.OverlayItem.Render("  "+t.Name))
			}
		}
		overlay = m.styles.OverlayBox.Render(strings.Join(lines, "\n"))
	} else if m.state == AlertView {
		var lines []string
		lines = append(lines, m.styles.OverlayTitle.Foreground(m.theme.Error).Render("Error / Alert"))
		lines = append(lines, m.styles.OverlayItem.Render(m.alertMsg))
		lines = append(lines, "")
		lines = append(lines, m.styles.OverlaySelItem.Render("  OK  "))
		overlay = m.styles.OverlayBox.BorderForeground(m.theme.Error).Render(strings.Join(lines, "\n"))
	} else if m.state == ConfirmView {
		var lines []string
		lines = append(lines, m.styles.OverlayTitle.Foreground(m.theme.Header).Render("Confirm Stop"))
		if m.pendingStop != nil {
			lines = append(lines, m.styles.OverlayItem.Render(fmt.Sprintf("Stop container '%s'?", m.pendingStop.Name)))
		}
		lines = append(lines, "")
		lines = append(lines, m.styles.OverlaySelItem.Render(" Yes (Enter/Y) ")+"  "+m.styles.OverlayItem.Render("No (Esc/N)"))
		overlay = m.styles.OverlayBox.Render(strings.Join(lines, "\n"))
	}

	// Place overlay in the center and use theme background for blank space
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay, lipgloss.WithWhitespaceChars(" "), lipgloss.WithWhitespaceBackground(m.theme.Background))
}
