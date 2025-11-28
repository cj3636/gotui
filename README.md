# GoTUI Dashboard

A modular Bubble Tea dashboard inspired by wtfutil and built with Charmbracelet components. Widgets resize with the terminal and can be extended easily.

## Widgets
- **Clock** – Live time updates.
- **Weather** – Current conditions from [wttr.in](https://wttr.in). Override location with `WTTR_LOCATION` and pass custom query strings (e.g., `format=3&M` for metric units) via `WTTR_PARAMS`.
- **Moon Phase** – wttr.in `/moon` view with customizable parameters through `WTTR_MOON_PARAMS`.
- **System** – CPU, memory, disk usage, Go runtime version, and SMART availability hint.
- **IP Info** – Active interface addresses.
- **Markdown** – Renders Markdown through Glow. Provide `MARKDOWN_PATH` to point at a local file.
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
   - `MARKDOWN_PATH` – Path to a Markdown file rendered by Glow.
   - `GITHUB_TOKEN` – Personal access token for private/public GitHub data.
   - `GITLAB_TOKEN` – Personal access token for GitLab profile data.

The dashboard adapts to your terminal size and uses two columns when space allows.

## Roadmap

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
