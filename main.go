package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"gotui/widgets"
)

func main() {
	model := widgets.NewDashboard()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("failed to start program: %v", err)
	}
}
