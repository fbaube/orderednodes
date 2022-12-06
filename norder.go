package orderednodes

import (
	"io"
)

// Norder is satisfied by *Nord NOT by Nord
type Norder interface {
	// PUBLIC METHODS
	SeqID() int
	// SetSeqID(int)
	Level() int
	RelFP() string
	AbsFP() string
	IsRoot() bool
	GetRoot() Norder
	IsDir() bool
	// SetIsRoot(bool)
	Parent() Norder
	HasKids() bool
	FirstKid() Norder
	LastKid() Norder
	PrevKid() Norder
	NextKid() Norder
	KidsAsSlice() []Norder
	AddKid(Norder) Norder
	ReplaceWith(Norder) Norder
	SetParent(Norder)
	SetPrevKid(Norder)
	SetNextKid(Norder)
	SetFirstKid(Norder)
	SetLastKid(Norder)
	LinePrefixString() string
	LineSummaryString() string
	GetLineSummaryFunc() StringFunc
	// SetLineSummaryFunc(StringFunc)
	PrintTree(io.Writer) error
	// PACKAGE METHODS
	setLevel(int)
}
