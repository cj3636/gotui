package widgets

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ClockWidget displays the current time
type ClockWidget struct {
	BaseWidget
	currentTime time.Time
}

// NewClockWidget creates a new clock widget
func NewClockWidget() *ClockWidget {
	return &ClockWidget{
		BaseWidget:  NewBaseWidget("🕐 Clock"),
		currentTime: time.Now(),
	}
}

// TickMsg is sent every second to update the clock
type TickMsg time.Time

// Init initializes the widget
func (w *ClockWidget) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages
func (w *ClockWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		w.currentTime = time.Time(msg)
		return w, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
	}
	return w, nil
}

// View renders the widget
func (w *ClockWidget) View() string {
	content := fmt.Sprintf(
		"%s\n\n%s\n%s",
		w.currentTime.Format("Monday"),
		w.currentTime.Format("January 2, 2006"),
		w.currentTime.Format("15:04:05"),
	)
	return w.RenderContent(content)
}
