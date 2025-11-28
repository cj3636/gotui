package widgets

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// WeatherWidget displays weather information from wttr.in
type WeatherWidget struct {
	BaseWidget
	location       string
	weatherData    string
	err            error
	lastUpdate     time.Time
	updateInterval time.Duration
}

// WeatherMsg contains weather data
type WeatherMsg struct {
	data string
	err  error
}

// WeatherRefreshMsg signals it's time to refresh the weather
type WeatherRefreshMsg struct{}

// NewWeatherWidget creates a new weather widget
func NewWeatherWidget(location string, refreshInterval int) *WeatherWidget {
	return &WeatherWidget{
		BaseWidget:     NewBaseWidget("🌤️  Weather"),
		location:       location,
		weatherData:    "Loading...",
		updateInterval: time.Duration(refreshInterval) * time.Second,
	}
}

// Init initializes the widget
func (w *WeatherWidget) Init() tea.Cmd {
	return w.fetchWeather()
}

// Update handles messages
func (w *WeatherWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case WeatherMsg:
		if msg.err != nil {
			w.err = msg.err
			w.weatherData = fmt.Sprintf("Error: %v", msg.err)
		} else {
			w.weatherData = msg.data
			w.err = nil
		}
		w.lastUpdate = time.Now()
		return w, tea.Tick(w.updateInterval, func(t time.Time) tea.Msg {
			return WeatherRefreshMsg{}
		})
	case WeatherRefreshMsg:
		return w, w.fetchWeather()
	}
	return w, nil
}

// View renders the widget
func (w *WeatherWidget) View() string {
	content := w.weatherData
	if w.lastUpdate.IsZero() {
		content = "Loading weather data..."
	} else {
		elapsed := time.Since(w.lastUpdate)
		if elapsed < time.Minute {
			content += fmt.Sprintf("\n\nUpdated: %ds ago", int(elapsed.Seconds()))
		} else {
			content += fmt.Sprintf("\n\nUpdated: %dm ago", int(elapsed.Minutes()))
		}
	}
	return w.RenderContent(content)
}

func (w *WeatherWidget) fetchWeather() tea.Cmd {
	return func() tea.Msg {
		// Format location for wttr.in (replace spaces with +)
		location := strings.ReplaceAll(w.location, " ", "+")

		// Request weather data in plain text format with custom formatting
		url := fmt.Sprintf("https://wttr.in/%s?format=%%l:+%%C+%%t\n%%w\n%%p\n%%h", location)

		resp, err := http.Get(url)
		if err != nil {
			return WeatherMsg{err: err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return WeatherMsg{err: err}
		}

		data := string(body)

		// Clean up the data
		lines := strings.Split(data, "\n")
		var cleaned []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleaned = append(cleaned, line)
			}
		}

		result := strings.Join(cleaned, "\n")

		return WeatherMsg{data: result}
	}
}
