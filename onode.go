package orderednodes

import (
	"fmt"
	"os"
	S "strings"
)

// Ignore https://godoc.org/golang.org/x/net/html#Node

// ONode is an Ordered node: the child nodes have a specific specified order.
// It lets us define such funcs as FirstKid(), NextKid(), PrevKid(), LastKid().
// They are defined in interface ONoder.
//
// If we build up a tree of ONode's when processing an os.DirFS, the strict
// ordering is not strictly needed, but it can be used (and relied upon)
// because io.fs.WalkDir is deterministic (using lexical order).
//
// *Implementation note:* We use a doubly-linked list, not a slice.
//
type ONode struct {
	parent            *ONode
	firstKid, lastKid *ONode
	prevKid, nextKid  *ONode
	isRoot bool
	// level is 0 for root node, and >0 for others.
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

// Available to ensure that assignments to/from root node are explicit.
type RootONode ONode

// ONoder is satisfied by *ONode NOT ONode
type ONoder interface {
	IsRoot() bool
	SeqId()  int
	Parent()   *ONode
	FirstKid() *ONode
	NextKid()  *ONode
	KidsAsSlice() []*ONode
	AddKid(*ONode)  *ONode // Noder; PASS BY REFERENCE
}

// Parent returns the parent, duh.
func (p *ONode) IsRoot() bool {
	return p.isRoot
}

// SeqId is duh.
func (p *ONode) SeqId() int {
	return p.seqId
}

// Parent returns the parent, duh.
func (p *ONode) Parent() *ONode {
	return p.parent
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *ONode) AddKid(aKid *ONode) *ONode {
	if aKid.prevKid != nil || aKid.nextKid != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if aKid.parent != nil && aKid.parent != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = p.firstKid
	var LK = p.lastKid
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		p.firstKid, p.lastKid = aKid, aKid
		aKid.parent = p
		aKid.prevKid, aKid.nextKid = nil, nil
		return aKid
	}
	// So, replace the last kid
	if LK != nil {
		if LK.parent != p {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.nextKid != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.nextKid = aKid
		aKid.prevKid = LK
		p.lastKid = aKid
		aKid.parent = p
		return aKid
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
	panic("AddKid: Chaos!")
}

// FirstKid provides read-only access for other packages. Can return nil.
func (p *ONode) FirstKid() *ONode {
	return p.firstKid
}

// NextKid provides read-only access for other packages. Can return nil.
func (p *ONode) NextKid() *ONode {
	return p.nextKid
}

// Echo implements Markupper.
func (p *ONode) Echo() string {
	panic("recursion") // return p.Echo()
}

// AsLinePrefix provides indentation and should start a line of display/debug.
// It does not end the string with (white)space.
func (p ONode) AsLinePrefix() string {
	if p.isRoot { // && p.Parent == nil
		return "[R]"
	} else if (p.level == 0 && p.Parent != nil) {
		return fmt.Sprintf("[%d]", p.seqId)
	} else {
		// (spaces)[lvl:seq]"
		// func S.Repeat(s string, count int) string
		return fmt.Sprintf("%s[%d:%02s]",
			S.Repeat("  ", p.level-1), p.level, p.seqId)
 	}
}

/* String implements Markupper.
func (p ONode) String() string {
	var s = p.Echo() +
		// fmt.Sprintf(" [d%d:%d] ", p.Depth, p.MatchingTagsIndex) +
		fmt.Sprintf(" [%d] ", p.level) + p.TagSummary.String()
	return s
}

// StringRecursively is fab.
func (p ONode) StringRecursively(s string, iLvl int) string {

	s += SU.GetIndent(iLvl) + p.String() + "\n" // p.GToken.Echo() +
	// fmt.Sprintf(" d%d::[%d] ", p.Depth, p.MatchingTagsIndex) +
	// p.TagSummary.String() + "\n"
	var kids []*ONode
	kids = p.KidsAsSlice()
	for _, k := range kids {
		s += k.StringRecursively(s, iLvl+1)
	}
	return s
}
*/

func (p *ONode) KidsAsSlice() []*ONode {
	var pp []*ONode
	c := p.firstKid
	for c != nil {
		pp = append(pp, c)
		c = c.nextKid
	}
	return pp
}
