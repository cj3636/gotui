package widgets

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

// GitlabWidget displays GitLab project information
type GitlabWidget struct {
	BaseWidget
	token          string
	projects       []string
	projectInfo    []ProjectInfo
	err            error
	lastUpdate     time.Time
	updateInterval time.Duration
}

// ProjectInfo contains project information
type ProjectInfo struct {
	Name       string
	Stars      int
	Forks      int
	OpenIssues int
	OpenMRs    int
}

// GitlabMsg contains GitLab information
type GitlabMsg struct {
	projects []ProjectInfo
	err      error
}

// NewGitlabWidget creates a new GitLab widget
func NewGitlabWidget(token string, projects []string, refreshInterval int) *GitlabWidget {
	return &GitlabWidget{
		BaseWidget:     NewBaseWidget("🦊 GitLab"),
		token:          token,
		projects:       projects,
		updateInterval: time.Duration(refreshInterval) * time.Second,
	}
}

// Init initializes the widget
func (w *GitlabWidget) Init() tea.Cmd {
	return w.fetchGitlabInfo()
}

// Update handles messages
func (w *GitlabWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case GitlabMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.projectInfo = msg.projects
			w.err = nil
		}
		w.lastUpdate = time.Now()
		return w, tea.Tick(w.updateInterval, func(t time.Time) tea.Msg {
			return w.fetchGitlabInfo()()
		})
	}
	return w, nil
}

// View renders the widget
func (w *GitlabWidget) View() string {
	var content string
	if w.err != nil {
		content = fmt.Sprintf("Error: %v", w.err)
	} else if len(w.projectInfo) == 0 {
		if len(w.projects) == 0 {
			content = "No projects configured"
		} else {
			content = "Loading..."
		}
	} else {
		var lines []string
		for _, proj := range w.projectInfo {
			lines = append(lines, fmt.Sprintf(
				"%s\n⭐ %d  🍴 %d  📝 %d  🔀 %d",
				proj.Name,
				proj.Stars,
				proj.Forks,
				proj.OpenIssues,
				proj.OpenMRs,
			))
		}
		content = strings.Join(lines, "\n\n")
	}
	return w.RenderContent(content)
}

func (w *GitlabWidget) fetchGitlabInfo() tea.Cmd {
	return func() tea.Msg {
		if len(w.projects) == 0 {
			return GitlabMsg{projects: []ProjectInfo{}}
		}

		var client *gitlab.Client
		var err error
		if w.token != "" {
			client, err = gitlab.NewClient(w.token)
		} else {
			client, err = gitlab.NewClient("")
		}
		if err != nil {
			return GitlabMsg{err: err}
		}

		var projects []ProjectInfo
		for _, projectName := range w.projects {
			project, _, err := client.Projects.GetProject(projectName, nil)
			if err != nil {
				return GitlabMsg{err: err}
			}

			// Get merge requests count
			openState := "opened"
			mrs, _, _ := client.MergeRequests.ListProjectMergeRequests(project.ID, &gitlab.ListProjectMergeRequestsOptions{
				State: &openState,
			})

			info := ProjectInfo{
				Name:       projectName,
				Stars:      project.StarCount,
				Forks:      project.ForksCount,
				OpenIssues: project.OpenIssuesCount,
				OpenMRs:    len(mrs),
			}
			projects = append(projects, info)
		}

		return GitlabMsg{projects: projects}
	}
}
