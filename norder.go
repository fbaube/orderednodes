package orderednodes

import (
	"io"
)

// Norder is satisfied by [*Nord] NOT by Nord.
type Norder interface {
	// SeqID() int
	// SetSeqID(int)
	// Level is zero-based (i.e. root nord's is 0) 
	Level() int
	// RelFP is rel.filepath for a file/dir, and for a DOM
	// node, it is meaningless, unless it is a [RootNorder],
	// for which it is the rel.path to the DOCUMENT
	RelFP() string
	// AbsFP is abs.filepath for a file/dir, and for 
	// a DOM node, the (abs.)path of the node w.r.t.
	// the document root, except for a [RootNorder],
	// for which it is the abs.path to the DOCUMENT
	AbsFP() string
	IsRoot() bool
	// Root should always return the root, at arena index 0 
	Root() RootNorder
	IsDir() bool
	// IsDirlike() bool // FIXME: add this!
	Parent() Norder
	HasKids() bool
	FirstKid() Norder
	LastKid() Norder
	PrevKid() Norder
	NextKid() Norder
	KidsAsSlice() []Norder
	// AddKid returns the kid, who knows his
	// own arena index (using [slices.Index])
	AddKid(Norder) Norder
	// AddKids returns the method target
	// - the parent of all the kids 
	AddKids([]Norder) Norder 
	ReplaceWith(Norder) Norder
	SetParent(Norder)
	SetPrevKid(Norder)
	SetNextKid(Norder)
	SetFirstKid(Norder)
	SetLastKid(Norder)
	LinePrefixString() string
	LineSummaryString() string
	LineSummaryFunc() StringFunc
	PrintTree(io.Writer) error
	// PACKAGE METHODS
	setLevel(int)
}
