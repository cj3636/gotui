# GoTUI Dashboard

<<<<<<< codex/add-tui-dashboard-app-in-go
GoTUI is a modular terminal dashboard inspired by wtfutil and built on Charmbracelet's Bubble Tea stack. It arranges independent widgets that adapt to your terminal size, using a single-column layout for narrow viewports and a responsive two-column layout when there is enough space.

## Features
- Bubble Tea dashboard shell with automatic resizing and simple rounded panel styling via Lip Gloss.
- Pluggable widget interface (`Widget` with `Init`, `Update`, `View`, `Title`) so components can be swapped or extended easily.
- Network-backed widgets for GitHub, GitLab, wttr.in weather, and moon phase data with token-aware messaging.
- Local system widgets for clock, Markdown rendering (Glamour/Glow), IP discovery, and system metrics (CPU, memory, disk, Go version with a SMART availability hint).
- Sensible defaults with environment variable overrides for tokens, locations, markdown path, and wttr.in query parameters.

## Widgets
- **Clock** – Live time updates via `tea.Tick`.
- **Weather** – Current conditions from [wttr.in](https://wttr.in). Override the location with `WTTR_LOCATION` (spaces become underscores), choose units with `WTTR_UNITS` (`m`, `u`, optional `M` for m/s wind), and pick view flags via `WTTR_VIEW` (e.g., `Fq1`). Default request: `https://wttr.in/89701?uFq1`.
- **Moon Phase** – wttr.in `/moon` view with customizable location (`WTTR_MOON_LOCATION`), units (`WTTR_MOON_UNITS`), and view flags (`WTTR_MOON_VIEW`, default `Fq1`).
- **System** – CPU, memory, disk usage, Go runtime version, and SMART availability hint using `gopsutil`.
- **IP Info** – Active interface addresses so you can quickly see reachable IPs.
- **Markdown** – Renders Markdown through Glamour (Glow's renderer). Provide `MARKDOWN_PATH` to point at a local file; otherwise renders a helpful default message.
- **GitHub** – Authenticated profile summary when `GITHUB_TOKEN` is set; surfaces status guidance when unauthenticated.
- **GitLab** – Authenticated profile summary when `GITLAB_TOKEN` is set; surfaces status guidance when unauthenticated.

## Configuration
Set environment variables before running the dashboard to customize widget behavior:

| Variable | Description | Default |
| --- | --- | --- |
| `WTTR_LOCATION` | Weather location/endpoint segment passed to wttr.in (spaces converted to `_`). | `89701` |
| `WTTR_UNITS` | Weather units: `m` (metric) or `u` (US); append `M` for wind in m/s. | `u` |
| `WTTR_VIEW` | Weather view flags from [wttr.in help](https://wttr.in/:help) such as `Fq1`. | `Fq1` |
| `WTTR_MOON_LOCATION` | wttr.in moon endpoint segment (typically `moon`). | `moon` |
| `WTTR_MOON_UNITS` | Units for moon endpoint (same as weather units). | `u` |
| `WTTR_MOON_VIEW` | View flags for moon endpoint. | `Fq1` |
| `MARKDOWN_PATH` | Local Markdown file to render in the Markdown widget. | *(embedded welcome copy)* |
| `GITHUB_TOKEN` | GitHub personal access token for private/public profile calls. | *(unauthenticated request)* |
| `GITLAB_TOKEN` | GitLab personal access token for profile calls. | *(unauthenticated request)* |
| `WIDGET_HEIGHT_<TITLE>` | Optional per-widget vertical sizing multiplier (e.g., `WIDGET_HEIGHT_WEATHER=2`). | `1` |

> Quit the dashboard with `q` or `Ctrl+C`.

## Running
1. Install Go 1.25 or newer.
2. Download dependencies and run the dashboard:
   ```bash
   go run ./...
   ```
3. Optionally build a binary:
   ```bash
   go build -o gotui ./...
   ./gotui
   ```
4. Export any environment variables from the configuration table above to tailor the widgets.

The dashboard adapts to your terminal size and uses two columns when space allows. Each widget self-reschedules with sensible refresh intervals (e.g., 30 minutes for wttr.in, 5 seconds for system stats).

## Architecture
- **Layout** – `widgets/dashboard.go` coordinates Bubble Tea updates, tracks `tea.WindowSizeMsg`, and renders panels with Lip Gloss. Column width is computed from the terminal width, falling back to a single column for narrow windows.
- **Widget contract** – Widgets implement `Init`, `Update`, `View`, and `Title`. The dashboard batches initial commands and fans out updates to each widget on every message.
- **Tick helper** – `Tick` in `dashboard.go` emits periodic `TickMsg` messages so widgets can refresh without blocking the main event loop.
- **Data fetching** – Network widgets use short-lived `http.Client` instances with timeouts. GitHub/GitLab widgets surface helpful messages when tokens are missing or unauthorized responses are returned.

## Development
- Format code with `gofmt -w .`.
- Resolve modules with `go mod tidy` (internet access required).
- Run tests (if/when added) with `go test ./...`.
- Keep widget implementations focused and independent; prefer small structs with simple view rendering to preserve readability inside the TUI.

## Roadmap
=======
A modular Bubble Tea dashboard inspired by wtfutil and built with Charmbracelet components. Widgets resize with the terminal and can be extended easily.

## Widgets
- **Clock** – Live time updates.
- **Weather** – Current conditions from [wttr.in](https://wttr.in). Override location with `WTTR_LOCATION` and pass custom query strings (e.g., `format=3&M` for metric units) via `WTTR_PARAMS`.
- **Moon Phase** – wttr.in `/moon` view with customizable parameters through `WTTR_MOON_PARAMS`.
- **System** – CPU, memory, disk usage, Go runtime version, and SMART availability hint.
- **IP Info** – Active interface addresses.
- **Markdown** – Renders Markdown through Glamour (the renderer used by Glow). Provide `MARKDOWN_PATH` to point at a local file.
- **GitHub** – Authenticated profile summary when `GITHUB_TOKEN` is set.
- **GitLab** – Authenticated profile summary when `GITLAB_TOKEN` is set.

## Running
1. Install Go 1.25 or newer.
2. Download dependencies and run the app:
   ```bash
   go run ./...
   ```
3. Optional environment variables:
   - `WTTR_LOCATION` – Weather location (city or query supported by wttr.in).
   - `WTTR_PARAMS` – Optional query string appended to the weather request (defaults to `format=3`).
   - `WTTR_MOON_LOCATION` – Override the `/moon` endpoint segment if you want a location-specific moon forecast.
   - `WTTR_MOON_PARAMS` – Optional query string appended to the moon request (defaults to `format=Moon:+%m`).
   - `MARKDOWN_PATH` – Path to a Markdown file rendered by Glamour (used by Glow).
   - `GITHUB_TOKEN` – Personal access token for private/public GitHub data.
   - `GITLAB_TOKEN` – Personal access token for GitLab profile data.

The dashboard adapts to your terminal size and uses two columns when space allows.

## Roadmap

>>>>>>> main
### Near-term widgets
- **Enhanced weather**: expose presets for rich formats (e.g., multi-line forecasts), caching to reduce repeated calls, and per-widget refresh intervals.
- **Moon phase**: add phase icons/emoji, configurable locales, and optional sunrise/sunset overlays using wttr.in fields.
- **Calendar/agenda**: integrate ICS parsing and simple TODO list display.
- **Notifications**: surface API errors or threshold alerts (e.g., high CPU) via a lightweight status bar.

### Core improvements
- **Configuration system**: introduce a structured config file (YAML/TOML) to define widget ordering, sizing hints, refresh cadences, and secrets (with env overrides).
- **Theme and layout**: selectable color schemes, compact/dense modes, and per-widget padding/border toggles.
- **Persistence/DB integration**: optional SQLite layer for caching API responses, storing historical metrics, and persisting user preferences. Start with a persistence interface, then add SQLite implementation with migration scaffolding and background pruning.
- **Plugin model**: define a registry for third-party widgets with discovery and sandboxed execution guards.
- **Observability**: structured logging toggle, debug overlay with last refresh timestamps, and a health pane summarizing widget status.

### Development milestones
1. **Config + themes**: scaffold config parsing, default theme profiles, and validation.
2. **Widget upgrades**: expand weather/moon capabilities, add calendar/agenda, and implement notification surface.
3. **Persistence layer**: introduce the storage interface, add SQLite integration, and thread caching into networked widgets (weather, GitHub/GitLab, IP lookup).
4. **Plugin registry**: allow registering external widgets via config, with safety rails and opt-in execution.
5. **Polish and QA**: integration tests for layout scaling, linting/formatting automation, and sample configs demonstrating advanced setups.
<<<<<<< codex/add-tui-dashboard-app-in-go

## Notes
- Markdown rendering uses Glamour, which powers Charmbracelet's Glow CLI. Glow can still be used to render Markdown externally; in the dashboard we keep rendering lightweight and non-fullscreen.
- SMART status is reported as a hint (smartctl presence) rather than querying drives directly to avoid portability issues.
=======
>>>>>>> main
