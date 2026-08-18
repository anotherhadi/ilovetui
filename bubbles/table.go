package bubbles

import (
	"charm.land/bubbles/v2/table"

	"github.com/anotherhadi/ilovetui/style"
)

func NewTable(opts ...table.Option) table.Model {
	t := table.New(opts...)
	s := table.DefaultStyles()
	s.Header = s.Header.Foreground(style.S.Primary)
	s.Cell = s.Cell.Foreground(style.S.Text)
	s.Selected = s.Selected.Foreground(style.S.Primary)
	t.SetStyles(s)
	return t
}
