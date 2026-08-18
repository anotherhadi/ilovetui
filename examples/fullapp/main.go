package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/drawer"
	"github.com/anotherhadi/ilovetui/examples/fullapp/metrics"
	"github.com/anotherhadi/ilovetui/examples/fullapp/overview"
	"github.com/anotherhadi/ilovetui/examples/fullapp/settings"
	"github.com/anotherhadi/ilovetui/examples/fullapp/sidebar"
	"github.com/anotherhadi/ilovetui/helpbar"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/notification"
	"github.com/anotherhadi/ilovetui/style"
)

var pages = []struct {
	item sidebar.NavItem
	new  func() tea.Model
}{
	{sidebar.NavItem{Icon: "󰋜", Name: "Overview"}, func() tea.Model { return overview.New() }},
	{sidebar.NavItem{Icon: "󰄨", Name: "Metrics"}, func() tea.Model { return metrics.New() }},
	{sidebar.NavItem{Icon: "󰒓", Name: "Settings"}, func() tea.Model { return settings.New() }},
}

type HelpProvider interface {
	HelpBindings() []key.Binding
}

func navItems() []sidebar.NavItem {
	items := make([]sidebar.NavItem, len(pages))
	for i, p := range pages {
		items[i] = p.item
	}
	return items
}

type keyMap struct {
	FocusSidebar key.Binding
	Close        key.Binding
	Help         key.Binding
	Quit         key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		FocusSidebar: key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("ctrl+b", "focus sidebar")),
		Close:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close")),
		Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:         key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q / ctrl+c", "quit")),
	}
}

type Model struct {
	sidebar sidebar.Model
	page    tea.Model
	help    helpbar.Model
	notif   notification.Model
	modal   modal.Model
	drawer  drawer.Model
	keys    keyMap

	sidebarWidth   int
	hideNav        bool
	contentFocused bool
	width, height  int
}

func NewModel() Model {
	keys := defaultKeyMap()
	return Model{
		sidebar:      sidebar.New(navItems()...),
		sidebarWidth: 24,
		page:         pages[0].new(),
		help:         helpbar.New(helpbar.WithToggle(keys.Help), helpbar.WithGlobal(keys.FocusSidebar, keys.Quit)),
		notif:        notification.New(),
		modal:        modal.New(),
		drawer:       drawer.New(),
		keys:         keys,
		hideNav:      true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.page.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m.resize()

	case sidebar.SelectMsg:
		m.page = pages[msg.Index].new()
		var cmd tea.Cmd
		m, cmd = m.resize()
		return m, tea.Batch(m.page.Init(), cmd)

	case sidebar.BlurMsg:
		m.contentFocused = true
		return m.resize()

	case tea.KeyPressMsg:
		if m.modal.Open() {
			if key.Matches(msg, m.keys.Close) {
				return m, modal.Close()
			}
			var cmd tea.Cmd
			m.modal, cmd = m.modal.Update(msg)
			return m, cmd
		}
		if m.drawer.Open() {
			if key.Matches(msg, m.keys.Close) {
				return m, drawer.Close()
			}
			var cmd tea.Cmd
			m.drawer, cmd = m.drawer.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help, _ = m.help.Update(msg)
			return m.resize()
		case key.Matches(msg, m.keys.FocusSidebar):
			m.contentFocused = !m.contentFocused
			return m.resize()
		}

		if m.contentFocused {
			var cmd tea.Cmd
			m.page, cmd = m.page.Update(msg)
			return m, cmd
		}
		var cmd tea.Cmd
		m.sidebar, cmd = m.sidebar.Update(msg)
		return m, cmd
	}

	var notifCmd, modalCmd, drawerCmd, sidebarCmd, pageCmd tea.Cmd
	m.notif, notifCmd = m.notif.Update(msg)
	m.modal, modalCmd = m.modal.Update(msg)
	m.drawer, drawerCmd = m.drawer.Update(msg)
	m.sidebar, sidebarCmd = m.sidebar.Update(msg)
	m.page, pageCmd = m.page.Update(msg)
	return m, tea.Batch(notifCmd, modalCmd, drawerCmd, sidebarCmd, pageCmd)
}

func (m Model) navWidth() int {
	if m.hideNav && m.contentFocused {
		return 0
	}
	return m.sidebarWidth
}

func (m Model) resize() (Model, tea.Cmd) {
	if m.width <= 0 || m.height <= 0 {
		return m, nil
	}
	m.help.SetWidth(m.width)
	inner := style.ContentHeight(m.height - m.help.Height(m.focusedHelp()...))

	navW := m.navWidth()
	m.sidebar.SetSize(max(navW-2, 0), inner)

	var cmd tea.Cmd
	m.page, cmd = m.page.Update(tea.WindowSizeMsg{
		Width:  max(m.width-navW-2, 0),
		Height: inner,
	})
	return m, cmd
}

func (m Model) focusedHelp() []key.Binding {
	if !m.contentFocused {
		return m.sidebar.HelpBindings()
	}
	if hp, ok := m.page.(HelpProvider); ok {
		return hp.HelpBindings()
	}
	return nil
}

func (m Model) View() tea.View {
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}
	helpBar := m.help.View(m.focusedHelp()...)
	bodyH := max(m.height-lipgloss.Height(helpBar), 0)

	navW := m.navWidth()
	body := style.RenderWithTitle(
		panel(m.contentFocused), m.sidebar.Selected().Name, m.page.View().Content, m.width-navW, bodyH)

	if navW > 0 {
		left := style.RenderWithTitle(
			panel(!m.contentFocused), "Menu", m.sidebar.View(), navW, bodyH)
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, body)
	}

	appView := lipgloss.JoinVertical(lipgloss.Left,
		body,
		helpBar,
	)

	view := tea.NewView(m.notif.Render(m.modal.Render(m.drawer.Render(appView))))
	view.AltScreen = true
	return view
}

func panel(focused bool) lipgloss.Style {
	if focused {
		return style.S.PanelFocused
	}
	return style.S.Panel
}

func main() {
	if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
