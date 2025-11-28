package widgets

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Widget is the base interface for all widgets
type Widget interface {
	// Init initializes the widget
	Init() tea.Cmd

	// Update handles messages and returns the updated widget
	Update(msg tea.Msg) (Widget, tea.Cmd)

	// View renders the widget
	View() string

	// SetSize sets the widget dimensions
	SetSize(width, height int)

	// Title returns the widget title
	Title() string
}

// BaseWidget provides common functionality for all widgets
type BaseWidget struct {
	width  int
	height int
	title  string
	style  lipgloss.Style
}

// NewBaseWidget creates a new base widget
func NewBaseWidget(title string) BaseWidget {
	return BaseWidget{
		title: title,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1),
	}
}

// SetSize sets the widget dimensions
func (w *BaseWidget) SetSize(width, height int) {
	w.width = width
	w.height = height
}

// Title returns the widget title
func (w *BaseWidget) Title() string {
	return w.title
}

// RenderContent renders content with the widget's style and dimensions
func (w *BaseWidget) RenderContent(content string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Align(lipgloss.Center)

	title := titleStyle.Render(w.title)

	// Calculate available space
	availableWidth := w.width - w.style.GetHorizontalFrameSize()
	availableHeight := w.height - w.style.GetVerticalFrameSize() - 1 // -1 for title

	if availableWidth < 1 {
		availableWidth = 1
	}
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Truncate content to fit
	contentStyle := lipgloss.NewStyle().
		Width(availableWidth).
		Height(availableHeight)

	renderedContent := contentStyle.Render(content)

	// Combine title and content
	combined := lipgloss.JoinVertical(lipgloss.Left, title, renderedContent)

	return w.style.
		Width(w.width - w.style.GetHorizontalFrameSize()).
		Height(w.height - w.style.GetVerticalFrameSize()).
		Render(combined)
}
