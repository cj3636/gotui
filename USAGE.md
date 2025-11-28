# GoTUI Usage Guide

This guide provides detailed information on how to use and customize GoTUI.

## Quick Start

1. **Build the application:**
   ```bash
   make build
   ```

2. **Run with default configuration:**
   ```bash
   ./gotui
   ```

3. **Create your own configuration:**
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your preferences
   ./gotui
   ```

## Configuration

### Configuration File Locations

GoTUI looks for configuration in the following order:
1. `./config.yaml` (current directory)
2. `~/.config/gotui/config.yaml` (user config directory)

### Basic Configuration

```yaml
# Minimal configuration
weather_location: "San Francisco"
layout:
  rows: 2
  cols: 2
```

### Advanced Configuration

```yaml
# Full configuration with all options
github_token: "ghp_your_token_here"
gitlab_token: "glpat-your_token_here"
weather_location: "Tokyo"

refresh_intervals:
  weather: 1800
  github: 300
  gitlab: 300
  system: 5
  ip: 3600

github_repos:
  - "torvalds/linux"
  - "microsoft/vscode"
  - "your-org/private-repo"

gitlab_projects:
  - "gitlab-org/gitlab"

text_file: "/var/log/syslog"
markdown_file: "./README.md"

layout:
  rows: 4
  cols: 3
```

## Widget Guide

### Clock Widget (🕐)

Displays the current time and date.

- **Updates**: Every second
- **Configuration**: None required
- **Always included**: Yes

**Example output:**
```
Thursday
November 28, 2025
14:35:42
```

### Calendar Widget (📅)

Shows a monthly calendar with the current day highlighted.

- **Updates**: Every minute
- **Configuration**: None required
- **Always included**: Yes

**Example output:**
```
November 2025

Su Mo Tu We Th Fr Sa
                1  2
 3  4  5  6  7  8  9
10 11 12 13 14 15 16
17 18 19 20 21 22 23
24 25 26 27 [28] 29 30
```

### Weather Widget (🌤️)

Displays weather information from wttr.in.

- **Updates**: Every 30 minutes (configurable)
- **Configuration**: `weather_location`
- **Data source**: wttr.in

**Configuration:**
```yaml
weather_location: "New York"  # City name, coordinates, or airport code
refresh_intervals:
  weather: 1800  # 30 minutes
```

**Example output:**
```
New York: Clear 72°F
Wind: 5mph
Precipitation: 0mm
Humidity: 45%
```

### GitHub Widget (🐙)

Monitor GitHub repositories with real-time statistics.

- **Updates**: Every 5 minutes (configurable)
- **Configuration**: `github_repos`, optional `github_token`
- **Features**: Stars, forks, issues, pull requests

**Configuration:**
```yaml
github_token: "ghp_your_token_here"  # Optional, for private repos
github_repos:
  - "owner/repository"
  - "another-owner/another-repo"
refresh_intervals:
  github: 300  # 5 minutes
```

**API Token Setup:**
1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes: `repo` (for private repos)
4. Copy the token to your config

**Example output:**
```
charmbracelet/bubbletea
⭐ 24,567  🍴 789  📝 45

your-org/private-repo
⭐ 123  🍴 12  📝 5
```

### GitLab Widget (🦊)

Track GitLab projects with detailed statistics.

- **Updates**: Every 5 minutes (configurable)
- **Configuration**: `gitlab_projects`, optional `gitlab_token`
- **Features**: Stars, forks, issues, merge requests

**Configuration:**
```yaml
gitlab_token: "glpat-your_token_here"  # Optional, for private projects
gitlab_projects:
  - "namespace/project"
  - "another-namespace/another-project"
refresh_intervals:
  gitlab: 300  # 5 minutes
```

**API Token Setup:**
1. Go to https://gitlab.com/-/profile/personal_access_tokens
2. Create a new token
3. Select scopes: `read_api`, `read_repository`
4. Copy the token to your config

**Example output:**
```
gitlab-org/gitlab
⭐ 1,234  🍴 567  📝 89  🔀 12
```

### System Resources Widget (💻)

Real-time system monitoring.

- **Updates**: Every 5 seconds (configurable)
- **Configuration**: Optional refresh interval
- **Monitors**: CPU, RAM, Disk usage

**Configuration:**
```yaml
refresh_intervals:
  system: 5  # 5 seconds
```

**Example output:**
```
CPU:   45.2%
RAM:   67.8% (5.4 GiB / 8.0 GiB)
Disk:  72.1% (144 GiB / 200 GiB)
OS:    linux/amd64
```

### IP Information Widget (🌐)

Display your public IP address and location.

- **Updates**: Every hour (configurable)
- **Configuration**: Optional refresh interval
- **Data source**: ipinfo.io

**Configuration:**
```yaml
refresh_intervals:
  ip: 3600  # 1 hour
```

**Example output:**
```
IP:       203.0.113.42
Location: San Francisco, CA
Country:  United States
Timezone: America/Los_Angeles
ISP:      Example ISP
```

### SMART Status Widget (💾)

Disk drive status and health information.

- **Updates**: On initialization
- **Configuration**: None required
- **Requirements**: smartmontools (optional)

**Example output:**
```
Disk Status:

Filesystem Size Used Avail Use%
/dev/sda1  200G 144G  56G  72%

Detected drives:
  /dev/sda
  /dev/sdb
```

### Text Viewer Widget (📄)

Display content from any text file.

- **Updates**: On initialization
- **Configuration**: `text_file`
- **Features**: Auto-truncation to 100 lines

**Configuration:**
```yaml
text_file: "/path/to/your/file.txt"
```

**Example output:**
```
Line 1: Your text content
Line 2: More content
Line 3: ...
```

### Markdown Viewer Widget (📝)

Render markdown files with beautiful formatting.

- **Updates**: On initialization
- **Configuration**: `markdown_file`
- **Features**: Syntax highlighting, auto-truncation

**Configuration:**
```yaml
markdown_file: "./README.md"
```

**Example output:**
```
  My Project

  This is a README file

  ## Features

  • Feature 1
  • Feature 2
  • Feature 3
```

## Layout Customization

### Grid System

GoTUI uses a grid layout system. Widgets are placed in order:

```yaml
layout:
  rows: 3    # Number of rows
  cols: 3    # Number of columns
```

**Example layouts:**

**2x2 Grid (4 widgets):**
```
┌────────┬────────┐
│ Clock  │ Weather│
├────────┼────────┤
│ System │   IP   │
└────────┴────────┘
```

**3x3 Grid (9 widgets):**
```
┌────────┬────────┬────────┐
│ Clock  │Calendar│ Weather│
├────────┼────────┼────────┤
│ GitHub │ GitLab │ System │
├────────┼────────┼────────┤
│   IP   │ SMART  │  Text  │
└────────┴────────┴────────┘
```

**1x4 Grid (4 widgets, horizontal):**
```
┌────────┬────────┬────────┬────────┐
│ Clock  │ System │ Weather│   IP   │
└────────┴────────┴────────┴────────┘
```

### Widget Order

Widgets appear in this order (if configured):
1. Clock
2. Calendar
3. Weather
4. GitHub
5. GitLab
6. System Resources
7. IP Information
8. SMART Status
9. Text Viewer
10. Markdown Viewer

## Keyboard Controls

- **`q`**: Quit the application
- **`Esc`**: Quit the application
- **`Ctrl+C`**: Quit the application

The application automatically handles terminal resizing.

## Tips and Tricks

### Optimizing Refresh Intervals

- **Weather**: Update every 30-60 minutes (data doesn't change frequently)
- **GitHub/GitLab**: 5-10 minutes (be mindful of API rate limits)
- **System**: 1-5 seconds (minimal overhead)
- **IP**: 1 hour or more (rarely changes)

### API Rate Limits

- **GitHub**: 60 requests/hour (unauthenticated), 5,000 requests/hour (authenticated)
- **GitLab**: 10 requests/second (authenticated)

Use tokens for higher rate limits and access to private repositories.

### Performance Optimization

1. Reduce the number of monitored repositories
2. Increase refresh intervals for less critical widgets
3. Use a smaller grid layout (fewer widgets)
4. Disable widgets you don't need by leaving their config empty

### Troubleshooting

**Widget shows "Error":**
- Check your internet connection
- Verify API tokens are valid
- Check file paths for text/markdown viewers
- Ensure refresh intervals aren't too aggressive

**GitHub/GitLab rate limiting:**
- Add an API token to your config
- Increase refresh intervals
- Reduce the number of monitored repositories

**SMART data not available:**
- Install smartmontools: `apt install smartmontools` (Linux)
- Some systems require root access for SMART data
- The widget will still show basic disk information

## Example Configurations

### Developer Setup
```yaml
weather_location: "San Francisco"
github_token: "your_token"
github_repos:
  - "your-org/backend"
  - "your-org/frontend"
  - "your-org/mobile"
text_file: "./TODO.txt"
markdown_file: "./CHANGELOG.md"
layout:
  rows: 3
  cols: 3
```

### System Monitor
```yaml
weather_location: "London"
refresh_intervals:
  system: 2
  ip: 7200
layout:
  rows: 2
  cols: 2
```

### Minimal Dashboard
```yaml
weather_location: "Tokyo"
layout:
  rows: 2
  cols: 2
```

## Advanced Usage

### Running in a tmux pane
```bash
tmux new-session -d -s gotui './gotui'
tmux attach -t gotui
```

### Running as a systemd service
Create `/etc/systemd/system/gotui.service`:
```ini
[Unit]
Description=GoTUI Dashboard
After=network.target

[Service]
Type=simple
User=your-username
ExecStart=/usr/local/bin/gotui
Restart=always

[Install]
WantedBy=multi-user.target
```

### Using with multiple configs
```bash
# Development config
./gotui  # Uses ./config.yaml

# Production config
cp config.prod.yaml config.yaml
./gotui
```

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/cj3636/gotui/issues
- Documentation: https://github.com/cj3636/gotui/blob/main/README.md
