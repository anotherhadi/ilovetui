package layout

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// stubPane is a minimal Pane used across the test suite: it records every
// message it receives (and counts Focus/Blur specifically) and can be told
// to return a fixed cmd on its next Update, so tests can assert on both
// sides of the layout <-> pane contract without a real component.
type stubPane struct {
	focusN, blurN int
	w, h          int
	msgs          []tea.Msg
	nextCmd       tea.Cmd
	help          []key.Binding
}

func newStub() *stubPane { return &stubPane{} }

func (p *stubPane) Init() tea.Cmd { return nil }

func (p *stubPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	p.msgs = append(p.msgs, msg)
	switch m := msg.(type) {
	case FocusMsg:
		p.focusN++
	case BlurMsg:
		p.blurN++
	case SizeMsg:
		p.w, p.h = m.Width, m.Height
	}
	cmd := p.nextCmd
	p.nextCmd = nil
	return p, cmd
}

func (p *stubPane) View() string { return "" }

// HelpBindings implements HelpProvider.
func (p *stubPane) HelpBindings() []key.Binding { return p.help }

func (p *stubPane) last() tea.Msg {
	if len(p.msgs) == 0 {
		return nil
	}
	return p.msgs[len(p.msgs)-1]
}
