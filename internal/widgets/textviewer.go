package widgets

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TextViewerWidget displays content from a text file
type TextViewerWidget struct {
	BaseWidget
	filename string
	content  string
	err      error
}

// NewTextViewerWidget creates a new text viewer widget
func NewTextViewerWidget(filename string) *TextViewerWidget {
	return &TextViewerWidget{
		BaseWidget: NewBaseWidget("📄 Text Viewer"),
		filename:   filename,
	}
}

// Init initializes the widget
func (w *TextViewerWidget) Init() tea.Cmd {
	if w.filename != "" {
		return w.loadFile()
	}
	return nil
}

// Update handles messages
func (w *TextViewerWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	return w, nil
}

// View renders the widget
func (w *TextViewerWidget) View() string {
	var content string
	if w.filename == "" {
		content = "No file configured"
	} else if w.err != nil {
		content = fmt.Sprintf("Error: %v", w.err)
	} else if w.content == "" {
		content = "Loading..."
	} else {
		content = w.content
	}
	return w.RenderContent(content)
}

func (w *TextViewerWidget) loadFile() tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(w.filename)
		if err != nil {
			w.err = err
			return nil
		}

		// Limit content to first 100 lines to avoid overwhelming the display
		lines := strings.Split(string(data), "\n")
		if len(lines) > 100 {
			lines = lines[:100]
			lines = append(lines, "... (truncated)")
		}

		w.content = strings.Join(lines, "\n")
		return nil
	}
}
