package widgets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type githubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
}

type githubWidget struct {
	user    githubUser
	message string
	err     error
}

// NewGitHubWidget constructs the GitHub widget.
func NewGitHubWidget() Widget { return &githubWidget{} }

func (g *githubWidget) Title() string { return "GitHub" }

func (g *githubWidget) Init() tea.Cmd { return g.fetch() }

func (g *githubWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return g, g.fetch()
	case githubMsg:
		g.user = msg.user
		g.message = msg.message
		g.err = msg.err
		return g, Tick(10 * time.Minute)
	}
	return g, nil
}

func (g *githubWidget) View(width, _ int) string {
	if g.err != nil {
		return fmt.Sprintf("Error: %v", g.err)
	}
	if g.user.Login == "" && g.message == "" {
		return "Loading profile..."
	}
	if g.user.Login == "" {
		return g.message
	}
	return fmt.Sprintf("User: %s\nName: %s\nRepos: %d\nFollowers: %d", g.user.Login, g.user.Name, g.user.PublicRepos, g.user.Followers)
}

type githubMsg struct {
	user    githubUser
	message string
	err     error
}

func (g *githubWidget) fetch() tea.Cmd {
	return func() tea.Msg {
		token := os.Getenv("GITHUB_TOKEN")
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
		if token != "" {
			req.Header.Set("Authorization", "token "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return githubMsg{err: err}
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized {
			return githubMsg{message: "Set GITHUB_TOKEN to load private data"}
		}
		if resp.StatusCode != http.StatusOK {
			return githubMsg{message: fmt.Sprintf("API status: %s", resp.Status)}
		}
		var user githubUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return githubMsg{err: err}
		}
		return githubMsg{user: user}
	}
}
