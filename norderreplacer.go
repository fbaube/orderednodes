package orderednodes

type NorderUpgradeFunc func(Norder) Norder

func ReplaceAllNorders(oldRoot Norder, nuf NorderUpgradeFunc) (newRoot Norder, err error) { // (*ONoder, error) {
	if (oldRoot.Parent() != nil) || !oldRoot.IsRoot() {
		println("ReplaceAllNorders: did not get a parentless root")
	}
	// HANDLE ROOT
	newRoot = nuf(oldRoot)
	oldRoot.ReplaceWith(newRoot)

	var old, new Norder
	old = newRoot.FirstKid()
	for old != nil {
		new = nuf(old)
		old.ReplaceWith(new)
		old = old.NextKid()
	}
	return nil, nil
}

func UpgradeNordToFilePropsNord(inNord Norder) Norder {
	return nil
}
