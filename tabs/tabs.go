package tabs

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/style"
)

type Tab interface {
	Init() tea.Cmd
	Update(tea.Msg) (Tab, tea.Cmd)
	View() string
}

type Item struct {
	Title string
	Model Tab
}

type KeyMap struct {
	Next key.Binding
	Prev key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next: key.NewBinding(
			key.WithKeys("right", "l", "tab"),
			key.WithHelp("→/tab", "next tab"),
		),
		Prev: key.NewBinding(
			key.WithKeys("left", "h", "shift+tab"),
			key.WithHelp("←/shift+tab", "previous tab"),
		),
	}
}

type Styles struct {
	// Border shape/padding for the tab boxes, color-less: the actual border
	// color is picked at render time from FocusedBorder/BlurredBorder
	// depending on Model.Focused(), so the whole tabs+content frame always
	// reads as one continuous, single-colored box.
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	// Title text styles: this is what actually distinguishes the active tab
	// from the others.
	ActiveTitle   lipgloss.Style
	InactiveTitle lipgloss.Style
	// Border shape/padding for the content pane, same color-less rule as
	// above.
	Content lipgloss.Style

	// The border family (rounded, normal, thick...) tabs and Content are
	// built from, snapshotted from style.S.BorderType at DefaultStyles()
	// time. Kept around so renderBar can pick the right per-position notch
	// glyph (corner vs. T-junction) for that same family at render time.
	BorderType lipgloss.Border

	FocusedBorder color.Color
	BlurredBorder color.Color
}

func DefaultStyles() Styles {
	bt := style.S.BorderType
	inactiveBorder, activeBorder := tabBorders(bt)

	return Styles{
		InactiveTab: lipgloss.NewStyle().
			Border(inactiveBorder, true).
			Padding(0, 1),
		ActiveTab: lipgloss.NewStyle().
			Border(activeBorder, true).
			Padding(0, 1),

		ActiveTitle:   lipgloss.NewStyle().Foreground(style.S.Primary).Bold(true),
		InactiveTitle: lipgloss.NewStyle().Foreground(style.S.Subtle),

		Content: lipgloss.NewStyle().
			Border(bt).
			UnsetBorderTop().
			Padding(1, 2),

		BorderType: bt,

		FocusedBorder: style.S.Primary,
		BlurredBorder: style.S.Subtle,
	}
}

// tabBorders derives the inactive/active tab border shapes from a border
// family, using its own junction glyphs (MiddleBottom, MiddleLeft,
// MiddleRight...) instead of hardcoded characters, so tabs follow
// style.S.BorderType instead of always looking rounded regardless of config.
//
// Inactive tabs get a plain "┴"-style bottom (a Content has UnsetBorderTop,
// so this line is what actually separates the bar from Content below).
// The active tab's bottom is left open (blank) with its corners swapped to
// the family's own BottomLeft/BottomRight glyphs, so its sides appear to
// flow straight down into Content.
func tabBorders(bt lipgloss.Border) (inactive, active lipgloss.Border) {
	inactive = bt
	inactive.BottomLeft = bt.MiddleBottom
	inactive.BottomRight = bt.MiddleBottom

	active = bt
	active.Bottom = " "
	active.BottomLeft = bt.BottomRight
	active.BottomRight = bt.BottomLeft

	return inactive, active
}

type Model struct {
	items   []Item
	active  int
	focused bool
	loop    bool
	width   int
	height  int

	styles Styles
	keyMap KeyMap
}

type Option func(*Model)

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

func WithKeyMap(k KeyMap) Option {
	return func(m *Model) {
		m.keyMap = k
	}
}

func WithActive(i int) Option {
	return func(m *Model) {
		m.active = i
	}
}

// WithLoop sets whether Next/Prev navigation wraps around: Next from the
// last tab goes to the first, Prev from the first tab goes to the last.
// On by default; pass false to clamp at either end instead.
func WithLoop(l bool) Option {
	return func(m *Model) {
		m.loop = l
	}
}

// WithFocus sets the initial focus state. A focused tabs bar renders its
// border in the accent color, a blurred one in the muted color, letting a
// host app with several panes show which one is currently active.
func WithFocus(f bool) Option {
	return func(m *Model) {
		m.focused = f
	}
}

func New(items []Item, opts ...Option) Model {
	m := Model{
		items:   items,
		focused: true,
		loop:    true,
		styles:  DefaultStyles(),
		keyMap:  DefaultKeyMap(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.active = clamp(m.active, 0, len(m.items)-1)

	return m
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.items))
	for i, item := range m.items {
		cmds[i] = item.Model.Init()
	}
	return tea.Batch(cmds...)
}

func (m Model) Active() int {
	return m.active
}

func (m *Model) SetActive(i int) {
	m.active = clamp(i, 0, len(m.items)-1)
}

func (m Model) Items() []Item {
	return m.items
}

func (m Model) ActiveItem() Item {
	return m.items[m.active]
}

func (m Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) Loop() bool {
	return m.loop
}

func (m *Model) SetLoop(l bool) {
	m.loop = l
}

func (m Model) Width() int {
	return m.width
}

// SetWidth sets the target outer width for the whole component. The tab bar
// itself always keeps its intrinsic width (the sum of its tab labels);
// Content stretches to Width if that's wider, so a host building a
// fullscreen app can make the content pane fill the terminal without the tab
// bar itself looking artificially stretched.
func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m Model) Height() int {
	return m.height
}

// SetHeight sets the target outer height for the whole component. Content's
// height is Height minus the bar's own (fixed) height; SetHeight is a no-op
// on the render until called, so the zero-value Model keeps auto-sizing
// Content to whatever the active item's View() returns.
func (m *Model) SetHeight(h int) {
	m.height = h
}

// SetSize is a shorthand for SetWidth followed by SetHeight.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// ContentWidth and ContentHeight report the usable inner area available to
// the active item's View() once Width/Height are set: Content's box size
// minus its own border and padding. tabs has no generic way to size an
// arbitrary Tab itself (the interface is intentionally minimal), so a host
// building a fullscreen app calls these after SetSize and forwards the
// result to its own Tab implementations, exactly as it would size any other
// nested bubbles component. ContentHeight returns 0 until SetHeight has been
// called (see SetHeight).
func (m Model) ContentWidth() int {
	if len(m.items) == 0 {
		return 0
	}
	w := lipgloss.Width(m.renderBar(m.collapsedSegments(m.width), m.styles.BlurredBorder, false))
	if m.width > w {
		w = m.width
	}
	if inner := w - m.styles.Content.GetHorizontalFrameSize(); inner > 0 {
		return inner
	}
	return 0
}

func (m Model) ContentHeight() int {
	if len(m.items) == 0 || m.height <= 0 {
		return 0
	}
	barHeight := lipgloss.Height(m.renderBar(m.collapsedSegments(m.width), m.styles.BlurredBorder, false))
	if inner := m.height - barHeight - m.styles.Content.GetVerticalFrameSize(); inner > 0 {
		return inner
	}
	return 0
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(keyMsg, m.keyMap.Next):
			m.active = m.step(1)
			return m, nil
		case key.Matches(keyMsg, m.keyMap.Prev):
			m.active = m.step(-1)
			return m, nil
		}
	}

	if len(m.items) == 0 {
		return m, nil
	}

	var cmd tea.Cmd
	m.items[m.active].Model, cmd = m.items[m.active].Model.Update(msg)
	return m, cmd
}

// tabSegment is a single box drawn on the bar: either a real item, or the
// synthetic "+N" segment standing in for tabs collapsed by collapsedSegments.
// It's never active and never has a backing Item.
type tabSegment struct {
	title    string
	isActive bool
	isMore   bool
}

func (m Model) segments() []tabSegment {
	segs := make([]tabSegment, len(m.items))
	for i, item := range m.items {
		segs[i] = tabSegment{title: item.Title, isActive: i == m.active}
	}
	return segs
}

// segmentWidth measures a segment as it would actually render, without
// needing a border color (color doesn't affect measured width).
func (m Model) segmentWidth(seg tabSegment) int {
	tabStyle := m.styles.InactiveTab
	titleStyle := m.styles.InactiveTitle
	if seg.isActive {
		tabStyle = m.styles.ActiveTab
		titleStyle = m.styles.ActiveTitle
	}
	return lipgloss.Width(tabStyle.Render(titleStyle.Render(seg.title)))
}

// collapsedSegments returns the full segment list unchanged if it already
// fits within budget (or budget is unset). Otherwise it keeps a contiguous
// window of tabs that always includes the active one - grown outward from
// active, alternating backward/forward, as far as it fits - and folds
// everything left out of that window into a single trailing "+N" segment.
func (m Model) collapsedSegments(budget int) []tabSegment {
	segs := m.segments()
	if budget <= 0 {
		return segs
	}

	widths := make([]int, len(segs))
	total := 0
	for i, seg := range segs {
		widths[i] = m.segmentWidth(seg)
		total += widths[i]
	}
	if total <= budget {
		return segs
	}

	moreWidth := m.segmentWidth(tabSegment{title: fmt.Sprintf("+%d", len(segs)-1), isMore: true})
	fitBudget := budget - moreWidth
	if fitBudget < widths[m.active] {
		// Not even room for active + the badge: guarantee active alone
		// fits, even if that leaves the badge slightly cramped.
		fitBudget = widths[m.active]
	}

	start, end := m.active, m.active
	used := widths[m.active]
	for {
		grew := false
		if start > 0 && used+widths[start-1] <= fitBudget {
			start--
			used += widths[start]
			grew = true
		}
		if end < len(segs)-1 && used+widths[end+1] <= fitBudget {
			end++
			used += widths[end]
			grew = true
		}
		if !grew {
			break
		}
	}

	hidden := len(segs) - (end - start + 1)
	if hidden <= 0 {
		return segs
	}

	visible := append([]tabSegment{}, segs[start:end+1]...)
	visible = append(visible, tabSegment{title: fmt.Sprintf("+%d", hidden), isMore: true})
	return visible
}

func (m Model) View() string {
	if len(m.items) == 0 {
		return ""
	}

	borderColor := m.styles.BlurredBorder
	if m.focused {
		borderColor = m.styles.FocusedBorder
	}

	segs := m.collapsedSegments(m.width)

	bar := m.renderBar(segs, borderColor, false)
	contentWidth := lipgloss.Width(bar)

	if m.width > contentWidth {
		// Re-render with the last tab's right edge treated as an interior
		// junction instead of the widget's outer edge, since the cap line
		// now continues past it into the extension.
		bar = m.extendBarCap(m.renderBar(segs, borderColor, true), m.width, borderColor)
		contentWidth = m.width
	}
	contentStyle := m.styles.Content.
		BorderForeground(borderColor).
		Width(contentWidth)
	if m.height > 0 {
		if h := m.height - lipgloss.Height(bar); h > 0 {
			contentStyle = contentStyle.Height(h)
		}
	}

	content := contentStyle.Render(m.items[m.active].Model.View())

	return lipgloss.JoinVertical(lipgloss.Left, bar, content)
}

// renderBar builds the tab bar. extendCap should be true when the caller
// already knows the cap line will be stretched past the last tab (see
// extendBarCap): in that case the last tab's right edge is drawn as an
// interior junction (bt.MiddleLeft) rather than the widget's outer edge,
// since the horizontal line continues past it instead of terminating there.
func (m Model) renderBar(segs []tabSegment, borderColor color.Color, extendCap bool) string {
	rendered := make([]string, len(segs))

	for i, seg := range segs {
		isFirst, isLast := i == 0, i == len(segs)-1

		tabStyle := m.styles.InactiveTab
		titleStyle := m.styles.InactiveTitle
		if seg.isActive {
			tabStyle = m.styles.ActiveTab
			titleStyle = m.styles.ActiveTitle
		}
		tabStyle = tabStyle.BorderForeground(borderColor)

		bt := m.styles.BorderType
		border, _, _, _, _ := tabStyle.GetBorder()
		switch {
		case isFirst && seg.isActive:
			border.BottomLeft = bt.Left
		case isFirst && !seg.isActive:
			border.BottomLeft = bt.MiddleLeft
		case isLast && seg.isActive && !extendCap:
			border.BottomRight = bt.Right
		case isLast && !seg.isActive && !extendCap:
			border.BottomRight = bt.MiddleRight
		}
		// extendCap: no BottomRight override at all, so the last tab falls
		// back to its type's plain default (already set in DefaultStyles:
		// bt.MiddleBottom for inactive, the swap-trick corner for active) -
		// same as every other, non-edge tab. The isLast-specific corners
		// above only make sense when this really is the widget's edge and
		// Content's own border aligns right below it; once the cap extends
		// past it, that's no longer true.
		tabStyle = tabStyle.Border(border)

		rendered[i] = tabStyle.Render(titleStyle.Render(seg.title))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

// extendBarCap stretches only the bar's bottom row out to width. That row
// doubles as Content's own top border (Content has UnsetBorderTop), so when
// Content is wider than the bar's natural width, it needs to reach all the
// way across or the frame looks broken open above the extra space.
// lipgloss.JoinVertical would otherwise pad the shorter bar rows with plain
// spaces, not border characters.
func (m Model) extendBarCap(bar string, width int, borderColor color.Color) string {
	gap := width - lipgloss.Width(bar)
	if gap <= 0 {
		return bar
	}

	bt := m.styles.BorderType
	fill := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat(bt.Bottom, gap-1) + bt.TopRight)

	lines := strings.Split(bar, "\n")
	lines[len(lines)-1] += fill

	return strings.Join(lines, "\n")
}

// step moves the active index by delta (+1 for Next, -1 for Prev). With Loop
// it wraps around at either end; otherwise it just clamps, so Next on the
// last tab (or Prev on the first) is a no-op.
func (m Model) step(delta int) int {
	n := len(m.items)
	if n == 0 {
		return 0
	}
	if m.loop {
		return ((m.active+delta)%n + n) % n
	}
	return clamp(m.active+delta, 0, n-1)
}

func clamp(v, low, high int) int {
	if high < low {
		return low
	}
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}
