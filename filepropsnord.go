package orderednodes

import (
	FU "github.com/fbaube/fileutils"
)

// Ignore https://godoc.org/golang.org/x/net/html#Node

// FilePropsNord is an Ordered Propertied Path node:
// NOT ONLY the child nodes have a specific specified order
// BUT ALSO each node has a filepath plus the file properties.
// This means Pthat every Parent node is a directory.
//
// It also means we can use the redundancy to do a lot of error checking.
// Also we can use fields of seqId's to store parent and kid seqId's,
// adding yet another layer of error checking and simplified access.
//
type FilePropsNord struct {
	Nord
	FU.PathProps
}

// Available to ensure that assignments to/from root node are explicit.
type RootFilePropsNord FilePropsNord

// type OPPNoder interface { }

// AbsFP is duh.
func (p *FilePropsNord) AbsFP() string {
	return p.PathProps.AbsFP()
}

// RelFP is duh.
func (p *FilePropsNord) RelFP() string {
	return p.PathProps.RelFP()
}

// type WalkONoderFunc func(pNode ONoder) error
// func WalkONoders(p ONoder, wfn WalkONoderFunc) error {

var pathToBeFound string
var pathIsFound Norder

func (p *FilePropsNord) FirstKid() Norder {
	var pp *Nord
	// var ok bool
	pp = &(p.Nord)
	pp, _ = pp.FirstKid().(*Nord)
	var ppp Norder
	var pppp *FilePropsNord
	ppp = pp
	pppp, _ = ppp.(*FilePropsNord)
	return pppp
	// return p.FirstKid().(*pOPPNode)
}

func (pRoot *FilePropsNord) FindONoderByPath(path string) Norder {
	println("FindNorderByPath:", path)
	pathToBeFound = path
	pathIsFound = nil
	var onr Norder
	// var ok bool
	onr = pRoot
	e := WalkNorders(onr, nvfFindPath)
	if e != nil {
		println("wfnFindPath ERR:", e.Error())
		return nil
	}
	if pathIsFound == nil {
		println("wfnFindPath: No luck")
		return nil
	}
	return pathIsFound
}

func nvfFindPath(p Norder) error {
	if pathToBeFound == p.RelFP() {
		pathIsFound = p
	}
	return nil
}
