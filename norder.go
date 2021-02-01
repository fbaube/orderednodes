package orderednodes

import "io"

// Norder is satisfied by *Nord NOT by Nord
type Norder interface {
	SeqId() int
	Level() int
	RelFP() string
	AbsFP() string
	IsRoot() bool
	GetRoot() Norder
	SetIsRoot(bool)
	Parent() Norder
	HasKids() bool
	FirstKid() Norder
	PrevKid() Norder
	NextKid() Norder
	KidsAsSlice() []Norder
	AddKid(Norder) Norder
	ReplaceWith(Norder) Norder
	SetLevel(int)
	SetParent(Norder)
	SetPrevKid(Norder)
	SetNextKid(Norder)
	LinePrefixString() string
	LineSummaryString() string
	PrintAll(io.Writer) error
}
