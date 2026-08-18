package bubbles

import (
	"charm.land/bubbles/v2/progress"

	"github.com/anotherhadi/ilovetui/style"
)

func NewProgress(opts ...progress.Option) progress.Model {
	allOpts := append([]progress.Option{
		progress.WithColors(style.S.Subtle, style.S.Primary),
	}, opts...)
	p := progress.New(allOpts...)
	p.EmptyColor = style.S.SubtleBg
	return p
}
