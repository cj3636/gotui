package widgets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type gitlabUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	WebURL   string `json:"web_url"`
}

type gitlabWidget struct {
	user    gitlabUser
	message string
	err     error
}

// NewGitLabWidget constructs the GitLab widget.
func NewGitLabWidget() Widget { return &gitlabWidget{} }

func (g *gitlabWidget) Title() string { return "GitLab" }

func (g *gitlabWidget) Init() tea.Cmd { return g.fetch() }

func (g *gitlabWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return g, g.fetch()
	case gitlabMsg:
		g.user = msg.user
		g.message = msg.message
		g.err = msg.err
		return g, Tick(10 * time.Minute)
	}
	return g, nil
}

func (g *gitlabWidget) View(width, _ int) string {
	if g.err != nil {
		return fmt.Sprintf("Error: %v", g.err)
	}
	if g.user.Username == "" && g.message == "" {
		return "Loading profile..."
	}
	if g.user.Username == "" {
		return g.message
	}
	return fmt.Sprintf("User: %s\nName: %s\nURL: %s", g.user.Username, g.user.Name, g.user.WebURL)
}

type gitlabMsg struct {
	user    gitlabUser
	message string
	err     error
}

func (g *gitlabWidget) fetch() tea.Cmd {
	return func() tea.Msg {
		token := os.Getenv("GITLAB_TOKEN")
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", "https://gitlab.com/api/v4/user", nil)
		if token != "" {
			req.Header.Set("PRIVATE-TOKEN", token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return gitlabMsg{err: err}
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized {
			return gitlabMsg{message: "Set GITLAB_TOKEN for private data"}
		}
		if resp.StatusCode != http.StatusOK {
			return gitlabMsg{message: fmt.Sprintf("API status: %s", resp.Status)}
		}
		var user gitlabUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return gitlabMsg{err: err}
		}
		return gitlabMsg{user: user}
	}
}
