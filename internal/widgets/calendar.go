package widgets

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// CalendarWidget displays a monthly calendar
type CalendarWidget struct {
	BaseWidget
	currentTime time.Time
}

// NewCalendarWidget creates a new calendar widget
func NewCalendarWidget() *CalendarWidget {
	return &CalendarWidget{
		BaseWidget:  NewBaseWidget("📅 Calendar"),
		currentTime: time.Now(),
	}
}

// Init initializes the widget
func (w *CalendarWidget) Init() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages
func (w *CalendarWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		w.currentTime = time.Time(msg)
		return w, tea.Tick(time.Minute, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
	}
	return w, nil
}

// View renders the widget
func (w *CalendarWidget) View() string {
	content := w.generateCalendar()
	return w.RenderContent(content)
}

func (w *CalendarWidget) generateCalendar() string {
	now := w.currentTime
	year, month, _ := now.Date()

	// Get the first day of the month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())

	// Get the last day of the month
	lastDay := firstDay.AddDate(0, 1, -1)

	var sb strings.Builder

	// Month and year header
	sb.WriteString(fmt.Sprintf("%s %d\n\n", month, year))

	// Day headers
	sb.WriteString("Su Mo Tu We Th Fr Sa\n")

	// Add padding for the first week
	startWeekday := int(firstDay.Weekday())
	for i := 0; i < startWeekday; i++ {
		sb.WriteString("   ")
	}

	// Add days
	today := now.Day()
	for day := 1; day <= lastDay.Day(); day++ {
		if day == today {
			sb.WriteString(fmt.Sprintf("[%2d]", day))
		} else {
			sb.WriteString(fmt.Sprintf("%2d ", day))
		}

		weekday := (startWeekday + day - 1) % 7
		if weekday == 6 && day != lastDay.Day() {
			sb.WriteString("\n")
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}
