package widgets

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// MarkdownWidget displays rendered markdown content
type MarkdownWidget struct {
	BaseWidget
	filename string
	content  string
	err      error
}

// MarkdownMsg contains rendered markdown content
type MarkdownMsg struct {
	content string
	err     error
}

// NewMarkdownWidget creates a new markdown viewer widget
func NewMarkdownWidget(filename string) *MarkdownWidget {
	return &MarkdownWidget{
		BaseWidget: NewBaseWidget("📝 Markdown"),
		filename:   filename,
	}
}

// Init initializes the widget
func (w *MarkdownWidget) Init() tea.Cmd {
	if w.filename != "" {
		return w.loadMarkdown()
	}
	return nil
}

// Update handles messages
func (w *MarkdownWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case MarkdownMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.content = msg.content
			w.err = nil
		}
	}
	return w, nil
}

// View renders the widget
func (w *MarkdownWidget) View() string {
	var content string
	if w.filename == "" {
		content = "No markdown file configured"
	} else if w.err != nil {
		content = fmt.Sprintf("Error: %v", w.err)
	} else if w.content == "" {
		content = "Loading..."
	} else {
		content = w.content
	}
	return w.RenderContent(content)
}

func (w *MarkdownWidget) loadMarkdown() tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(w.filename)
		if err != nil {
			return MarkdownMsg{err: err}
		}

		// Use glamour to render markdown
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(60),
		)
		if err != nil {
			return MarkdownMsg{err: err}
		}

		rendered, err := r.Render(string(data))
		if err != nil {
			return MarkdownMsg{err: err}
		}

		// Limit content to avoid overwhelming the display
		lines := strings.Split(rendered, "\n")
		if len(lines) > 100 {
			lines = lines[:100]
			lines = append(lines, "... (truncated)")
		}

		content := strings.Join(lines, "\n")
		return MarkdownMsg{content: content}
	}
}
