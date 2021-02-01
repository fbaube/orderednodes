package orderednodes

import (
	FU "github.com/fbaube/fileutils"
)

type FileNord struct {
	Nord
	argPath string // relFP
	absFP   FU.AbsFilePath
}
