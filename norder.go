package orderednodes

// Norder is satisfied by *Nord NOT by Nord
type Norder interface {
	SeqId() int
	RelFP() string
	AbsFP() string
	IsRoot() bool
	GetRoot() Norder
	Parent() Norder
	HasKids() bool
	FirstKid() Norder
	PrevKid() Norder
	NextKid() Norder
	KidsAsSlice() []Norder
	AddKid(Norder) Norder      // *ONode // Noder; PASS BY REFERENCE
	ReplaceWith(Norder) Norder // *ONode // Noder; PASS BY REFERENCE
	SetParent(Norder)
	SetPrevKid(Norder)
	SetNextKid(Norder)
	LinePrefix() string
	LineSummary() string
}
