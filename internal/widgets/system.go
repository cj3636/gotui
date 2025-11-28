package widgets

import (
	"fmt"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemWidget displays system resource information
type SystemWidget struct {
	BaseWidget
	cpuPercent     float64
	memPercent     float64
	memUsed        uint64
	memTotal       uint64
	diskPercent    float64
	diskUsed       uint64
	diskTotal      uint64
	updateInterval time.Duration
}

// SystemMsg contains system information
type SystemMsg struct {
	cpuPercent  float64
	memPercent  float64
	memUsed     uint64
	memTotal    uint64
	diskPercent float64
	diskUsed    uint64
	diskTotal   uint64
}

// SystemRefreshMsg signals it's time to refresh system info
type SystemRefreshMsg struct{}

// NewSystemWidget creates a new system resource widget
func NewSystemWidget(refreshInterval int) *SystemWidget {
	return &SystemWidget{
		BaseWidget:     NewBaseWidget("💻 System Resources"),
		updateInterval: time.Duration(refreshInterval) * time.Second,
	}
}

// Init initializes the widget
func (w *SystemWidget) Init() tea.Cmd {
	return w.fetchSystemInfo()
}

// Update handles messages
func (w *SystemWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case SystemMsg:
		w.cpuPercent = msg.cpuPercent
		w.memPercent = msg.memPercent
		w.memUsed = msg.memUsed
		w.memTotal = msg.memTotal
		w.diskPercent = msg.diskPercent
		w.diskUsed = msg.diskUsed
		w.diskTotal = msg.diskTotal
		return w, tea.Tick(w.updateInterval, func(t time.Time) tea.Msg {
			return SystemRefreshMsg{}
		})
	case SystemRefreshMsg:
		return w, w.fetchSystemInfo()
	}
	return w, nil
}

// View renders the widget
func (w *SystemWidget) View() string {
	content := fmt.Sprintf(
		"CPU:  %5.1f%%\n"+
			"RAM:  %5.1f%% (%s / %s)\n"+
			"Disk: %5.1f%% (%s / %s)\n"+
			"OS:   %s/%s",
		w.cpuPercent,
		w.memPercent,
		formatBytes(w.memUsed),
		formatBytes(w.memTotal),
		w.diskPercent,
		formatBytes(w.diskUsed),
		formatBytes(w.diskTotal),
		runtime.GOOS,
		runtime.GOARCH,
	)
	return w.RenderContent(content)
}

func (w *SystemWidget) fetchSystemInfo() tea.Cmd {
	return func() tea.Msg {
		// Get CPU percentage
		cpuPercentages, _ := cpu.Percent(time.Second, false)
		cpuPercent := 0.0
		if len(cpuPercentages) > 0 {
			cpuPercent = cpuPercentages[0]
		}

		// Get memory stats
		memStats, _ := mem.VirtualMemory()
		memPercent := 0.0
		memUsed := uint64(0)
		memTotal := uint64(0)
		if memStats != nil {
			memPercent = memStats.UsedPercent
			memUsed = memStats.Used
			memTotal = memStats.Total
		}

		// Get disk stats
		diskStats, _ := disk.Usage("/")
		diskPercent := 0.0
		diskUsed := uint64(0)
		diskTotal := uint64(0)
		if diskStats != nil {
			diskPercent = diskStats.UsedPercent
			diskUsed = diskStats.Used
			diskTotal = diskStats.Total
		}

		return SystemMsg{
			cpuPercent:  cpuPercent,
			memPercent:  memPercent,
			memUsed:     memUsed,
			memTotal:    memTotal,
			diskPercent: diskPercent,
			diskUsed:    diskUsed,
			diskTotal:   diskTotal,
		}
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
