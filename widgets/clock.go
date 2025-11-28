package widgets

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// clockWidget displays the current time.
type clockWidget struct {
	now time.Time
}

// NewClockWidget returns an initialized clock widget.
func NewClockWidget() Widget {
	return &clockWidget{now: time.Now()}
}

func (c *clockWidget) Title() string { return "Clock" }

func (c *clockWidget) Init() tea.Cmd { return Tick(time.Second) }

func (c *clockWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		c.now = time.Now()
		return c, Tick(time.Second)
	}
	return c, nil
}

func (c *clockWidget) View(width, _ int) string {
	return c.now.Format("Mon Jan _2 15:04:05 MST")
}
