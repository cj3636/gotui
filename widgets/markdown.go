package widgets

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

const defaultMarkdown = "# Welcome to GoTUI\n\n" +
	"Use `MARKDOWN_PATH` to point at a local file. This widget renders Markdown using Glow so you can keep notes, dashboards, and runbooks nearby.\n"

// markdownWidget shows rendered markdown content using Glow.
type markdownWidget struct {
	content string
	width   int
}

// NewMarkdownWidget constructs the markdown widget.
func NewMarkdownWidget() Widget {
	return &markdownWidget{}
}

func (m *markdownWidget) Title() string { return "Markdown" }

func (m *markdownWidget) Init() tea.Cmd { return m.load() }

func (m *markdownWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case markdownMsg:
		m.content = msg.content
	}
	return m, nil
}

func (m *markdownWidget) View(width, _ int) string {
	if m.content == "" {
		return "Waiting for markdown data..."
	}
	rendererWidth := width
	if rendererWidth < 20 {
		rendererWidth = 20
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(rendererWidth),
	)
	if err != nil {
		return m.content
	}
	rendered, err := r.Render(strings.TrimSpace(m.content))
	if err != nil {
		return m.content
	}
	return rendered
}

type markdownMsg struct{ content string }

func (m *markdownWidget) load() tea.Cmd {
	return func() tea.Msg {
		path := os.Getenv("MARKDOWN_PATH")
		if path == "" {
			return markdownMsg{content: defaultMarkdown}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return markdownMsg{content: fmt.Sprintf("Failed to read %s: %v", path, err)}
		}
		return markdownMsg{content: string(data)}
	}
}
