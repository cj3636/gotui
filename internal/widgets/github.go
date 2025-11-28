package widgets

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// GithubWidget displays GitHub repository information
type GithubWidget struct {
	BaseWidget
	token          string
	repos          []string
	repoInfo       []RepoInfo
	err            error
	lastUpdate     time.Time
	updateInterval time.Duration
}

// RepoInfo contains repository information
type RepoInfo struct {
	Name        string
	Stars       int
	Forks       int
	OpenIssues  int
	OpenPRs     int
	Description string
}

// GithubMsg contains GitHub information
type GithubMsg struct {
	repos []RepoInfo
	err   error
}

// GithubRefreshMsg signals it's time to refresh GitHub info
type GithubRefreshMsg struct{}

// NewGithubWidget creates a new GitHub widget
func NewGithubWidget(token string, repos []string, refreshInterval int) *GithubWidget {
	return &GithubWidget{
		BaseWidget:     NewBaseWidget("🐙 GitHub"),
		token:          token,
		repos:          repos,
		updateInterval: time.Duration(refreshInterval) * time.Second,
	}
}

// Init initializes the widget
func (w *GithubWidget) Init() tea.Cmd {
	return w.fetchGithubInfo()
}

// Update handles messages
func (w *GithubWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case GithubMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.repoInfo = msg.repos
			w.err = nil
		}
		w.lastUpdate = time.Now()
		return w, tea.Tick(w.updateInterval, func(t time.Time) tea.Msg {
			return GithubRefreshMsg{}
		})
	case GithubRefreshMsg:
		return w, w.fetchGithubInfo()
	}
	return w, nil
}

// View renders the widget
func (w *GithubWidget) View() string {
	var content string
	if w.err != nil {
		content = fmt.Sprintf("Error: %v", w.err)
	} else if len(w.repoInfo) == 0 {
		if len(w.repos) == 0 {
			content = "No repositories configured"
		} else {
			content = "Loading..."
		}
	} else {
		var lines []string
		for _, repo := range w.repoInfo {
			lines = append(lines, fmt.Sprintf(
				"%s\n⭐ %d  🍴 %d  📝 %d",
				repo.Name,
				repo.Stars,
				repo.Forks,
				repo.OpenIssues,
			))
		}
		content = strings.Join(lines, "\n\n")
	}
	return w.RenderContent(content)
}

func (w *GithubWidget) fetchGithubInfo() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var client *github.Client
		if w.token != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: w.token},
			)
			tc := oauth2.NewClient(ctx, ts)
			client = github.NewClient(tc)
		} else {
			client = github.NewClient(nil)
		}

		var repos []RepoInfo
		for _, repoName := range w.repos {
			parts := strings.Split(repoName, "/")
			if len(parts) != 2 {
				continue
			}
			owner, repo := parts[0], parts[1]

			repoData, _, err := client.Repositories.Get(ctx, owner, repo)
			if err != nil {
				return GithubMsg{err: err}
			}

			// Get pull requests count
			prs, _, _ := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
				State: "open",
			})

			info := RepoInfo{
				Name:       repoName,
				Stars:      repoData.GetStargazersCount(),
				Forks:      repoData.GetForksCount(),
				OpenIssues: repoData.GetOpenIssuesCount(),
				OpenPRs:    len(prs),
			}
			repos = append(repos, info)
		}

		return GithubMsg{repos: repos}
	}
}
