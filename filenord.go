package orderednodes

import (
	FU "github.com/fbaube/fileutils"
)

type FileNord struct {
	Nord
	Path string // relFP
	FU.AbsFilePath
}
