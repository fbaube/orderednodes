package orderednodes

type InspectorFunc func(pNode Norder) error

// func InspectTree used to be func WalkNorders
// .
func InspectTree(p Norder, f InspectorFunc) error {
	var e error
	if e = f(p); e != nil {
		return e
	}
	pKid := p.FirstKid()
	for pKid != nil {
		if e = InspectTree(pKid, f); e != nil {
			return e
		}
		pKid = pKid.NextKid()
	}
	return nil
}
