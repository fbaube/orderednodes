package orderednodes

// The goal is to use this code to write some new iterator -based code.
// It's complicated cos of the use of callbacks (or whatever you want to
// call them), but consider that the stdlib itself uses these kinds of
// callbacks in several places for walking hierarchies like filesystems.
//
// So let's document and analyse how the stdlib works - all three types
// of Walkers - and then use analogies to write an iterator-based version
// that is simple to comprehend.
//
// Anyways FP.WalkDir is the best. It doesn't try to hide symlinks away.
// And we are quite happy to use fileutils/FSItem 

import (
	"io/fs"
	"os"
	FP "path/filepath"
)

const (
	Sep     = os.PathSeparator
	ListSep = os.PathListSeparator
)

// SkipDir as a return value (from [NordProcFunc]) 
// says to skip the directory named in the call. 
// It is not returned as an error by any function.
var SkipDir error = fs.SkipDir

// SkipAll as a return value (from [NordProcFunc]) 
// says skip all remaining files and directories.
// It is not returned as an error by any function.
var SkipAll error = fs.SkipAll

// ================================================
// user API call is: 
// ================================================

// var lstat = os.Lstat // for testing

// WalkNorderTree walks the tree of [ON.Norder]s, calling `NorderProcFunc`
// for each Norder in the tree, including the input Norder. 
//
//?? All errors that arise visiting files and directories are filtered by fn:
//?? see the [fs.WalkDirFunc] documentation for details.
//
// The paths associated with the `Norder`s are, at each level downward,
// pretty much assumed to be in lexical order.
// .
func WalkNorderTree(root Norder, fn fs.WalkDirFunc) error {
	info, err := os.Lstat(root.AbsFP())
	if err != nil {
		err = fn(root.AbsFP(), nil, err)
	} else {
		err = walkDir(root.AbsFP(), fs.FileInfoToDirEntry(info), fn)
	}
	if err == SkipDir || err == SkipAll {
		return nil
	}
	return err
}

// walkDir recursively descends path, calling walkDirFn.
func walkDir(path string, d fs.DirEntry, walkDirFn fs.WalkDirFunc) error {
	if err := walkDirFn(path, d, nil); err != nil || !d.IsDir() {
		if err == SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(path, d, err)
		if err != nil {
			if err == SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}

	for _, d1 := range dirs {
		path1 := FP.Join(path, d1.Name())
		if err := walkDir(path1, d1, walkDirFn); err != nil {
			if err == SkipDir {
				break
			}
			return err
		}
	}
	return nil
}

// ========

/*

// Preorder returns an iterator over all 
// the nodes beneath (and including) the 
// specified root, in depth-first preorder.
//
// It stops iterating when geting the first
// error, so it avoids some design issues. 
//
// For greater control over the traversal
// of each subtree, use [Inspect].
// . 
func Preorder(root Norder) iter.Seq[Norder] {
	return func(yield func(Norder) bool) {
		ok := true
		Inspect(root, func(n Norder) bool {
			if n != nil {
				// yield must not be called once ok is false.
				ok = ok && yield(n)
			}
			return ok
		})
	}
}

*/

