package bubbles

import (
	"charm.land/bubbles/v2/list"

	"github.com/anotherhadi/ilovetui/style"
)

func NewList(items []list.Item, width, height int) list.Model {
	m := list.New(items, NewDefaultDelegate(), width, height)
	m.Styles = themedListStyles()
	return m
}

func themedListStyles() list.Styles {
	s := list.DefaultStyles(true)
	s.Title = s.Title.Background(style.S.Primary).Foreground(style.S.Background)
	s.Spinner = s.Spinner.Foreground(style.S.Primary)
	s.Filter = themedTextInputStyles()
	s.DefaultFilterCharacterMatch = s.DefaultFilterCharacterMatch.Foreground(style.S.Primary)
	s.StatusBar = s.StatusBar.Foreground(style.S.Muted)
	s.StatusEmpty = s.StatusEmpty.Foreground(style.S.Subtle)
	s.StatusBarActiveFilter = s.StatusBarActiveFilter.Foreground(style.S.Text)
	s.StatusBarFilterCount = s.StatusBarFilterCount.Foreground(style.S.Subtle)
	s.NoItems = s.NoItems.Foreground(style.S.Subtle)
	s.ArabicPagination = s.ArabicPagination.Foreground(style.S.Subtle)
	s.ActivePaginationDot = s.ActivePaginationDot.Foreground(style.S.Primary)
	s.InactivePaginationDot = s.InactivePaginationDot.Foreground(style.S.Subtle)
	s.DividerDot = s.DividerDot.Foreground(style.S.Subtle)
	return s
}

func NewDefaultDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(style.S.Text)
	d.Styles.NormalDesc = d.Styles.NormalDesc.Foreground(style.S.Muted)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		BorderForeground(style.S.Primary).
		Foreground(style.S.Primary)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		BorderForeground(style.S.Primary).
		Foreground(style.S.Primary)
	d.Styles.DimmedTitle = d.Styles.DimmedTitle.Foreground(style.S.Subtle)
	d.Styles.DimmedDesc = d.Styles.DimmedDesc.Foreground(style.S.SubtleBg)
	d.Styles.FilterMatch = d.Styles.FilterMatch.Foreground(style.S.Primary)
	return d
}
