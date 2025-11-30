package widgets

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultWeatherLocation = "89701"
	defaultWeatherUnits    = "u"
	defaultWeatherView     = "Fq1"
	defaultMoonLocation    = "moon"
	defaultMoonView        = "Fq1"
)

// wttrWidget pulls conditions from wttr.in endpoints.
type wttrWidget struct {
	title    string
	location string
	units    string
	view     string
	summary  string
	err      error
}

// NewWeatherWidget constructs the weather widget using the WTTR_LOCATION and
// WTTR_PARAMS environment variables when provided.
func NewWeatherWidget() Widget {
	cfg := wttrConfig{
		location: readLocation("WTTR_LOCATION", defaultWeatherLocation),
		units:    readUnits("WTTR_UNITS", defaultWeatherUnits),
		view:     readView("WTTR_VIEW", defaultWeatherView),
	}
	return &wttrWidget{title: "Weather", location: cfg.location, units: cfg.units, view: cfg.view}
}

// NewMoonWidget constructs a moon phase widget that points to the wttr.in/moon
// endpoint. Location and params can be overridden with WTTR_MOON_LOCATION and
// WTTR_MOON_PARAMS respectively.
func NewMoonWidget() Widget {
	cfg := wttrConfig{
		location: readLocation("WTTR_MOON_LOCATION", defaultMoonLocation),
		units:    readUnits("WTTR_MOON_UNITS", defaultWeatherUnits),
		view:     readView("WTTR_MOON_VIEW", defaultMoonView),
	}
	return &wttrWidget{title: "Moon Phase", location: cfg.location, units: cfg.units, view: cfg.view}
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
	loc := sanitizeLocation(w.location)
	params := strings.TrimSpace(w.units + w.view)

	return func() tea.Msg {
		client := &http.Client{Timeout: 5 * time.Second}
		path := fmt.Sprintf("https://wttr.in/%s", loc)
		if params != "" {
			path = path + "?" + params
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

type wttrConfig struct {
	location string
	units    string
	view     string
}

func readLocation(envKey, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(envKey)); v != "" {
		return v
	}
	return fallback
}

func sanitizeLocation(loc string) string {
	loc = strings.TrimSpace(loc)
	loc = strings.ReplaceAll(loc, " ", "_")
	if loc == "" {
		return "_"
	}
	return loc
}

func readUnits(envKey, fallback string) string {
	val := strings.TrimSpace(os.Getenv(envKey))
	if val == "" {
		val = fallback
	}

	lower := strings.ToLower(val)
	base := ""
	if strings.Contains(lower, "m") {
		base = "m"
	}
	if strings.Contains(lower, "u") {
		base = "u"
	}
	if base == "" {
		base = string(fallback[0])
	}

	windMetric := ""
	if strings.Contains(strings.ToUpper(val), "M") {
		windMetric = "M"
	} else if len(fallback) > 1 && strings.Contains(strings.ToUpper(fallback), "M") {
		windMetric = "M"
	}

	return base + windMetric
}

func readView(envKey, fallback string) string {
	val := strings.TrimSpace(os.Getenv(envKey))
	if val == "" {
		val = fallback
	}
	allowed := "0123456789AdFnqQ"
	var builder strings.Builder
	for _, ch := range val {
		if strings.ContainsRune(allowed, ch) {
			builder.WriteRune(ch)
		}
	}
	cleaned := builder.String()
	if cleaned == "" {
		cleaned = fallback
	}
	return cleaned
}
