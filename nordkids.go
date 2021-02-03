package orderednodes

import (
	"fmt"
	"os"
)

// HasKids is duh.
func (p *Nord) HasKids() bool {
	return p.firstKid != nil && p.lastKid != nil
}

// Parent returns the parent, duh.
func (p *Nord) Parent() Norder {
	return p.parent
}

// SetParent has no side effects.
func (p *Nord) SetParent(p2 Norder) {
	p.parent = p2
}

// SetPrevKid has no side effects.
func (p *Nord) SetPrevKid(p2 Norder) {
	p.prevKid = p2
}

// SetNextKid has no side effects.
func (p *Nord) SetNextKid(p2 Norder) {
	p.nextKid = p2
}

// SetFirstKid has no side effects.
func (p *Nord) SetFirstKid(p2 Norder) {
	p.firstKid = p2
}

// SetLastKid has no side effects.
func (p *Nord) SetLastKid(p2 Norder) {
	p.lastKid = p2
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *Nord) AddKid(aKid Norder) Norder { // returns aKid
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	if aKid.PrevKid() != nil || aKid.NextKid() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = p.firstKid
	var LK = p.lastKid
	// Set the level now
	aKid.setLevel(p.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		p.firstKid, p.lastKid = aKid, aKid
		aKid.SetParent(p)
		aKid.SetPrevKid(nil)
		aKid.SetNextKid(nil)
		return aKid
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != p {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.NextKid() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNextKid(aKid) // LK.nextKid = aKid
		aKid.SetPrevKid(LK) // aKid.prevKid = LK
		p.lastKid = aKid
		aKid.SetParent(p)
		return aKid
	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
	panic("AddKid: Chaos!")
}

// FirstKid provides read-only access for other packages. Can return nil.
func (p *Nord) FirstKid() Norder {
	return p.firstKid
}

// LastKid provides read-only access for other packages. Can return nil.
func (p *Nord) LastKid() Norder {
	return p.lastKid
}

// PrevKid provides read-only access for other packages. Can return nil.
func (p *Nord) PrevKid() Norder {
	return p.prevKid
}

// NextKid provides read-only access for other packages. Can return nil.
func (p *Nord) NextKid() Norder {
	return p.nextKid
}

func (p *Nord) KidsAsSlice() []Norder {
	var pp []Norder
	c := p.FirstKid() // p.firstKid
	for c != nil {
		pp = append(pp, c)
		c = c.NextKid() // c.nextKid
	}
	return pp
}
