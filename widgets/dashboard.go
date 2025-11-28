package widgets

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Widget defines the rendering and update contract for dashboard components.
type Widget interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Widget, tea.Cmd)
	View(width, height int) string
	Title() string
}

// Dashboard is the root Bubble Tea model that arranges multiple widgets.
type Dashboard struct {
	widgets []Widget
	width   int
	height  int
}

// NewDashboard bootstraps the dashboard with the default widget set.
func NewDashboard() Dashboard {
	ws := []Widget{
		NewClockWidget(),
		NewWeatherWidget(),
		NewMoonWidget(),
		NewSystemWidget(),
		NewIPWidget(),
		NewMarkdownWidget(),
		NewGitHubWidget(),
		NewGitLabWidget(),
	}
	return Dashboard{widgets: ws}
}

// Init starts all widget initialization commands.
func (d Dashboard) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(d.widgets))
	for i, w := range d.widgets {
		cmds[i] = w.Init()
	}
	return tea.Batch(cmds...)
}

// Update dispatches messages to all widgets and handles window resizing.
func (d Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = m.Width
		d.height = m.Height
	}

	for i, w := range d.widgets {
		updated, cmd := w.Update(msg)
		d.widgets[i] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return d, tea.Batch(cmds...)
}

// View renders widgets in one or two columns depending on terminal width.
func (d Dashboard) View() string {
	if d.width == 0 {
		// Layout isn't ready until we receive the first WindowSizeMsg.
		return "Loading layout..."
	}

	useColumns := d.width >= 90
	columnWidth := d.width - 6 // account for padding and borders
	if useColumns {
		columnWidth = (d.width / 2) - 5
	}

	boxes := make([]string, len(d.widgets))
	for i, w := range d.widgets {
		boxes[i] = renderPanel(w.Title(), w.View(columnWidth-4, d.height/len(d.widgets)))
	}

	if !useColumns {
		return strings.Join(boxes, "\n")
	}

	// Two-column layout
	var rows []string
	for i := 0; i < len(boxes); i += 2 {
		left := boxes[i]
		right := ""
		if i+1 < len(boxes) {
			right = boxes[i+1]
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	panelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

func renderPanel(title, body string) string {
	header := titleStyle.Render(title)
	return panelStyle.Render(header + "\n" + body)
}

// TickMsg is used for widgets that need periodic updates.
type TickMsg time.Time

// Tick returns a Tea command that emits a TickMsg after the provided duration.
func Tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return TickMsg(t) })
}
