package widgets

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ipWidget struct {
	addresses []string
}

// NewIPWidget constructs the IP address widget.
func NewIPWidget() Widget { return &ipWidget{} }

func (i *ipWidget) Title() string { return "IP Info" }

func (i *ipWidget) Init() tea.Cmd { return i.refresh() }

func (i *ipWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		return i, i.refresh()
	case ipMsg:
		data := msg.(ipMsg)
		i.addresses = data.addresses
		return i, Tick(1 * time.Minute)
	}
	return i, nil
}

func (i *ipWidget) View(width, _ int) string {
	if len(i.addresses) == 0 {
		return "Discovering network interfaces..."
	}
	return strings.Join(i.addresses, "\n")
}

type ipMsg struct{ addresses []string }

func (i *ipWidget) refresh() tea.Cmd {
	return func() tea.Msg {
		var addrs []string
		interfaces, err := net.Interfaces()
		if err != nil {
			return ipMsg{addresses: []string{fmt.Sprintf("Error: %v", err)}}
		}
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp == 0 {
				continue
			}
			iaddrs, _ := iface.Addrs()
			for _, addr := range iaddrs {
				addrs = append(addrs, fmt.Sprintf("%s: %s", iface.Name, addr.String()))
			}
		}
		sort.Strings(addrs)
		if len(addrs) == 0 {
			addrs = []string{"No active addresses"}
		}
		return ipMsg{addresses: addrs}
	}
}
