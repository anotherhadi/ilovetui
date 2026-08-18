package bubbles

import (
	"charm.land/bubbles/v2/filepicker"

	"github.com/anotherhadi/ilovetui/style"
)

func NewFilePicker() filepicker.Model {
	f := filepicker.New()
	s := filepicker.DefaultStyles()
	s.Cursor = s.Cursor.Foreground(style.S.Primary)
	s.DisabledCursor = s.DisabledCursor.Foreground(style.S.Subtle)
	s.Symlink = s.Symlink.Foreground(style.S.Warning)
	s.Directory = s.Directory.Foreground(style.S.Primary)
	s.File = s.File.Foreground(style.S.Text)
	s.DisabledFile = s.DisabledFile.Foreground(style.S.Subtle)
	s.Permission = s.Permission.Foreground(style.S.Muted)
	s.Selected = s.Selected.Foreground(style.S.Primary)
	s.DisabledSelected = s.DisabledSelected.Foreground(style.S.Subtle)
	s.FileSize = s.FileSize.Foreground(style.S.Muted)
	s.EmptyDirectory = s.EmptyDirectory.Foreground(style.S.Subtle)
	f.Styles = s
	return f
}
