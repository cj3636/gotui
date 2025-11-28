package widgets

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultWeatherLocation = "San Francisco"
	defaultWeatherParams   = "format=3"
	defaultMoonLocation    = "moon"
	defaultMoonParams      = "format=Moon:+%m"
)

// wttrWidget pulls conditions from wttr.in endpoints.
type wttrWidget struct {
	title    string
	location string
	params   string
	summary  string
	err      error
}

// NewWeatherWidget constructs the weather widget using the WTTR_LOCATION and
// WTTR_PARAMS environment variables when provided.
func NewWeatherWidget() Widget {
	loc := os.Getenv("WTTR_LOCATION")
	if loc == "" {
		loc = defaultWeatherLocation
	}
	params := os.Getenv("WTTR_PARAMS")
	if params == "" {
		params = defaultWeatherParams
	}
	return &wttrWidget{title: "Weather", location: loc, params: params}
}

// NewMoonWidget constructs a moon phase widget that points to the wttr.in/moon
// endpoint. Location and params can be overridden with WTTR_MOON_LOCATION and
// WTTR_MOON_PARAMS respectively.
func NewMoonWidget() Widget {
	loc := os.Getenv("WTTR_MOON_LOCATION")
	if loc == "" {
		loc = defaultMoonLocation
	}
	params := os.Getenv("WTTR_MOON_PARAMS")
	if params == "" {
		params = defaultMoonParams
	}
	return &wttrWidget{title: "Moon Phase", location: loc, params: params}
}

func (w *wttrWidget) Title() string { return w.title }

func (w *wttrWidget) Init() tea.Cmd { return w.fetch() }

func (w *wttrWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch m := msg.(type) {
	case TickMsg:
		return w, w.fetch()
	case weatherMsg:
		w.summary = m.summary
		w.err = m.err
		return w, Tick(30 * time.Minute)
	}
	return w, nil
}

func (w *wttrWidget) View(width, _ int) string {
	if w.err != nil {
		return fmt.Sprintf("Location: %s\nError: %v", w.location, w.err)
	}
	if w.summary == "" {
		return fmt.Sprintf("Location: %s\nLoading forecast...", w.location)
	}
	return fmt.Sprintf("Location: %s\n%s", w.location, w.summary)
}

type weatherMsg struct {
	summary string
	err     error
}

func (w *wttrWidget) fetch() tea.Cmd {
	loc := w.location
	return func() tea.Msg {
		client := &http.Client{Timeout: 5 * time.Second}
		path := fmt.Sprintf("https://wttr.in/%s", url.PathEscape(loc))
		if w.params != "" {
			path = path + "?" + w.params
		}
		resp, err := client.Get(path)
		if err != nil {
			return weatherMsg{err: err}
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return weatherMsg{err: err}
		}
		return weatherMsg{summary: string(body)}
	}
}
