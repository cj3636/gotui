# GoTUI - Terminal Dashboard

A powerful, modular terminal user interface (TUI) dashboard application built with Go and the [Charmbracelet](https://charm.sh/) toolkit. Inspired by [wtfutil](https://github.com/wtfutil/wtf), GoTUI provides real-time monitoring and information display in a beautiful terminal interface.

## Features

### 🎨 Widgets

- **🕐 Clock**: Real-time clock display with date and time
- **📅 Calendar**: Monthly calendar view with current day highlighting
- **🌤️ Weather**: Live weather information from wttr.in
- **🐙 GitHub**: Monitor repositories with stars, forks, issues, and PR counts (supports private repos with API token)
- **🦊 GitLab**: Track projects with statistics (supports private projects with API token)
- **💻 System Resources**: Real-time CPU, memory, and disk usage monitoring
- **🌐 IP Information**: Display your public IP address and location details
- **💾 SMART Status**: Disk drive status information
- **📄 Text Viewer**: Display content from any text file
- **📝 Markdown Viewer**: Render markdown files with Glow-powered rendering

### ✨ Key Features

- **Automatic Resizing**: All widgets automatically adjust to terminal window size
- **Grid Layout**: Configurable rows and columns for widget placement
- **Refresh Intervals**: Customizable refresh rates for each widget type
- **API Support**: Full support for GitHub and GitLab private repositories using API tokens
- **Modular Design**: Each widget is an independent component that can be easily extended

## Installation

### Prerequisites

- Go 1.21 or higher
- A terminal with support for 256 colors

### Build from Source

```bash
git clone https://github.com/cj3636/gotui.git
cd gotui
go build -o gotui .
```

### Run

```bash
./gotui
```

## Configuration

GoTUI uses a YAML configuration file. By default, it looks for `config.yaml` in:
1. Current directory
2. `~/.config/gotui/config.yaml`

### Example Configuration

```yaml
# GitHub API Token (optional - for private repositories)
github_token: ""

# GitLab API Token (optional - for private repositories)
gitlab_token: ""

# Weather location (city name or coordinates)
weather_location: "New York"

# Refresh intervals (in seconds)
refresh_intervals:
  weather: 1800      # 30 minutes
  github: 300        # 5 minutes
  gitlab: 300        # 5 minutes
  system: 5          # 5 seconds
  ip: 3600           # 1 hour

# GitHub repositories to monitor (format: "owner/repo")
github_repos:
  - "charmbracelet/bubbletea"
  - "charmbracelet/lipgloss"

# GitLab projects to monitor (format: "namespace/project")
gitlab_projects: []

# Text file to display
text_file: ""

# Markdown file to display
markdown_file: "example.md"

# Widget layout configuration
layout:
  rows: 3
  cols: 3
```

### Configuration Options

#### API Tokens

- **github_token**: GitHub Personal Access Token for accessing private repositories
  - Create at: https://github.com/settings/tokens
  - Required scopes: `repo` (for private repos)
  
- **gitlab_token**: GitLab Personal Access Token for accessing private projects
  - Create at: https://gitlab.com/-/profile/personal_access_tokens
  - Required scopes: `read_api`, `read_repository`

#### Refresh Intervals

Control how often each widget updates its data:
- `weather`: Weather data refresh interval
- `github`: GitHub repository data refresh interval
- `gitlab`: GitLab project data refresh interval
- `system`: System resources refresh interval
- `ip`: IP information refresh interval

#### Layout

- `rows`: Number of rows in the widget grid (default: 3)
- `cols`: Number of columns in the widget grid (default: 3)

The total number of widgets displayed is `rows × cols`. Widgets are placed left-to-right, top-to-bottom.

## Usage

### Keyboard Controls

- `q`, `Esc`, or `Ctrl+C`: Quit the application

### Widget Display

The application automatically displays widgets in a grid layout based on your configuration. Widgets are shown in the following order:

1. Clock
2. Calendar
3. Weather (if location configured)
4. GitHub (if repositories configured)
5. GitLab (if projects configured)
6. System Resources
7. IP Information
8. SMART Status
9. Text Viewer (if file configured)
10. Markdown Viewer (if file configured)

## Architecture

### Project Structure

```
gotui/
├── main.go                 # Application entry point
├── config.yaml             # Configuration file
├── internal/
│   ├── app/
│   │   └── app.go         # Main application model
│   ├── config/
│   │   └── config.go      # Configuration management
│   └── widgets/
│       ├── widget.go      # Base widget interface
│       ├── clock.go       # Clock widget
│       ├── calendar.go    # Calendar widget
│       ├── weather.go     # Weather widget
│       ├── github.go      # GitHub widget
│       ├── gitlab.go      # GitLab widget
│       ├── system.go      # System resources widget
│       ├── ip.go          # IP information widget
│       ├── smart.go       # SMART status widget
│       ├── textviewer.go  # Text viewer widget
│       └── markdown.go    # Markdown viewer widget
└── go.mod
```

### Widget System

All widgets implement the `Widget` interface:

```go
type Widget interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Widget, tea.Cmd)
    View() string
    SetSize(width, height int)
    Title() string
}
```

This modular design makes it easy to:
- Add new widgets
- Customize existing widgets
- Replace widgets with custom implementations

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [gopsutil](https://github.com/shirou/gopsutil) - System information
- [go-github](https://github.com/google/go-github) - GitHub API client
- [go-gitlab](https://github.com/xanzy/go-gitlab) - GitLab API client

## Contributing

Contributions are welcome! Here are some ways you can contribute:

- Add new widgets
- Improve existing widgets
- Fix bugs
- Improve documentation
- Add tests

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Inspired by [wtfutil](https://github.com/wtfutil/wtf)
- Built with the amazing [Charmbracelet](https://charm.sh/) toolkit
- Weather data from [wttr.in](https://wttr.in)
- IP information from [ipinfo.io](https://ipinfo.io)
