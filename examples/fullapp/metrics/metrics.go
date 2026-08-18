// Package metrics is the fullapp example's second page: a spinner and a few
// gauges. It has no keys either, but unlike overview it runs a command, so
// it's what proves the shell keeps feeding ticks to a pane that doesn't have
// focus.
package metrics

import (
	"fmt"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/style"
)

// gaugeWidth is how wide a bar renders, before the label in front of it.
const gaugeWidth = 30

type gauge struct {
	name    string
	percent float64
}

// Model is the page. Values are hardcoded - this is a layout test, not a
// monitoring tool.
type Model struct {
	spinner spinner.Model
	bar     progress.Model
	gauges  []gauge

	width, height int
}

func New() Model {
	bar := bubbles.NewProgress()
	bar.SetWidth(gaugeWidth)
	return Model{
		spinner: bubbles.NewSpinner(spinner.WithSpinner(spinner.MiniDot)),
		bar:     bar,
		gauges: []gauge{
			{"cpu", 0.42},
			{"memory", 0.71},
			{"disk", 0.13},
		},
	}
}

// Init starts the spinner ticking. The shell returns it from its own Init,
// or the spinner never starts.
func (m Model) Init() tea.Cmd { return m.spinner.Tick }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
		return m, nil
	}

	// Every other message goes to the spinner, whose own tick keeps it
	// turning - including while another pane has focus.
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	rows := []string{style.S.Bold.Render(m.spinner.View() + " collecting")}
	for _, g := range m.gauges {
		// ViewAs renders a given percentage without animating toward it,
		// which is all a fixed value needs.
		rows = append(rows, fmt.Sprintf("%-8s %s", g.name, m.bar.ViewAs(g.percent)))
	}

	return tea.NewView(lipgloss.NewStyle().
		Width(m.width).Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, rows...)))
}
