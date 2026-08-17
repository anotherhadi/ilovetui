// Package third is the "Third" page's own workspace: a single pane showing
// modal.Show being called from deep inside a nested layout.Model, exactly
// like second/second.go does for notification.Show - modal has no
// dependency on layout either, so nothing here needs a reference to the
// modal.Model that actually renders it (see main.go's model.View).
package third

import (
	"charm.land/bubbles/v2/key"
	bubbleslist "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
	"github.com/anotherhadi/ilovetui/layout"
	"github.com/anotherhadi/ilovetui/modal"
)

// action is a list entry: the title/content Show gets called with when it's
// selected, and whether closing it needs an explicit "y" (a destructive
// confirm) or any key at all (a plain info popup).
type action struct {
	title, content string
	confirm        bool
}

func (a action) Title() string       { return a.title }
func (a action) Description() string { return "" }
func (a action) FilterValue() string { return a.title }

type pane struct {
	list    bubbleslist.Model
	w, h    int
	focused bool
	// open/confirm track the modal this pane itself opened, so it knows to
	// swallow keys instead of forwarding them to the list underneath while
	// it's up - modal.Model has no notion of "focus" of its own, the pane
	// that triggered it is responsible for gating input while it's open.
	open    bool
	confirm bool
}

func newPane() *pane {
	items := []bubbleslist.Item{
		action{title: "Delete file", content: "This can't be undone.\n\ny: confirm  esc: cancel", confirm: true},
		action{title: "About", content: "modal renders a centered popup over a\nflat-dimmed background.\n\nesc: close"},
	}
	list := bubbles.NewList(items, 0, 0)
	list.Title = "Modal"
	// Same reasoning as sidebar/second: redundant with layout's own
	// centralized help bar (see HelpBindings below).
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

	if key, ok := msg.(tea.KeyPressMsg); ok {
		if p.open {
			if !p.confirm || key.String() == "y" || key.String() == "esc" {
				p.open = false
				return p, modal.Close()
			}
			return p, nil
		}

		if key.String() == "enter" {
			if selected, ok := p.list.SelectedItem().(action); ok {
				p.open, p.confirm = true, selected.confirm
				return p, modal.Show(selected.title+"?", selected.content)
			}
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
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open modal")),
		key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "go down")),
		key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "go up")),
	}
}

// NewWorkspace builds the "Third" page's own nested layout.Model. No
// layout.AsRoot() - see first.NewWorkspace's doc comment.
func NewWorkspace() layout.Model {
	return layout.New(layout.Leaf("content", newPane()))
}
