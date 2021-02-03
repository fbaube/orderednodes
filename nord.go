package orderednodes

import (
	"fmt"
	"io"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
)

// StringFunc is actually: func (*Norder) FuncName() string
type StringFunc func(Norder) string

// NOTE: Defining NewNord(Path) and NewRootNord() could remove the need
// for several of the setters defined below.

// NOTE: Ignore https://godoc.org/golang.org/x/net/html#Node
// (and many other available implementations of "Node" data structure).

// Nord is a Node but with ordered children nodes: the child nodes have a
// specific specified order. It lets us define such funcs as FirstKid(),
// NextKid(), PrevKid(), LastKid(). They are defined in interface Norder.
// A Nord also contains its own relative and absolute paths.
//
// There are three use cases identified for Nords: files & directories, DOM
// markup, and [Lw]DITA map files. Note that we never have two identically
// named files in the same directory, but that we might (for example) have
// multiple sibling <p> tags. So when representing markup, a map from paths
// to Nords will fail unless the tags are made unique with subscript indices,
// such as "[1]", "[2]", like those used in (for example) JQuery.
//
// If we build up a tree of Nords when processing an os.DirFS, the strict
// ordering provided by DirFS is not strictly needed, BUT it can anyways
// be used (and relied upon) because io.fs.WalkDir is deterministic (using
// lexical order). It means that a given Nord will always appear AFTER the
// Nord for its directory has appeared, which makes it much easier to build
// a tree.
//
// *Implementation note:* We use a doubly-linked list, not a slice.
// Since a Nord does not store a complete set of pointers to all of its kids,
// for example in a slice, it is not feasible to define a simpler variant of
// Node (with unordered kids) that could then be embedded in Nord. Nor is it
// simple to get a kid count.
//
type Nord struct {
	// Path is the relative path of this Nord. The last element
	// of the Path is this Nord's own label, i.e. FP.Base(Path)
	path string // relFP
	// absPath is the same as path, but rooted in the root node's path.
	// For a file, it is rooted at the filesystem root.
	// For markup, it is rooted at the beginning of the document.
	absPath FU.AbsFilePath

	parent            Norder // level up
	firstKid, lastKid Norder // level down
	prevKid, nextKid  Norder // level same
	isRoot            bool   // level topmost
	// level is equivalent to the number of "/" filepath seperators.
	// Therefore it is 0 for root node (where IsRoot() is true and Parent()
	// is nil), 0 for top-level files & directories, and >0 for others.
	// Reserve negative numbers for future (ab)use.
	level int
	// seqID is a unique ID under this node's tree's root. It does not need
	// to be the same as (say) the index of this Nord in a slice of Nord's,
	// but it probably is. Its use is optional, and also it can be used in
	// other ways in structs that embed Nord.
	seqID int
	// parSeqID and kidSeqID's can add a layer of error checking and
	// simplified access. Their use is optional.
	// kidSeqIds when empty is ",", otherwise e.g. ",1,4,56,". The
	// seqIds should be in the same order as the Kid nodes themselves.
	// The bracketing by commas makes searching simpler.
	parSeqID, kidSeqID string
	lineSummaryFunc    StringFunc
}

type nordCreationState struct {
	nexSeqID      int // reset to 0 when doing another tree ?
	rootPath      string
	summaryString StringFunc
}

var pNCS *nordCreationState = new(nordCreationState)

func NewRootNord(rootPath string, smryFunc StringFunc) *Nord {
	p := new(Nord)
	// p.lineSummaryFunc = NordSummaryString // func
	if pNCS.nexSeqID != 0 {
		println("newRootNord: seqID is not zero")
	}
	if rootPath == "" {
		println("newRootNord: missing root path")
	}
	p.seqID = pNCS.nexSeqID
	pNCS.nexSeqID += 1
	p.path = "."
	// This next stmt assumes *files* not DOM
	p.absPath = FU.AbsFP(FP.Clean(rootPath))
	p.isRoot = true
	pNCS.summaryString = smryFunc
	println("newRootNord:", p.absPath)
	return p
}

func NewNord(aPath string) *Nord {
	if aPath == "" {
		println("newNord: missing path")
	}
	p := new(Nord)
	// p.lineSummaryFunc = NordSummaryString // func
	p.seqID = pNCS.nexSeqID
	pNCS.nexSeqID += 1
	p.path = aPath
	p.absPath = FU.AbsFP(FP.Join(pNCS.rootPath, aPath))
	p.lineSummaryFunc = pNCS.summaryString
	return p
}

func (p *Nord) GetLineSummaryFunc() StringFunc {
	return p.lineSummaryFunc
}

func (p *Nord) SetLineSummaryFunc(sf StringFunc) {
	p.lineSummaryFunc = sf
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
func (p *Nord) SeqID() int {
	return p.seqID
}

// Level is duh.
func (p *Nord) Level() int {
	return p.level
}

// Path is duh.
func (p *Nord) Path() string { return p.path }

// RelFP is duh.
func (p *Nord) RelFP() string { return p.path }

// AbsFP is not valid until set by the embedding struct.
func (p *Nord) AbsFP() string { return string(p.absPath) }

// Parent returns the parent, duh.
func (p *Nord) Parent() Norder {
	return p.parent
}

/*
// SetSeqId is duh.
func (p *Nord) SetSeqId(i int) {
	p.seqId = i
}
*/

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

// SetFirstKid has no side effects.
func (p *Nord) SetFirstKid(p2 Norder) {
	p.firstKid = p2
}

// SetLastKid has no side effects.
func (p *Nord) SetLastKid(p2 Norder) {
	p.lastKid = p2
}

/*
// SetParent has no side effects.
func (p *ONode) SetParent(p2 ONoder) {
	p.parent = p2
}
*/

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *Nord) AddKid(aKid Norder) Norder { // returns aKid
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
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
		/*
			if aKid.Parent() != p {
				panic("BAD PARENT 1")
			}
			println("OK PARENT 1")
		*/
		return aKid
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
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
		/*
			if aKid.Parent() != p {
				panic("BAD PARENT 2")
			}
			println("OK PARENT 2")
		*/
		return aKid
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
	panic("AddKid: Chaos!")
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func AddKid2(par, kid Norder) { // Norder { // returns aKid
	// fmt.Printf("nord: ptrs? par<%T> kid<%T> \n", par, kid)
	if kid.PrevKid() != nil || kid.NextKid() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", par, kid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if kid.Parent() != nil && kid.Parent() != par {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", par, kid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = par.FirstKid()
	var LK = par.LastKid()
	// Set the level now
	kid.SetLevel(par.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		par.SetFirstKid(kid)
		par.SetLastKid(kid)
		kid.SetParent(par)
		kid.SetPrevKid(nil)
		kid.SetNextKid(nil)
		/*
			if kid.Parent() != par {
				panic("BAD PARENT 1")
			}
			println("OK PARENT 1")
		*/
		return
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != par {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.NextKid() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNextKid(kid) // LK.nextKid = aKid
		kid.SetPrevKid(LK) // aKid.prevKid = LK
		par.SetLastKid(kid)
		kid.SetParent(par)
		/*
			if kid.Parent() != par {
				panic("BAD PARENT 2")
			}
			println("OK PARENT 2")
		*/
		return
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
	panic("AddKid: Chaos!")
}

// ===

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func AddKid3(par, kid *Nord) { // Norder { // returns aKid
	// fmt.Printf("nord: ptrs? par<%T> kid<%T> \n", par, kid)
	if kid.PrevKid() != nil || kid.NextKid() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", par, kid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if kid.Parent() != nil && kid.Parent() != par {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", par, kid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = par.firstKid
	var LK = par.lastKid
	// Set the level now
	kid.SetLevel(par.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		par.firstKid = kid
		par.lastKid = kid
		kid.parent = par
		kid.prevKid = nil
		kid.nextKid = nil
		/*
			if kid.Parent() != par {
				panic("BAD PARENT 1")
			}
			println("OK PARENT 1")
		*/
		return
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != par {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.NextKid() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNextKid(kid) // LK.nextKid = aKid
		kid.prevKid = LK   // aKid.prevKid = LK
		par.lastKid = kid
		kid.parent = par
		/*
			if kid.Parent() != par {
				panic("BAD PARENT 2")
			}
			println("OK PARENT 2")
		*/
		return
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", par, kid)
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

// LastKid provides read-only access for other packages. Can return nil.
func (p *Nord) LastKid() Norder {
	return p.lastKid
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
	} else if p.Level() == 0 && p.Parent() != nil {
		return fmt.Sprintf("[%d]", p.seqID)
	} else {
		// (spaces)[lvl:seq]"
		// func S.Repeat(s string, count int) string
		return fmt.Sprintf("%s[%02d:%02d]",
			S.Repeat("  ", p.level-1), p.level, p.seqID)
	}
}

func yn(b bool) string {
	if b {
		return "Y"
	} else {
		return "n"
	}
}

// func (p *Nord) NordSummaryString() string {
func (p *Nord) LineSummaryString() string {
	var sb S.Builder
	if p.IsRoot() {
		sb.WriteString("ROOT ")
	}
	/*
		if p.PrevKid() != nil {
			sb.WriteString("P ")
		}
		if p.Parent() == nil {
			sb.WriteString("NOPARENT ")
		}
		if p.NextKid() != nil {
			sb.WriteString("N ")
		}
		if p.HasKids() {
			sb.WriteString("kid(s) ")
		}
	*/
	if p.path == "" {
		sb.WriteString("NOPATH")
	} else {
		sb.WriteString(p.path)
	}
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
	// var F StringFunc
	// F = p.GetLineSummaryFunc()
	// fmt.Fprintf(printAllTo, "%s %s (%T) \n", p.LinePrefixString(), F(p), p)
	// fmt.Fprintf(printAllTo, "%s %s (%T) \n", p.LinePrefixString(), "?", p)
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
