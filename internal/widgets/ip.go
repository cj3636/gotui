package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// IPWidget displays IP information
type IPWidget struct {
	BaseWidget
	ipInfo         IPInfo
	err            error
	lastUpdate     time.Time
	updateInterval time.Duration
}

// IPInfo contains IP address information
type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
}

// IPMsg contains IP information
type IPMsg struct {
	info IPInfo
	err  error
}

// NewIPWidget creates a new IP information widget
func NewIPWidget(refreshInterval int) *IPWidget {
	return &IPWidget{
		BaseWidget:     NewBaseWidget("🌐 IP Information"),
		updateInterval: time.Duration(refreshInterval) * time.Second,
	}
}

// Init initializes the widget
func (w *IPWidget) Init() tea.Cmd {
	return w.fetchIPInfo()
}

// Update handles messages
func (w *IPWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case IPMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.ipInfo = msg.info
			w.err = nil
		}
		w.lastUpdate = time.Now()
		return w, tea.Tick(w.updateInterval, func(t time.Time) tea.Msg {
			cmd := w.fetchIPInfo()
			return cmd()
		})
	}
	return w, nil
}

// View renders the widget
func (w *IPWidget) View() string {
	var content string
	if w.err != nil {
		content = fmt.Sprintf("Error: %v", w.err)
	} else if w.ipInfo.IP == "" {
		content = "Loading IP info..."
	} else {
		content = fmt.Sprintf(
			"IP:       %s\n"+
				"Location: %s, %s\n"+
				"Country:  %s\n"+
				"Timezone: %s\n"+
				"ISP:      %s",
			w.ipInfo.IP,
			w.ipInfo.City,
			w.ipInfo.Region,
			w.ipInfo.Country,
			w.ipInfo.Timezone,
			w.ipInfo.Org,
		)
	}
	return w.RenderContent(content)
}

func (w *IPWidget) fetchIPInfo() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get("https://ipinfo.io/json")
		if err != nil {
			return IPMsg{err: err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return IPMsg{err: err}
		}

		var info IPInfo
		if err := json.Unmarshal(body, &info); err != nil {
			return IPMsg{err: err}
		}

		return IPMsg{info: info}
	}
}
