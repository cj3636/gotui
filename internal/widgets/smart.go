package widgets

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// SMARTWidget displays SMART drive status
type SMARTWidget struct {
	BaseWidget
	smartData string
	err       error
}

// SMARTMsg contains SMART data
type SMARTMsg struct {
	data string
	err  error
}

// NewSMARTWidget creates a new SMART status widget
func NewSMARTWidget() *SMARTWidget {
	return &SMARTWidget{
		BaseWidget: NewBaseWidget("💾 SMART Status"),
	}
}

// Init initializes the widget
func (w *SMARTWidget) Init() tea.Cmd {
	return w.fetchSMARTData()
}

// Update handles messages
func (w *SMARTWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case SMARTMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.smartData = msg.data
			w.err = nil
		}
	}
	return w, nil
}

// View renders the widget
func (w *SMARTWidget) View() string {
	var content string
	if w.err != nil {
		content = fmt.Sprintf("Error: %v\n\nNote: SMART data requires\nsmartmontools and root access", w.err)
	} else if w.smartData == "" {
		content = "Loading SMART data..."
	} else {
		content = w.smartData
	}
	return w.RenderContent(content)
}

func (w *SMARTWidget) fetchSMARTData() tea.Cmd {
	return func() tea.Msg {
		// Check if smartctl is available
		var cmd *exec.Cmd

		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			// Try to get basic disk info without root
			cmd = exec.Command("df", "-h", "/")
		} else {
			return SMARTMsg{err: fmt.Errorf("SMART monitoring not supported on %s", runtime.GOOS)}
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			return SMARTMsg{err: err}
		}

		// Parse disk information
		lines := strings.Split(string(output), "\n")
		var result []string
		result = append(result, "Disk Status:")
		if len(lines) > 1 {
			// Add headers and data
			fields := strings.Fields(lines[0])
			if len(fields) >= 5 {
				result = append(result, "")
				result = append(result, strings.Join(fields[:5], " "))
			}

			dataFields := strings.Fields(lines[1])
			if len(dataFields) >= 5 {
				result = append(result, strings.Join(dataFields[:5], " "))
			}
		}

		// Try to get smartctl info (will fail without root)
		smartCmd := exec.Command("smartctl", "--scan")
		smartOutput, err := smartCmd.CombinedOutput()
		if err == nil && len(smartOutput) > 0 {
			result = append(result, "")
			result = append(result, "Detected drives:")
			driveLines := strings.Split(string(smartOutput), "\n")
			for _, line := range driveLines {
				if strings.TrimSpace(line) != "" {
					parts := strings.Fields(line)
					if len(parts) > 0 {
						result = append(result, "  "+parts[0])
					}
				}
			}
		} else {
			result = append(result, "")
			result = append(result, "Install smartmontools")
			result = append(result, "for detailed SMART data")
		}

		data := strings.Join(result, "\n")
		return SMARTMsg{data: data}
	}
}
