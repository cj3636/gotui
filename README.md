# Go TUI Dashboard

A modular and extensible Terminal User Interface (TUI) dashboard application built with Go and Bubble Tea.

## Overview

Go TUI Dashboard is a comprehensive terminal-based dashboard that provides various widgets for monitoring and interacting with different services and system resources. The application is designed to be modular, allowing easy extension with new widgets and features.

## Planned Features

### Widgets
- **Weather Widget**: Real-time weather information using `wttr.in`
- **GitHub/GitLab Integration**: Monitor repositories, snippets, gists, and To-Dos
- **System Monitoring**: Track CPU, memory, and disk usage
- **SMART Drive Status**: Monitor drive health and status
- **IP Information**: Display network and IP details
- **Text File Viewer**: View and navigate text files
- **Markdown Rendering**: Beautiful markdown display using `glow`
- **Clock/Calendar**: Current time and calendar information

### Technical Stack
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the core TUI framework
- **Interactions**: [Gum](https://github.com/charmbracelet/gum) for simple user interactions
- **Markdown**: [Glow](https://github.com/charmbracelet/glow) for markdown rendering

### Design Principles
- **Modular Architecture**: Widgets are self-contained and easy to add or remove
- **Responsive Design**: Dashboard automatically resizes and scales with terminal window
- **Extensible**: Easy to add new widgets and features

## Getting Started

### Prerequisites
- Go 1.24 or higher

### Installation

```bash
git clone https://github.com/cj3636/gotui.git
cd gotui
go build
```

### Running

```bash
./gotui
```

## Development

This project is currently in early development. Contributions and feature requests are welcome!

## License

TBD
