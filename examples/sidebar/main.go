// Command sidebar is the whole app-shell pattern in one file: a sidebar on
// the left, an ordinary tea.Model on the right, a global help bar at the
// bottom, tab to move focus between the two. No layout package involved -
// lipgloss.JoinHorizontal/JoinVertical and style.RenderWithTitle already do
// all of it, and the panes stay plain tea.Models.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/helpbar"
	"github.com/anotherhadi/ilovetui/style"
)

// sidebarWidth is the sidebar's total width, border included.
const sidebarWidth = 24

// HelpProvider is the only contract in this pattern, and it's optional: a
// pane that implements it gets its own bindings listed in the global help
// bar while it's focused. A pane that doesn't just contributes nothing.
type HelpProvider interface {
	HelpBindings() []key.Binding
}

// ---------------------------------------------------------------- shell keys

type keyMap struct {
	Focus key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Focus: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
		Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
	}
}

// ---------------------------------------------------------------- the shell

type model struct {
	sidebar list.Model
	content tea.Model
	help    helpbar.Model
	keys    keyMap

	contentFocused bool
	w, h           int
}

func newModel() model {
	items := []list.Item{page("Overview"), page("Metrics"), page("Settings")}
	sidebar := bubbles.NewList(items, sidebarWidth-2, 0)
	sidebar.SetShowTitle(false)
	sidebar.SetShowStatusBar(false)
	sidebar.SetShowHelp(false)

	keys := defaultKeyMap()
	return model{
		sidebar: sidebar,
		content: newCounter(),
		help:    helpbar.New(helpbar.WithToggle(keys.Help), helpbar.WithGlobal(keys.Focus, keys.Quit)),
		keys:    keys,
	}
}

func (m model) Init() tea.Cmd { return m.content.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		return m.resize()

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			// Expanding the bar takes rows away from the panes.
			m.help, _ = m.help.Update(msg)
			return m.resize()
		case key.Matches(msg, m.keys.Focus):
			m.contentFocused = !m.contentFocused
			return m, nil
		}
		// Only the focused pane sees key presses.
		if m.contentFocused {
			var cmd tea.Cmd
			m.content, cmd = m.content.Update(msg)
			return m, cmd
		}
		var cmd tea.Cmd
		m.sidebar, cmd = m.sidebar.Update(msg)
		return m, cmd
	}

	// Everything else (ticks, HTTP responses...) goes to both, so a blurred
	// pane keeps working.
	var sidebarCmd, contentCmd tea.Cmd
	m.sidebar, sidebarCmd = m.sidebar.Update(msg)
	m.content, contentCmd = m.content.Update(msg)
	return m, tea.Batch(sidebarCmd, contentCmd)
}

// bodyHeight is the height left for the two panes once the help bar has
// taken its share. resize and View both go through it so they can't disagree.
func (m model) bodyHeight() int {
	return max(m.h-m.help.Height(m.focusedHelp()...), 0)
}

func (m model) resize() (model, tea.Cmd) {
	if m.w <= 0 || m.h <= 0 {
		return m, nil
	}
	m.help.SetWidth(m.w)
	inner := style.ContentHeight(m.bodyHeight())

	m.sidebar.SetSize(sidebarWidth-2, inner)

	var cmd tea.Cmd
	m.content, cmd = m.content.Update(tea.WindowSizeMsg{
		Width: max(m.w-sidebarWidth-2, 0), Height: inner,
	})
	return m, cmd
}

// focusedHelp is the focused pane's own bindings, if it offers any.
func (m model) focusedHelp() []key.Binding {
	if m.contentFocused {
		if hp, ok := m.content.(HelpProvider); ok {
			return hp.HelpBindings()
		}
		return nil
	}
	return []key.Binding{m.sidebar.KeyMap.CursorUp, m.sidebar.KeyMap.CursorDown}
}

func (m model) View() tea.View {
	if m.w <= 0 || m.h <= 0 {
		return tea.NewView("")
	}
	helpBar := m.help.View(m.focusedHelp()...)
	bodyH := m.bodyHeight()

	left := style.RenderWithTitle(
		panel(!m.contentFocused), "Menu", m.sidebar.View(), sidebarWidth, bodyH)
	right := style.RenderWithTitle(
		panel(m.contentFocused), "Content", m.content.View().Content, m.w-sidebarWidth, bodyH)

	view := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
		helpBar,
	))
	view.AltScreen = true
	return view
}

// panel picks the bordered panel style matching a pane's focus state.
func panel(focused bool) lipgloss.Style {
	if focused {
		return style.S.PanelFocused
	}
	return style.S.Panel
}

// ------------------------------------------------------- the right-hand pane

// counter is an ordinary tea.Model - nothing about it knows it's living in a
// shell. It implements HelpProvider purely to appear in the help bar.
type counter struct {
	n    int
	w, h int
	keys struct{ Inc, Dec key.Binding }
}

func newCounter() *counter {
	c := &counter{}
	c.keys.Inc = key.NewBinding(key.WithKeys("+", "k"), key.WithHelp("+/k", "increment"))
	c.keys.Dec = key.NewBinding(key.WithKeys("-", "j"), key.WithHelp("-/j", "decrement"))
	return c
}

func (c *counter) Init() tea.Cmd { return nil }

func (c *counter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.w, c.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, c.keys.Inc):
			c.n++
		case key.Matches(msg, c.keys.Dec):
			c.n--
		}
	}
	return c, nil
}

func (c *counter) View() tea.View {
	return tea.NewView(lipgloss.NewStyle().
		Width(c.w).Height(c.h).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(fmt.Sprintf("count: %d", c.n)))
}

func (c *counter) HelpBindings() []key.Binding {
	return []key.Binding{c.keys.Inc, c.keys.Dec}
}

// ------------------------------------------------------------- sidebar items

type page string

func (p page) Title() string       { return string(p) }
func (p page) Description() string { return "" }
func (p page) FilterValue() string { return string(p) }

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
