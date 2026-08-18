# helpbar

A responsive help bar: one line of key bindings that expands, on toggle, into a multi-column view
reflowed to use as many columns as the available width allows.

- Bindings render in order: the `WithToggle` binding first, then `WithGlobal` bindings (set once),
  then contextual bindings passed to `View` at render time.
- Disabled bindings (`key.Binding.SetEnabled(false)`) are dropped before layout, so the reflow never
  budgets width for something that won't be drawn.
- `Height(contextual...)` always matches `lipgloss.Height` of `View` for the same arguments, so a
  host can reserve exactly the right amount of space. An empty bar renders `""` and takes 0 rows.

See `examples/helpbar`.
