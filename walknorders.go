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

func InspectTreeWithPreAndPost(p Norder,
	f0 InspectorFunc, f1 InspectorFunc) error {

	var e error
	// PRE
	if e = f0(p); e != nil {
		return e
	}
	// KIDS
	pKid := p.FirstKid()
	for pKid != nil {
		if e = InspectTreeWithPreAndPost(pKid, f0, f1); e != nil {
			return e
		}
		pKid = pKid.NextKid()
	}
	// POST
	if e = f1(p); e != nil {
		return e
	}
	return nil
}
