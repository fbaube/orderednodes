package orderednodes

import (
	"fmt"
	"io/fs"
	"os"
)

type PropPathFS struct {
	InputFS fs.FS
	root    *OPPNode
}

func NewPropPathFS(path string) *PropPathFS {
	// var e error
	var ppfs *PropPathFS
	fmt.Println("on.newppfs:", path)
	ppfs = new(PropPathFS)
	ppfs.InputFS = os.DirFS(path)
	// func WalkDir(fsys FS, root string, fn WalkDirFunc) error
	fs.WalkDir(ppfs.InputFS, ".", myWalkFn)
	return ppfs
}

// type WalkDirFunc func(path string, d DirEntry, err error) error
func myWalkFn(path string, d fs.DirEntry, err error) error {
	fmt.Println("Walking:", path)
	return nil
}
