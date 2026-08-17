package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/tabs"
)

// pane is a minimal tabs.Tab implementation. It keeps its own counter to
// show that each tab's model has independent state that persists across
// switches, and its own width/height to show how a host propagates size
// down to a wrapped Tab (see model.Update's tea.WindowSizeMsg case).
type pane struct {
	name          string
	count         int
	width, height int
}

func newPane(name string) pane {
	return pane{name: name}
}

func (p pane) Init() tea.Cmd {
	return nil
}

func (p pane) Update(msg tea.Msg) (tabs.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width, p.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if msg.String() == "+" {
			p.count++
		}
	}
	return p, nil
}

func (p pane) View() string {
	return fmt.Sprintf("%s\n\npress + to increment: %d\n(content area: %dx%d)", p.name, p.count, p.width, p.height)
}

var docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

type model struct {
	tabs tabs.Model
}

func newModel() model {
	items := []tabs.Item{
		{Title: "Lip Gloss", Model: newPane("Lip Gloss")},
		{Title: "Blush", Model: newPane("Blush")},
		{Title: "Eye Shadow", Model: newPane("Eye Shadow")},
		{Title: "Mascara", Model: newPane("Mascara")},
		{Title: "Foundation", Model: newPane("Foundation")},
	}

	return model{tabs: tabs.New(items)}
}

func (m model) Init() tea.Cmd {
	return m.tabs.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Size the tabs component to fill the terminal, net of docStyle's own
		// frame. The tab bar itself keeps its natural width; Content stretches.
		m.tabs.SetSize(
			msg.Width-docStyle.GetHorizontalFrameSize(),
			msg.Height-docStyle.GetVerticalFrameSize(),
		)

		// tabs has no generic way to size an arbitrary Tab itself, so forward
		// the actual usable content area as a WindowSizeMsg: tabs.Update
		// already routes non-key messages to the active item's Update, so
		// this reaches pane.Update's own tea.WindowSizeMsg case above.
		var cmd tea.Cmd
		m.tabs, cmd = m.tabs.Update(tea.WindowSizeMsg{
			Width:  m.tabs.ContentWidth(),
			Height: m.tabs.ContentHeight(),
		})
		return m, cmd
	}

	var cmd tea.Cmd
	m.tabs, cmd = m.tabs.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	view := tea.NewView(docStyle.Render(m.tabs.View()))
	view.AltScreen = true
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
