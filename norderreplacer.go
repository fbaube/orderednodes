package orderednodes

type ReplacerFunc func(Norder) Norder

func ReplaceTree(oldRoot Norder, f ReplacerFunc) (newRoot Norder, err error) { // (*ONoder, error) {
	if (oldRoot.Parent() != nil) || !oldRoot.IsRoot() {
		println("ReplaceTree: did not get a parentless root")
	}
	// HANDLE ROOT
	newRoot = f(oldRoot)
	oldRoot.ReplaceWith(newRoot)

	var old, new Norder
	old = newRoot.FirstKid()
	for old != nil {
		new = f(old)
		old.ReplaceWith(new)
		old = old.NextKid()
	}
	return nil, nil
}

/*
func UpgradeNordToFilePropsNord(inNord Norder) Norder {
	return nil
}
*/
