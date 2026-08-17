package layout

// Direction is the axis a Split divides its space along.
type Direction int

const (
	// Horizontal places the first child on the left, the second on the right.
	Horizontal Direction = iota
	// Vertical stacks the first child on top of the second.
	Vertical
)

// Node is one slot in the layout tree: either a Leaf (holding a Pane) or a
// Split (dividing its space between two child Nodes).
// Build a tree with Leaf and Split/HSplit/VSplit, then hand the root to New.
//
// Node is pointer-based on purpose: panes need a stable identity across
// arbitrarily nested splits, so there's no flat index to keep in sync the
// way a slice-based component would.
type Node struct {
	id string // "" means unaddressable; always non-empty for a Leaf

	leaf  bool
	model Pane // set when leaf

	dir           Direction // set when !leaf
	ratio         float64   // proportion of space given to first; set when !leaf
	min, max      int       // clamp on first's resolved cell size, in cells; 0 = unset
	first, second *Node     // set when !leaf
}

// Leaf wraps a single pane. id must be non-empty and unique within the tree
// it ends up in: it's how the pane is targeted later by SendMsg,
// RequestFocusMsg, SplitLeaf, CloseLeaf and Resize.
func Leaf(id string, model Pane) *Node {
	return &Node{id: id, leaf: true, model: model}
}

// Split divides its space between first and second along dir, giving ratio
// (0 to 1) of it to first and the rest to second. id may be "" if the split
// itself never needs to be addressed by Resize; it plays no role in pane
// addressing (only Leaf ids do).
func Split(id string, dir Direction, ratio float64, first, second *Node) *Node {
	return &Node{id: id, dir: dir, ratio: ratio, first: first, second: second}
}

// HSplit is Split with Horizontal and no id: layout.HSplit(0.3, left, right).
func HSplit(ratio float64, first, second *Node) *Node {
	return Split("", Horizontal, ratio, first, second)
}

// VSplit is Split with Vertical and no id: layout.VSplit(0.3, top, bottom).
func VSplit(ratio float64, first, second *Node) *Node {
	return Split("", Vertical, ratio, first, second)
}

// WithID sets the id used to address this node later (currently only
// meaningful on a Split, for Resize; a Leaf already gets its id from Leaf).
func (n *Node) WithID(id string) *Node {
	n.id = id
	return n
}

// WithMinimum clamps first's resolved size to never go below cells. No-op on
// a Leaf, which has no size of its own to constrain.
func (n *Node) WithMinimum(cells int) *Node {
	if !n.leaf {
		n.min = cells
	}
	return n
}

// WithMaximum clamps first's resolved size to never exceed cells. No-op on a
// Leaf. Setting both WithMinimum and WithMaximum to the same value fixes
// first's size regardless of ratio.
func (n *Node) WithMaximum(cells int) *Node {
	if !n.leaf {
		n.max = cells
	}
	return n
}

// findNode returns the node with the given id anywhere in the tree rooted at
// n, without mutating anything. id == "" never matches (it's the
// unaddressable sentinel, and several Splits may share it).
func findNode(n *Node, id string) (*Node, bool) {
	if n == nil || id == "" {
		return nil, false
	}
	if n.id == id {
		return n, true
	}
	if n.leaf {
		return nil, false
	}
	if found, ok := findNode(n.first, id); ok {
		return found, true
	}
	return findNode(n.second, id)
}

// splitConfig collects SplitOption values applied by SplitLeaf.
type splitConfig struct {
	id    string
	ratio float64
}

// SplitOption configures the Split node SplitLeaf creates in place of the
// leaf it splits.
type SplitOption func(*Node, *splitConfig)

// WithSplitID addresses the split SplitLeaf creates, so it can later be
// targeted by Resize.
func WithSplitID(id string) SplitOption {
	return func(_ *Node, c *splitConfig) { c.id = id }
}

// WithSplitRatio overrides SplitLeaf's default 50/50 split.
func WithSplitRatio(ratio float64) SplitOption {
	return func(_ *Node, c *splitConfig) { c.ratio = ratio }
}

// WithSplitMinimum clamps the size of the original (pre-split) leaf's side
// of the new split. Equivalent to calling (*Node).WithMinimum on the split
// SplitLeaf produces.
func WithSplitMinimum(cells int) SplitOption {
	return func(n *Node, _ *splitConfig) { n.WithMinimum(cells) }
}

// WithSplitMaximum clamps the size of the original (pre-split) leaf's side
// of the new split. Equivalent to calling (*Node).WithMaximum on the split
// SplitLeaf produces.
func WithSplitMaximum(cells int) SplitOption {
	return func(n *Node, _ *splitConfig) { n.WithMaximum(cells) }
}

// splitLeaf replaces the leaf identified by id with a Split holding the
// original leaf as first and a new Leaf(newID, newModel) as second. Returns
// the (possibly new) tree root and whether id was found and was a leaf.
func splitLeaf(root *Node, id string, dir Direction, newID string, newModel Pane, opts ...SplitOption) (*Node, bool) {
	target, ok := findNode(root, id)
	if !ok || !target.leaf {
		return root, false
	}

	cfg := splitConfig{ratio: 0.5}
	split := Split("", dir, 0.5, target, Leaf(newID, newModel))
	for _, opt := range opts {
		opt(split, &cfg)
	}
	split.id = cfg.id
	split.ratio = cfg.ratio

	newRoot, _ := replaceNode(root, id, func(*Node) *Node { return split })
	return newRoot, true
}

// replaceNode returns a copy of the tree rooted at n with the node
// identified by id swapped for transform's result. Only the path down to
// that node is cloned, everything else is shared. Returns n unchanged and
// false if id isn't found.
func replaceNode(n *Node, id string, transform func(*Node) *Node) (*Node, bool) {
	if n == nil || id == "" {
		return n, false
	}
	if n.id == id {
		return transform(n), true
	}
	if n.leaf {
		return n, false
	}
	if newFirst, ok := replaceNode(n.first, id, transform); ok {
		clone := *n
		clone.first = newFirst
		return &clone, true
	}
	if newSecond, ok := replaceNode(n.second, id, transform); ok {
		clone := *n
		clone.second = newSecond
		return &clone, true
	}
	return n, false
}

// closeLeaf removes the leaf identified by id, promoting its sibling to take
// the place of their parent Split. Returns the (possibly new) tree root and
// whether id was found as a direct child of some Split (the tree's own root
// leaf, with no parent, can never be closed this way).
func closeLeaf(root *Node, id string) (*Node, bool) {
	if root == nil || root.leaf || id == "" {
		return root, false
	}
	if root.first.leaf && root.first.id == id {
		return root.second, true
	}
	if root.second.leaf && root.second.id == id {
		return root.first, true
	}
	if newFirst, ok := closeLeaf(root.first, id); ok {
		clone := *root
		clone.first = newFirst
		return &clone, true
	}
	if newSecond, ok := closeLeaf(root.second, id); ok {
		clone := *root
		clone.second = newSecond
		return &clone, true
	}
	return root, false
}
