package orderednodes

import (
	"fmt"
	"os"
	FP "path/filepath"
	L "github.com/fbaube/mlog"
	FU "github.com/fbaube/fileutils"
)

// StringFunc is used by interface Norder, so its
// method signature actually (MAYBE!) looks like:
// func (*Nord) FuncName() string
type StringFunc func(Norder) string

// NOTE that in this file, we endeavor to ignore the many, many
// other implementatoins "out there" of the Node data structure
// type, including https://godoc.org/golang.org/x/net/html#Node

// Nord is shorthand for "ordered node" (or "ordinal node") - a Node with
// ordered children nodes: the child nodes have a specific specified order.
// Ordering is essential for representation of content.
//
// (It can help in other contexts too, such as filesystem operations,
// but it turns out that the Go stdlib generally returns directory
// items in lexical order and walks directories in lexical order.)
// 
// The ordering lets us define funcs like FirstKid, NextKid, PrevKid,
// LastKid. They are defined in interface [Norder]. A Nord for file/dir
// also contains its own relative path (relative to its inbatch) and
// absolute path. 
//
// Alternatives Everywhere
// 
// There are three use cases identified for Nords:
//  - files & directories (here ordering is less important) 
//  - XML/HTML/DOM markup (this creates problems handling same-named siblings)
//  - [Lw]DITA map files and other ToC's (these should be an ideal use case) 
//
// Also there are two distinct memory management nodes for allocating and
// linking nodes:
//  - The "traditional method" of allocating nodes individually, and linking
//    them using pointers. Using this method, both deletions and insertions
//    are relatively simple.
//     - Such nodes can also be loaded into a map, for random access to nodes
//       based on path. 
//  - The "new-fangled way" called an "arena", where we put all our nodes in
//    a big slice, and link them using indices. This method is much kinder
//    on memory management, but might becomes clumsy when we need dynamic
//    node management.
//     - Deletions are easy if we just zero out the slice entry; we cannot
//       then do compaction because it would require updating all indices
//       past the first point of compaction.
//     - Insertions are costly. However note that in this implementation,
//       for a node's kids, we use a linked list rather than a slice, so
//       this makes it easier to append a new node at the end of the
//       arena-slice and then update indices, wherever they may be
//       elsewhere in the slice.
//     - In any case, if an arena-slice has to grow (because of a call 
//       to append), it might be moved elsewhere in memory, which would
//       invalidate all ptrs to other Nords! If this happens, we should
//       trust only the "traditional method".
//
// Also there are multiple ways to represent node trees in our SQLite DBMS,
// and multiple ways to walk a node tree, so there is a unavoidable complexity
// wherever we look. 
// 
// NOTE: DOM markup exhibits name duplication: we never have two same-named
// files in the same directory, but we might (e.g.) have multiple sibling
// <p> tags. So when representing markup, a map from paths to Nords would 
// fail unless the tags are made unique with subscript indices, such as
// "[1]", "[2]", like those used in (e.g.) JQuery.
//
// NOTE: Using Nords for files & dirs exhibits strong typing. Dirs are
// dirs and files are files and never the twain shall meet. This means
// that dirs cannot contain own-content and that files can never be
// non-leaf nodes. (Note tho that symlinks have aspects of both.) 
// However this dir/file/etc typing is too complex to handle here in 
// a Nord, so it is handled instead by a struct type that embeds Nord, 
// such as [fileutils.FSItem]. 
//
// If we build up a tree of Nords when processing an [os.DirFS], the 
// strict ordering provided by DirFS is not strictly needed, BUT it 
// can anyways be used (and relied upon) because the three flavors 
// of WalkDir are deterministic (using lexical order). It means that 
// a given Nord will always appear AFTER the Nord for its directory
// has appeared, which makes it much easier to build a tree.
//
// Link fields are lower-cased so that other packages cannot damage links. 
//
// NOTE: This implementation stores pointers to child nodes in a doubly
// linked list, not a slice, and a Nord does not have a complete set of
// pointers to all of its kids. Therefore (a) it is not simple to get a
// kid count, because it requires a list traversal, and (b) it is not
// feasible to use this same code to define a simpler, more efficient 
// variant of Nord that has unordered kids. 
// .
type Nord struct {

     	// ----------------------------------
     	//  Substructure: Materialized Paths
     	// ----------------------------------	
	// relPath is the relative path of this Nord, relative to its 
	// tree's root Nord, which is the "local root" shared w other 
	// Nords in the same interconnected tree. (That is to say, a 
	// local root is the highest/topmost node of a directory tree 
	// imported in a single batch.) The last element of the relPath 
	// is this Nord's own name/label, analagous to FP.Base(Path).
	relPath string 
	// absPath is the same as path, except that it is rooted in 
	// - i.e. it is traced back to the root of - a local file 
	// system (or documwnt). For a file or dir in a filesystem, 
	// it is rooted at the filesystem root. For a markup node 
	// or a map/ToC file, it is rooted at the document start.
	absPath FU.AbsFilePath	
	// isRoot true has a relPath of "." and an absPath that is the 
	// rooted absolute path of this root node w.r.t. the external
	// environment (for a file or dir, the file system root; for
	// a markup node, the absolute path of the containing file. 
	isRoot bool
	// isDir is obvious for files & dirs BUT not (yet) for symlinks. 
	isDir  bool
	// level is equal to the number of "/" filepath separators
	// *separating* path elements (i.e. not including any leading 
	// or trailing separators). Therefore it is 0 for an XML docu-
	// ment root node or the local root of a file & dir tree (where 
	// in both cases, isRoot() is true and parent() is nil)), and it 
	// is >0 for others. Reserve negative numbers for future (ab)use.
	level int
     	// ----------------------------------
	//  Temporarily unused: three fields	
     	// ----------------------------------
	// seqID is a unique ID under this node's tree's root. It does not 
	// need to be the same as (say) the index of this Nord in a slice 
	// of Nord's, but it probably is. Its use is optional, and also 
	// it can be used in other ways in structs that embed Nord.
	// >> seqID int
	// parSeqID and kidSeqID's can add a layer of error checking 
	// and simplified access. Their use is optional.
	// kidSeqIds when empty is ",", otherwise e.g. ",1,4,56,". 
	// the seqIds should (well, can) be in the same order as 
	// the Kid nodes themselves. The bracketing by commas makes 
	// searching simpler (",%d,").
	// >> parSeqID, kidSeqID string
     	// ---------------------------------
	//  Substructure for Adjacency List 
	//   based on Go ptrs (not indices)
     	// ---------------------------------	
	parent            Norder // level up
	firstKid, lastKid Norder // level down
	prevKid, nextKid  Norder // level same (rename "Kid" => "Peer" ?)
     	// ---------------------------------
	//  Substructure for Adjacency List
	//  based on indices into "arena"
	//  slice (not using ptrs) 
     	// ---------------------------------
	// Alternative 1 (used): use the same idea 
	// as when using ptrs, i.e. point to parent 
	// and to peers in own list and to first and
	// last kids in the linked list of all kids. 
	iParent             int // level up
	iFirstKid, iLastKid int // level down
	iPrevKid, iNextKid  int // level same (rename "Kid" => "Peer" ?)
	// Alternative 2 (unused): a single string 
	// field holds all applicable indices.
	// kidIdxs when empty is ",", else e.g. ",1,4,56,". 
	// The kidIdxs should be in the same order as the 
	// Kid nodes themselves. The bracketing by commas 
	// makes searching simpler (",%d,").
	// >> parIdx, kidIdxs string

	// This is a handy utility function, but 
	// maybe just implementing+using interface
	// Stringser would be a better idea. 
	lineSummaryFunc StringFunc
}

// RootNord is defined, so that assignments
// to/from a root node have to be explicit.
type RootNord Nord

// IsDir does NOT work, because we are not setting bool isDir yet.
// It is set (or not set) in embedding structs.
func (p *Nord) IsDir() bool {
	return p.isDir
}

// NewRootNord verifies it got a directory, and then sets the bools
// [isRoot] and [isDir]. Note that the passed-in field [rootPath] is
// set elsewhere, and must be set in the global [NordEng] before any
// child Nord is created using [NewNord].
func NewRootNord(rootPath string, smryFunc StringFunc) *Nord {
	// L.L.Debug("NewRootNord: starting seqID: %d", NordEng.nexSeqID)
	if rootPath == "" {
		L.L.Error("NewRootNord: missing root path")
		return nil 
	}
	// NOTE the next stmts assume *filesystem* not XML DOM 
	asAbsPath := FU.EnsureTrailingPathSep(FP.Clean(rootPath))
	// Verify that it is in fact a directory
	if !FU.IsDirAndExists(asAbsPath) {
		L.L.Error("NewRootNord: path is not a dir: " + asAbsPath)
		return nil
	}
	p := NewNord(asAbsPath)
	if p == nil { return nil }

	// CHECK THE PATHS
	L.L.Debug("RootNode's abs: " + p.absPath.S())
	L.L.Debug("RootNode's rel: " + p.relPath)

	p.absPath = FU.AbsFP(asAbsPath)
	// For the relative path, try to trim the entire
	// RootNode RootPath off of this absolute path.
	// func CutPrefix(s, prefix string) (after string, found bool):
	// It returns s without the provided leading prefix 
	// string and reports whether it found the prefix.
	// If s dusn't start with prefix, CutPrefix returns (s, false).
	// If prefix is the empty string, CutPrefix returns (s, true).
	
	p.relPath = p.absPath.S()
	p.isRoot = true
	p.isDir = true 
	return p
}

// NewNord expects a relative path (!!), and does not either
// (a) set/unset the bool [isDir] or (b) load file content,
// because these are expensive operations that can and should
// be done elsewhere, and also (c) they do not apply if this
// is being used for XML DOM. 
func NewNord(aRelPath string) *Nord {
	if aRelPath == "" {
		L.L.Error("NewNord: missing path")
		return nil 
	}
	p := new(Nord)
	// p.lineSummaryFunc = NordSummaryString // func
	// p.seqID = NordEng.nexSeqID
	// NordEng.nexSeqID += 1
	// L.L.Debug("NewNord: seqID is now %d", NordEng.nexSeqID)
	p.relPath = aRelPath
	asAbsPath := FP.Join(NordEng.rootPath, aRelPath)
	if FU.IsDirAndExists(asAbsPath) {
	   asAbsPath = FU.EnsureTrailingPathSep(asAbsPath)
	   }
	p.absPath = FU.AbsFilePath(asAbsPath) 
	p.lineSummaryFunc = NordEng.summaryString
	// p.isDir =... sorry, not done here 
	return p
}

func (p *Nord) LineSummaryFunc() StringFunc {
	return p.lineSummaryFunc
}

/*
func (p *Nord) SetLineSummaryFunc(sf StringFunc) {
	p.lineSummaryFunc = sf
}
*/

// IsRoot is duh.
func (p *Nord) IsRoot() bool {
	return p.isRoot
}

// Root walks the tree upward until [IsRoot] is true,
// so it does not use any global variables.
func (p *Nord) Root() RootNorder {
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

/*
// SeqId is duh.
func (p *Nord) SeqID() int {
	return p.seqID
} */

// Level is duh.
func (p *Nord) Level() int {
	return p.level
}

// AbsFP is duh.
func (p *Nord) AbsFP() string { return p.absPath.S() }

// RelFP is duh.
func (p *Nord) RelFP() string { return p.relPath }

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
