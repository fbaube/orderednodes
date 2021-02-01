package orderednodes

import (
	"fmt"
	"io"
	"os"
	S "strings"
)

// Ignore https://godoc.org/golang.org/x/net/html#Node

// Nord is a Node but with ordered children nodes: the child nodes have a
// specific specified order. It lets us define such funcs as FirstKid(),
// NextKid(), PrevKid(), LastKid(). They are defined in interface Norder.
//
// If we build up a tree of Nords when processing an os.DirFS, the strict
// ordering is not strictly needed, BUT it can anyways be used (and relied
// upon) because io.fs.WalkDir is deterministic (using lexical order).
//
// *Implementation note:* We use a doubly-linked list, not a slice.
// Since a Nord does not store a complete set of pointers to all of its
// kids, for example in a slice, it is not feasible to define a simpler
// Node (with unordered kids) that could then be embedded in Nord.
//
type Nord struct {
	parent            Norder
	firstKid, lastKid Norder
	prevKid, nextKid  Norder
	isRoot            bool
	// level is equivalent to the number of "/" filepath seperators.
	// Therefore it is 0 for root node (where IsRoot() is true),
	// 0 for top-level files & directories, and >0 for others.
	level int
	// seqId is a unique ID under this node's tree's root. It does not need
	// to be the same as (say) the index of this ONode in a slice of ONode's.
	// Its use is optional, and also it can be used in other ways in structs
	// that embed ONode.
	seqId int
	// parSeqId and kidSeqIds can add a layer of error checking and
	// simplified access. Their use is optional.
	// kidSeqIds when empty is ",", otherwise e.g. ",1,4,56,". The
	// seqIds should be in the same order as the Kid nodes themselves.
	// The bracketing by commas makes searching simpler.
	parSeqId, kidSeqIds string
}

// RootNord is available to make explicit assignments to/from root node.
type RootNord Nord

// IsRoot is duh.
func (p *Nord) IsRoot() bool {
	return p.isRoot
}

// SetIsRoot is duh.
func (p *Nord) SetIsRoot(b bool) {
	p.isRoot = b
}

// GetRoot is duh.
func (p *Nord) GetRoot() Norder {
	if p.IsRoot() {
		return p
	}
	var ondr Norder
	ondr = p
	for !ondr.IsRoot() {
		ondr = ondr.Parent()
	}
	return ondr
}

// HasKids is duh.
func (p *Nord) HasKids() bool {
	return p.firstKid != nil && p.lastKid != nil
}

// SeqId is duh.
func (p *Nord) SeqId() int {
	return p.seqId
}

// Level is duh.
func (p *Nord) Level() int {
	return p.level
}

// RelFP is a dummy.
func (p *Nord) RelFP() string { return "" }

// AbsFP is a dummy.
func (p *Nord) AbsFP() string { return "" }

// Parent returns the parent, duh.
func (p *Nord) Parent() Norder {
	return p.parent
}

// Setlevel is duh.
func (p *Nord) SetLevel(i int) {
	p.level = i
}

// SetParent has no side effects.
func (p *Nord) SetParent(p2 Norder) {
	p.parent = p2
}

// SetPrevKid has no side effects.
func (p *Nord) SetPrevKid(p2 Norder) {
	p.prevKid = p2
}

// SetNextKid has no side effects.
func (p *Nord) SetNextKid(p2 Norder) {
	p.nextKid = p2
}

/*
// SetParent has no side effects.
func (p *ONode) SetParent(p2 ONoder) {
	p.parent = p2
}
*/

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *Nord) AddKid(aKid Norder) Norder { // NOTE aKid was *Nord
	if aKid.PrevKid() != nil || aKid.NextKid() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = p.firstKid
	var LK = p.lastKid
	// Set the level now
	aKid.SetLevel(p.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		p.firstKid, p.lastKid = aKid, aKid
		aKid.SetParent(p)
		aKid.SetPrevKid(nil)
		aKid.SetNextKid(nil)
		return aKid
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != p {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.NextKid() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNextKid(aKid) // LK.nextKid = aKid
		aKid.SetPrevKid(LK) // aKid.prevKid = LK
		p.lastKid = aKid
		aKid.SetParent(p)
		return aKid
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
	panic("AddKid: Chaos!")
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (pOld *Nord) ReplaceWith(pNew Norder) Norder {
	// REPLACE SIBLINGS' SIBBLE-LINKS
	// REPLACE KIDS' PARENT-LINKS
	// REPLACE PARENT'S KID-LINK

	// We require that pNew has no existing links
	if pNew.PrevKid() != nil || pNew.NextKid() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in ReplaceWith: Tag<< %+v >> new<< %+v >>\n", pOld, pNew)
		panic("ReplaceWith(K) can't cos K has siblings")
	}
	if pNew.Parent() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in ReplaceWith: Tag<< %+v >> new<< %+v >>\n", pOld, pNew)
		panic("E.ReplaceWith(K) can't cos K has non-P parent")
	}
	// REPLACE SIBLINGS' SIBBLE-LINKS
	prv := pOld.PrevKid()
	if prv != nil {
		pNew.SetPrevKid(prv)
		prv.SetNextKid(pNew)
	}
	nxt := pOld.NextKid()
	if nxt != nil {
		pNew.SetNextKid(nxt)
		nxt.SetPrevKid(pNew)
	}
	// REPLACE KIDS' PARENT-LINKS
	if pOld.FirstKid() != nil {
		crntKid := pOld.firstKid // FirstKid()
		for crntKid != nil {
			crntKid.SetParent(pNew)
			pNew.AddKid(crntKid)
			crntKid = crntKid.NextKid()
		}
	}
	// REPLACE PARENT'S KID-LINK

	return pNew
}

// FirstKid provides read-only access for other packages. Can return nil.
func (p *Nord) FirstKid() Norder {
	return p.firstKid
}

// PrevKid provides read-only access for other packages. Can return nil.
func (p *Nord) PrevKid() Norder {
	return p.prevKid
}

// NextKid provides read-only access for other packages. Can return nil.
func (p *Nord) NextKid() Norder {
	return p.nextKid
}

// Echo implements Markupper.
func (p *Nord) Echo() string {
	panic("recursion") // return p.Echo()
}

// LinePrefixString provides indentation and should start a line of display/debug.
// It does not end the string with (white)space.
func (p Nord) LinePrefixString() string {
	if p.isRoot { // && p.Parent == nil
		return "[R]"
	} else if p.level == 0 && p.Parent != nil {
		return fmt.Sprintf("[%d]", p.seqId)
	} else {
		// (spaces)[lvl:seq]"
		// func S.Repeat(s string, count int) string
		return fmt.Sprintf("%s[%d:%02d]",
			S.Repeat("  ", p.level-1), p.level, p.seqId)
	}
}

func yn(b bool) string {
	if b {
		return "Y"
	} else {
		return "n"
	}
}
func (p Nord) LineSummaryString() string {
	var sb S.Builder
	if p.IsRoot() {
		sb.WriteString("ROOT ")
	}
	if p.PrevKid() != nil {
		sb.WriteString("prv ")
	}
	if p.NextKid() != nil {
		sb.WriteString("nxt ")
	}
	if p.HasKids() {
		sb.WriteString("kid(s) ")
	}
	sb.WriteString(fmt.Sprintf("id:%d L%d", p.seqId, p.level))
	return (sb.String())
}

var printAllTo io.Writer

func (p *Nord) PrintAll(w io.Writer) error {
	if w == nil {
		return nil
	}
	printAllTo = w
	e := WalkNorders(p, nvfPrintOneLiner)
	if e != nil {
		println("nvfPrintOneLiner ERR:", e.Error())
		return e
	}
	return nil
}

func nvfPrintOneLiner(p Norder) error {
	fmt.Fprintf(printAllTo, "%s %s \n",
		p.LinePrefixString(), p.LineSummaryString())
	return nil
}

/* String implements Markupper.
func (p ONode) String() string {
	var s = p.Echo() +
		// fmt.Sprintf(" [d%d:%d] ", p.Depth, p.MatchingTagsIndex) +
		fmt.Sprintf(" [%d] ", p.level) + p.TagSummary.String()
	return s
}
*/

func (p *Nord) KidsAsSlice() []Norder {
	var pp []Norder
	c := p.FirstKid() // p.firstKid
	for c != nil {
		pp = append(pp, c)
		c = c.NextKid() // c.nextKid
	}
	return pp
}
