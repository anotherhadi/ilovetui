// Command fullapp demonstrates a more realistic app shell: the sidebar
// leaf isn't a plain placeholder pane, it's a full layout.Model of its own
// (sidebar.New(), a themed list, see sidebar/sidebar.go) embedded exactly
// like the "router" one below - layout.Model implements Pane, so any
// package can hand back one to be wrapped in a Leaf, no adapter needed.
//
// It also owns the app's single notification.Model (see second/second.go,
// which triggers toasts via notification.Show without ever holding a
// reference to it) and composites its toasts over the whole rendered
// layout, top-right - notification has no dependency on layout, so this is
// the same Render(background string) pattern as examples/notification, just
// with layout.Model.View() as the background instead of a plain string.
//
// Same story for the app's single modal.Model (see third/third.go), except
// modal is composited last, on top of notif's own output: a modal is meant
// to command the whole screen's attention, so it should flatten any toast
// already showing to the same dim gray as everything else, not leave it
// floating on top in full color.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/router"
	"github.com/anotherhadi/ilovetui/examples/layout/fullapp/sidebar"
	"github.com/anotherhadi/ilovetui/layout"
	"github.com/anotherhadi/ilovetui/modal"
	"github.com/anotherhadi/ilovetui/notification"
)

// model is the actual top-level tea.Model: layout itself reserves no quit
// key (that's an app policy, not layout's to make), so the host wraps it
// and handles ctrl+c/q itself, same as any other custom component in this
// repo (see examples/tabs).
type model struct {
	layout layout.Model
	notif  notification.Model
	modal  modal.Model
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.layout.Init(), m.notif.Init(), m.modal.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	updated, layoutCmd := m.layout.Update(msg)
	m.layout = updated.(layout.Model)

	var notifCmd, modalCmd tea.Cmd
	m.notif, notifCmd = m.notif.Update(msg)
	m.modal, modalCmd = m.modal.Update(msg)

	return m, tea.Batch(layoutCmd, notifCmd, modalCmd)
}

func (m model) View() tea.View {
	view := tea.NewView(m.modal.Render(m.notif.Render(m.layout.View())))
	view.AltScreen = true
	return view
}

func main() {
	root := layout.HSplit(0.25,
		layout.Leaf("sidebar", sidebar.New()),
		layout.Leaf("router", router.New()),
	).WithMaximum(20)

	m := model{
		layout: layout.New(root, layout.AsRoot()),
		notif:  notification.New(notification.WithPosition(notification.TopRight)),
		modal:  modal.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
