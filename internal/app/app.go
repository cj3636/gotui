package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cj3636/gotui/internal/config"
	"github.com/cj3636/gotui/internal/widgets"
)

// Model is the main application model
type Model struct {
	config  *config.Config
	widgets []widgets.Widget
	width   int
	height  int
	ready   bool
}

// NewModel creates a new application model
func NewModel(cfg *config.Config) Model {
	// Create widgets based on configuration
	var widgetList []widgets.Widget

	// Add clock and calendar
	widgetList = append(widgetList, widgets.NewClockWidget())
	widgetList = append(widgetList, widgets.NewCalendarWidget())

	// Add weather widget
	if cfg.WeatherLocation != "" {
		widgetList = append(widgetList, widgets.NewWeatherWidget(
			cfg.WeatherLocation,
			cfg.RefreshIntervals.Weather,
		))
	}

	// Add GitHub widget
	if len(cfg.GithubRepos) > 0 {
		widgetList = append(widgetList, widgets.NewGithubWidget(
			cfg.GithubToken,
			cfg.GithubRepos,
			cfg.RefreshIntervals.Github,
		))
	}

	// Add GitLab widget
	if len(cfg.GitlabProjects) > 0 {
		widgetList = append(widgetList, widgets.NewGitlabWidget(
			cfg.GitlabToken,
			cfg.GitlabProjects,
			cfg.RefreshIntervals.Gitlab,
		))
	}

	// Add system resources widget
	widgetList = append(widgetList, widgets.NewSystemWidget(cfg.RefreshIntervals.System))

	// Add IP information widget
	widgetList = append(widgetList, widgets.NewIPWidget(cfg.RefreshIntervals.IP))

	// Add SMART status widget
	widgetList = append(widgetList, widgets.NewSMARTWidget())

	// Add text viewer widget
	if cfg.TextFile != "" {
		widgetList = append(widgetList, widgets.NewTextViewerWidget(cfg.TextFile))
	}

	// Add markdown viewer widget
	if cfg.MarkdownFile != "" {
		widgetList = append(widgetList, widgets.NewMarkdownWidget(cfg.MarkdownFile))
	}

	return Model{
		config:  cfg,
		widgets: widgetList,
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, widget := range m.widgets {
		cmds = append(cmds, widget.Init())
	}
	return tea.Batch(cmds...)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateWidgetSizes()
		return m, nil
	}

	// Update all widgets
	var cmds []tea.Cmd
	for i, widget := range m.widgets {
		updatedWidget, cmd := widget.Update(msg)
		m.widgets[i] = updatedWidget
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate grid layout
	rows := m.config.Layout.Rows
	cols := m.config.Layout.Cols

	// Calculate widget dimensions
	widgetWidth := m.width / cols
	widgetHeight := m.height / rows

	// Build grid
	var gridRows []string
	widgetIndex := 0

	for row := 0; row < rows; row++ {
		var rowWidgets []string
		for col := 0; col < cols; col++ {
			if widgetIndex < len(m.widgets) {
				widget := m.widgets[widgetIndex]
				widget.SetSize(widgetWidth, widgetHeight)
				rowWidgets = append(rowWidgets, widget.View())
				widgetIndex++
			} else {
				// Empty placeholder
				emptyStyle := lipgloss.NewStyle().
					Width(widgetWidth).
					Height(widgetHeight)
				rowWidgets = append(rowWidgets, emptyStyle.Render(""))
			}
		}
		gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, rowWidgets...))
	}

	view := lipgloss.JoinVertical(lipgloss.Left, gridRows...)

	// Add help text at the bottom
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)
	
	help := helpStyle.Render("Press 'q', 'Esc', or 'Ctrl+C' to quit")
	
	return view + "\n" + help
}

// updateWidgetSizes updates the size of all widgets
func (m *Model) updateWidgetSizes() {
	rows := m.config.Layout.Rows
	cols := m.config.Layout.Cols

	widgetWidth := m.width / cols
	widgetHeight := m.height / rows

	for _, widget := range m.widgets {
		widget.SetSize(widgetWidth, widgetHeight)
	}
}

// LoadConfig loads the application configuration
func LoadConfig() (*config.Config, error) {
	configPath, err := config.FindConfigFile()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		// Return default config if file doesn't exist
		return &config.Config{
			WeatherLocation: "New York",
			RefreshIntervals: config.RefreshIntervals{
				Weather: 1800,
				Github:  300,
				Gitlab:  300,
				System:  5,
				IP:      3600,
			},
			GithubRepos: []string{},
			Layout: config.Layout{
				Rows: 3,
				Cols: 3,
			},
		}, nil
	}

	return cfg, nil
}
