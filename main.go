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
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cj3636/gotui/internal/app"
)

func main() {
	// Load configuration
	config, err := app.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create the main application model
	model := app.NewModel(config)

	// Initialize the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
