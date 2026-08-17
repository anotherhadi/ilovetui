// Command fullapp demonstrates a more realistic app shell: the sidebar
// leaf isn't a plain placeholder pane, it's a full layout.Model of its own
// (sidebar.New(), a themed list, see sidebar/sidebar.go) embedded exactly
// like the "router" one below - layout.Model implements Pane, so any
// package can hand back one to be wrapped in a Leaf, no adapter needed.
package first

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/layout"
)

type pane struct {
	id      string
	w, h    int
	focused bool
}

func newPane(id string) *pane { return &pane{id: id} }

func (p *pane) Init() tea.Cmd { return nil }

func (p *pane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.w, p.h = msg.Width, msg.Height
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	}
	return p, nil
}

func (p *pane) View() string {
	content := lipgloss.NewStyle().
		Width(p.w - 2).Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(p.id)
	return layout.Bordered(p.focused, p.w, p.h, content)
}

// HelpBindings implements layout.HelpProvider, so every pane - at the top
// level or nested three levels deep, doesn't matter - shows up correctly
// in the single, outer help bar.
func (p *pane) HelpBindings() []key.Binding {
	return []key.Binding{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "use "+p.id))}
}

// editorPane keeps a counter (press +) that's otherwise irrelevant to the
// app but proves the point of router.Model owning pages permanently:
// switch to "Second" and back, and it's still there. A plain pane (like
// terminal below) would just get reconstructed from scratch every time if
// something recreated it on each switch instead of keeping it alive.
type editorPane struct {
	w, h    int
	focused bool
	count   int
}

func newEditorPane() *editorPane { return &editorPane{} }

func (p *editorPane) Init() tea.Cmd { return nil }

func (p *editorPane) Update(msg tea.Msg) (layout.Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		p.w, p.h = msg.Width, msg.Height
	case layout.FocusMsg:
		p.focused = true
	case layout.BlurMsg:
		p.focused = false
	case tea.KeyPressMsg:
		if msg.String() == "+" {
			p.count++
		}
	}
	return p, nil
}

func (p *editorPane) View() string {
	content := lipgloss.NewStyle().
		Width(p.w - 2).Height(p.h - 2).
		AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(fmt.Sprintf("editor\n\npress + to increment: %d", p.count))
	return layout.Bordered(p.focused, p.w, p.h, content)
}

// HelpBindings implements layout.HelpProvider.
func (p *editorPane) HelpBindings() []key.Binding {
	return []key.Binding{key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "increment"))}
}

// NewWorkspace builds the "First" page's own nested layout.Model: note the
// lack of layout.AsRoot() here - only the outermost Model (see main) should
// render a help bar, or focused help would show up twice. second.
// NewWorkspace and sidebar.New() below are built the same way, same reason.
//
// This Model is built exactly once, by router.New(), and kept alive for
// the whole app's lifetime - router.Model just changes which page is
// rendered/routed to, it never reconstructs one. That's what lets
// editorPane's counter above survive switching to "Second" and back; see
// package router's own doc comment for why that ownership lives there
// and not in the sidebar.
func NewWorkspace() layout.Model {
	root := layout.VSplit(0.7,
		layout.Leaf("editor", newEditorPane()),
		layout.Leaf("terminal", newPane("terminal")),
	)
	return layout.New(root)
}
