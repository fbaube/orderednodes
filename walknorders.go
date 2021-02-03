package orderednodes

type NorderVisiterFunc func(pNode Norder) error

func WalkNorders(p Norder, nvf NorderVisiterFunc) error {
	var e error
	if e = nvf(p); e != nil {
		return e
	}
	pKid := p.FirstKid()
	for pKid != nil {
		if e = WalkNorders(pKid, nvf); e != nil {
			return e
		}
		pKid = pKid.NextKid()
	}
	return nil
}
