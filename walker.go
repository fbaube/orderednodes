package orderednodes

type NorderVisiterFunc func(pNode Norder) error

func WalkNorderTree(p Norder, nvf NorderVisiterFunc) error {
	// println("on.walker: HANDLE SKIPDIR")
	if err := nvf(p); err != nil {
		return err
	}
	pKid := p.FirstKid()
	for pKid != nil {
		if err := nvf(pKid); err != nil {
			return err
		}
		pKid = pKid.NextKid()
	}
	return nil
}
