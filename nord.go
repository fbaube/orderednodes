package orderednodes

import (
	"fmt"
	L "github.com/fbaube/mlog"
	"os"
	FP "path/filepath"

	FU "github.com/fbaube/fileutils"
)

// StringFunc is actually: func (*Norder) FuncName() string
type StringFunc func(Norder) string

// NOTE: Defining NewNord(Path) and NewRootNord() could remove the need
// for several of the setters defined below.

// NOTE: Ignore https://godoc.org/golang.org/x/net/html#Node
// (and many other available implementations of "Node" data structure).

// NOTE: Compared to usage for XML, usage for files & dirs is much more
// strictly typed. Dirs are dirs and files are files and never the twain
// shalll meet. This means that dirs can never contain content and that
// files can never be non-leaf nodes. However this kind of typing is too
// complex to handle here in an OrderedNode, so it is handled instead by
// (for example) an instance of [fileutils.FSItem]. 
//
// Nord is a Node but with ordered children nodes: the child nodes have 
// a specific specified order. It lets us define such funcs as FirstKid(),
// NextKid(), PrevKid(), LastKid(). They are defined in interface Norder.
// A Nord also contains its own relative (to its inbatch) and absolute paths.
//
// There are three use cases identified for Nords: files & directories, DOM
// markup, and [Lw]DITA map files. Note that we never have two identically
// named files in the same directory, but that we might (for example) have
// multiple sibling <p> tags. So when representing markup, a map from paths
// to Nords will fail unless the tags are made unique with subscript indices,
// such as "[1]", "[2]", like those used in (for example) JQuery.
//
// If we build up a tree of Nords when processing an [os.DirFS], the strict
// ordering provided by DirFS is not strictly needed, BUT it can anyways be
// used (and relied upon) because [io.fs.WalkDir] is deterministic (using
// lexical order). It means that a given Nord will always appear AFTER the
// Nord for its directory has appeared, which makes it much easier to build
// a tree.
//
// Fields are lower-cased so that other packages cannot disrupt links. 
//
// *Implementation note:* We use a doubly-linked list, not a slice.
// Since a Nord does not store a complete set of pointers to all of its kids,
// for example in a slice, it is not feasible to define a simpler variant of
// Node (with unordered kids) that could then be embedded in Nord. Nor is it
// simple to get a kid count.
// .
type Nord struct {
	// Path is the relative path of this Nord, relative to the root
	// Nord, which is the "local root" shared with its peers. (For
	// example: the root of a directory tree imported in a single
	// batch.) The last element of the Path is this Nord's own
	// label, i.e. FP.Base(Path).
	path string // relFP
	// absPath is the same as path, but rooted in - and in a 
	// dir tree, incorporating - the root node's absolute path.
	// For a file, it is rooted at the filesystem root.
	// For markup, it is rooted at the beginning of the document.
	absPath FU.AbsFilePath

	parent            Norder // level up
	firstKid, lastKid Norder // level down
	prevKid, nextKid  Norder // level same (rename "Kid" => "Peer" ?) 
	isRoot            bool   // level topmost
	// level is equal to the number of "/" filepath separators
	// separating path elements (i.e. not including any trailing
	// separator). Therefore it is 0 for XML root node (where
	// IsRoot() is true and Parent() is nil), 0 for top-level
	// files & directories, and >0 for others. Reserve negative
	// numbers for future (ab)use.
	level int
	// seqID is a unique ID under this node's tree's root. It does not 
	// need to be the same as (say) the index of this Nord in a slice 
	// of Nord's, but it probably is. Its use is optional, and also 
	// it can be used in other ways in structs that embed Nord.
	seqID int
	// parSeqID and kidSeqID's can add a layer of error checking and
	// simplified access. Their use is optional.
	// kidSeqIds when empty is ",", otherwise e.g. ",1,4,56,". The
	// seqIds should be in the same order as the Kid nodes themselves.
	// The bracketing by commas makes searching simpler.
	parSeqID, kidSeqID string
	lineSummaryFunc    StringFunc
	isDir              bool
}

// RootNord is defined so as to require 
// explicit assignments to/from a root node.
type RootNord Nord

// IsDir does NOT work, because we are not setting bool isDir yet.
// It is set (or not set) in embedding structs.
func (p *Nord) IsDir() bool {
	return p.isDir
}

// NewRootNord verified directoryness and
// then sets the bools [isRoot] and [isDir].
func NewRootNord(rootPath string, smryFunc StringFunc) *Nord {
	p := new(Nord)
	// p.lineSummaryFunc = NordSummaryString // func
	L.L.Info("NewRootNord: starting seqID: %d", NordEng.nexSeqID)
	if rootPath == "" {
		L.L.Error("NewRootNord: missing root path")
		return nil 
	}
	p.seqID = NordEng.nexSeqID
	NordEng.nexSeqID += 1
	p.path = "."
	// NOTE thenext stmt assumes *files* not DOM
	p.absPath = FU.AbsFP(FP.Clean(rootPath))
	// Verify that it is in fact a directory
	fm := FU.NewFileMeta(p.absPath.S())
	if !fm.IsDir() {
		L.L.Error("NewRootNord: path is not a dir: " + p.absPath.S())
		return nil
	}
	p.isRoot = true
	p.isDir = true 
	NordEng.summaryString = smryFunc
	return p
}

// NewNord does not set (or unset) the bool [isDir] because checking 
// it is an expensive operation that can and should be done elsewhere. 
func NewNord(aPath string) *Nord {
	if aPath == "" {
		println("newNord: missing path")
	}
	p := new(Nord)
	// p.lineSummaryFunc = NordSummaryString // func
	p.seqID = NordEng.nexSeqID
	NordEng.nexSeqID += 1
	// L.L.Dbg("NewNord: seqID is now %d", NordEng.nexSeqID)
	p.path = aPath
	p.absPath = FU.AbsFP(FP.Join(NordEng.rootPath, aPath))
	p.lineSummaryFunc = NordEng.summaryString
	// p.isDir =
	return p
}

func (p *Nord) GetLineSummaryFunc() StringFunc {
	return p.lineSummaryFunc
}

func (p *Nord) SetLineSummaryFunc(sf StringFunc) {
	p.lineSummaryFunc = sf
}

// IsRoot is duh.
func (p *Nord) IsRoot() bool {
	return p.isRoot
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

// setlevel is duh.
func (p *Nord) setLevel(i int) {
	p.level = i
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
