package widgets

import (
	"math"
	"os"
	"strconv"
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

	sizeHints map[string]int
	columns   int
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

	return Dashboard{
		widgets:   ws,
		sizeHints: loadWidgetHeights(ws),
		columns:   defaultColumns,
	}
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
	case tea.KeyMsg:
		if m.Type == tea.KeyCtrlC || m.String() == "q" || m.String() == "Q" {
			return d, tea.Quit
		}
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

	columns := d.columns
	if columns < 1 {
		columns = 1
	}
	if d.width < 90 {
		columns = 1
	}

	columnWidth := calculateColumnWidth(d.width, columns)
	innerWidth := columnWidth - panelStyle.GetHorizontalFrameSize()
	if innerWidth < 10 {
		innerWidth = 10
	}

	// Determine how much vertical space a single row consumes in each column
	colUnits := make([]int, columns)
	for i, w := range d.widgets {
		colIdx := i % columns
		colUnits[colIdx] += d.heightUnitFor(w)
	}
	maxUnits := maxInt(colUnits...)
	if maxUnits == 0 {
		maxUnits = 1
	}
	baseHeight := int(math.Max(5, float64(d.height)/float64(maxUnits)))

	boxes := make([]string, len(d.widgets))
	for i, w := range d.widgets {
		heightUnits := d.heightUnitFor(w)
		widgetHeight := baseHeight * heightUnits
		boxes[i] = renderPanelSized(columnWidth, widgetHeight, w.Title(), w.View(innerWidth, widgetHeight))
	}

	if columns == 1 {
		return strings.Join(boxes, "\n")
	}

	var rows []string
	for i := 0; i < len(boxes); i += columns {
		rowPanels := make([]string, 0, columns)
		for c := 0; c < columns; c++ {
			if i+c < len(boxes) {
				rowPanels = append(rowPanels, boxes[i+c])
			}
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowPanels...))
	}
	return strings.Join(rows, "\n")
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	panelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

func renderPanelSized(width, height int, title, body string) string {
	style := panelStyle.Copy().Width(width)
	if height > 0 {
		style = style.Height(height)
	}

	frameHeight := style.GetVerticalFrameSize()
	innerHeight := height - frameHeight
	if innerHeight < 1 {
		innerHeight = 1
	}

	header := titleStyle.Render(title)
	bodyLines := clampLines(body, innerHeight-1)

	content := header
	if bodyLines != "" {
		content += "\n" + bodyLines
	}

	return style.Render(content)
}

// TickMsg is used for widgets that need periodic updates.
type TickMsg time.Time

// Tick returns a Tea command that emits a TickMsg after the provided duration.
func Tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func calculateColumnWidth(totalWidth, columns int) int {
	if columns <= 1 {
		return totalWidth - 2
	}
	gap := 2
	return (totalWidth - gap*(columns-1)) / columns
}

func clampLines(content string, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:maxLines], "\n")
}

const defaultColumns = 2

func loadWidgetHeights(widgets []Widget) map[string]int {
	sizes := make(map[string]int, len(widgets))
	for _, w := range widgets {
		key := widgetEnvKey(w.Title())
		sizes[w.Title()] = 1
		if env := getenv(key); env != "" {
			if v, err := strconv.Atoi(env); err == nil && v > 0 {
				sizes[w.Title()] = v
			}
		}
	}
	return sizes
}

func widgetEnvKey(title string) string {
	normalized := strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, title)
	return "WIDGET_HEIGHT_" + strings.ToUpper(normalized)
}

func getenv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	return strings.TrimSpace(v)
}

func (d Dashboard) heightUnitFor(w Widget) int {
	if v, ok := d.sizeHints[w.Title()]; ok && v > 0 {
		return v
	}
	return 1
}

func maxInt(values ...int) int {
	max := 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
