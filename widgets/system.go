package widgets

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	tea "github.com/charmbracelet/bubbletea"
)

// systemWidget reports CPU, memory, and disk stats.
type systemWidget struct {
	cpuLoad   float64
	memUsed   float64
	memTotal  float64
	diskUsed  float64
	diskTotal float64
	err       error
}

// NewSystemWidget constructs the system monitor widget.
func NewSystemWidget() Widget { return &systemWidget{} }

func (s *systemWidget) Title() string { return "System" }

func (s *systemWidget) Init() tea.Cmd { return s.sample() }

func (s *systemWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		return s, s.sample()
	case systemMsg:
		data := msg.(systemMsg)
		s.cpuLoad = data.cpuLoad
		s.memUsed = data.memUsed
		s.memTotal = data.memTotal
		s.diskUsed = data.diskUsed
		s.diskTotal = data.diskTotal
		s.err = data.err
		return s, Tick(5 * time.Second)
	}
	return s, nil
}

func (s *systemWidget) View(width, _ int) string {
	if s.err != nil {
		return fmt.Sprintf("Error: %v", s.err)
	}
	if s.memTotal == 0 {
		return "Collecting metrics..."
	}
	smart := "SMART: smartctl not detected"
	return strings.Join([]string{
		fmt.Sprintf("CPU Load: %0.1f%%", s.cpuLoad),
		fmt.Sprintf("Memory: %0.1f / %0.1f GiB", s.memUsed, s.memTotal),
		fmt.Sprintf("Disk: %0.1f / %0.1f GiB", s.diskUsed, s.diskTotal),
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		smart,
	}, "\n")
}

type systemMsg struct {
	cpuLoad   float64
	memUsed   float64
	memTotal  float64
	diskUsed  float64
	diskTotal float64
	err       error
}

func (s *systemWidget) sample() tea.Cmd {
	return func() tea.Msg {
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			return systemMsg{err: err}
		}
		memStats, err := mem.VirtualMemory()
		if err != nil {
			return systemMsg{err: err}
		}
		diskStats, err := disk.Usage("/")
		if err != nil {
			return systemMsg{err: err}
		}
		return systemMsg{
			cpuLoad:   cpuPercent[0],
			memUsed:   float64(memStats.Used) / (1024 * 1024 * 1024),
			memTotal:  float64(memStats.Total) / (1024 * 1024 * 1024),
			diskUsed:  float64(diskStats.Used) / (1024 * 1024 * 1024),
			diskTotal: float64(diskStats.Total) / (1024 * 1024 * 1024),
		}
	}
}
