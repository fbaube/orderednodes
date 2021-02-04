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
