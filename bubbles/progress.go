package bubbles

import (
	"charm.land/bubbles/v2/progress"

	"github.com/anotherhadi/ilovetui/style"
)

// NewProgress returns a progress.Model with a themed fill, blending from
// style.S.Subtle to style.S.Primary. Any opts are forwarded to progress.New;
// pass progress.WithColors to override the default blend.
func NewProgress(opts ...progress.Option) progress.Model {
	allOpts := append([]progress.Option{
		progress.WithColors(style.S.Subtle, style.S.Primary),
	}, opts...)
	p := progress.New(allOpts...)
	p.EmptyColor = style.S.SubtleBg
	return p
}
