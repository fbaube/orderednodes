package orderednodes

// NordEngine tracks the state of a Nord tree being assembled,
// for example when a directory is specified for recursive analysis.
type NordEngine struct {
	// nexSeqID should be reset to 0 when starting another tree ?
	// No, because every single entity (dir/file) gets one,
	// even if it is listed on the CLI as an individual file.
	// nexSeqID      int
	rootPath      string
	summaryString StringFunc
}

// NordEng is a package global, which is dodgy and not re-entrant.
// The solution probably involves currying. 
var NordEng *NordEngine = new(NordEngine)

