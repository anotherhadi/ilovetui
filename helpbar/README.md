# Helpbar

A responsive help bar: one line of key bindings that expands, on demand, into a multi-column view reflowed to use as many columns as the available width allows.

## Quick start

```go
import "github.com/anotherhadi/ilovetui/helpbar"

type model struct {
	help helpbar.Model
	keys keyMap
	h    int
}

func newModel() model {
	keys := defaultKeyMap()
	return model{
		help: helpbar.New(
			helpbar.WithToggle(keys.Help),             // '?' expands/collapses
			helpbar.WithGlobal(keys.Focus, keys.Quit), // always shown
		),
		keys: keys,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.h = msg.Height
		m.help.SetWidth(msg.Width)
	case tea.KeyPressMsg:
		m.help, _ = m.help.Update(msg) // flips ShowAll on the toggle binding
	}
	return m, nil
}

func (m model) View() tea.View {
	bar := m.help.View(m.focused().HelpBindings()...)
	body := renderBody(m.h - lipgloss.Height(bar))
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, body, bar))
}
```

## Global vs. contextual bindings

Bindings come from two places, and they render in this order:

1. The `WithToggle` binding, so the way to expand the bar always leads.
2. `WithGlobal` bindings - what your app reserves for itself (quit, switch pane...), set once.
3. Contextual bindings, passed to `View` at render time.

Contextual bindings are meant to change from render to render, so the bar can track whatever component currently has focus. Nothing in the package knows what "focus" means for your app - you just hand it a different slice:

```go
func (m model) focusedHelp() []key.Binding {
	if m.contentFocused {
		return m.content.HelpBindings()
	}
	return []key.Binding{m.sidebar.KeyMap.CursorUp, m.sidebar.KeyMap.CursorDown}
}
```

Disabled bindings (`key.Binding.SetEnabled(false)`) are dropped before layout, so the reflow never
budgets width for something that won't be drawn.

## Reserving room for the bar

The bar's height depends on its content and on whether it's expanded, so ask it rather than assuming
a fixed number of rows:

```go
body := m.height - m.help.Height(contextual...)
```

`Height` is exactly `lipgloss.Height` of what `View` returns for the same arguments, so the two can
never disagree about where the bar begins. An empty bar (no bindings, or no width set) renders `""`
and takes zero rows.

## Styling

Defaults come from the shared `style` theme. Override with `WithStyles`:

```go
helpbar.New(helpbar.WithStyles(myStyles)) // help.Styles from charm.land/bubbles/v2/help
```

## Examples

- `examples/sidebar` uses it as an app-wide bar tracking the focused pane.
- `examples/app` full app with a single help bar
