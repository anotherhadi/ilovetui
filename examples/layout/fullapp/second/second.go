// Package second is the "Second" page's own workspace: deliberately a
// different shape from first's (a single pane holding a themed list, not an
// editor/terminal split), so swapping between the two via sidebar.Model's
// "enter" case is visibly a real change of sub-app, not just a relabeled
// copy of first. It also demonstrates notification.Show being called from
// deep inside a nested layout.Model, with nothing wiring it to the
// notification.Model that actually renders it - see main.go's model.View.
package second

import (
	"charm.land/bubbles/v2/key"
	bubbleslist "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/layout"
	"github.com/anotherhadi/ilovetui/notification"
)

// kind is a list entry: a notification.Kind plus the title/message Show
// gets called with when it's selected.
type kind struct {
	title   string
	message string
	kind    notification.Kind
}

func (k kind) Title() string       { return k.title }
func (k kind) Description() string { return "" }
func (k kind) FilterValue() string { return k.title }

type pane struct {
	list    bubbleslist.Model
	w, h    int
	focused bool
}

func newPane() *pane {
	items := []bubbleslist.Item{
		kind{title: "Info", message: "Just so you know.", kind: notification.Info},
		kind{title: "Success", message: "Config written to disk.", kind: notification.Success},
		kind{title: "Warning", message: "Disk space getting low on /dev/sda1.", kind: notification.Warning},
		kind{title: "Error", message: "Failed to reach the remote host.", kind: notification.Error},
	}
	list := bubbles.NewList(items, 0, 0)
	list.Title = "Notify"
	// Same reasoning as sidebar.Model: the built-in help footer is
	// redundant with layout's own centralized help bar (see HelpBindings
	// below).
	list.SetShowHelp(false)
	return &pane{list: list}
}

func (p *pane) Init() tea.Cmd { return nil }

func (p *pane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.w, p.h = msg.Width, msg.Height
		p.list.SetSize(p.w-2, p.h-2)
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	}

	if !p.focused {
		return p, nil
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
		if selected, ok := p.list.SelectedItem().(kind); ok {
			return p, notification.Show(selected.title, selected.message, selected.kind)
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p *pane) View() string {
	return layout.Bordered(p.focused, p.w, p.h, p.list.View())
}

// HelpBindings implements layout.HelpProvider.
func (p *pane) HelpBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "show notification")),
		key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "go down")),
		key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "go up")),
	}
}

// NewWorkspace builds the "Second" page's own nested layout.Model. No
// layout.AsRoot() - see first.NewWorkspace's doc comment.
func NewWorkspace() layout.Model {
	return layout.New(layout.Leaf("content", newPane()))
}
