package orderednodes

type NorderNewFunc func(Norder) Norder // (*ONode) *ONode

func ReplaceAllNorders(oldRoot Norder, nrf NorderNewFunc) (newRoot Norder, err error) { // (*ONoder, error) {
	if (oldRoot.Parent() != nil) || !oldRoot.IsRoot() {
		panic("ReplaceAllNorders: did not get a parentless root")
	}
	// HANDLE ROOT
	newRoot = nrf(oldRoot)
	oldRoot.ReplaceWith(newRoot)

	var old, new Norder
	old = newRoot.FirstKid()
	for old != nil {
		new = nrf(old)
		old.ReplaceWith(new)
		old = old.NextKid()
	}
	return nil, nil
}
