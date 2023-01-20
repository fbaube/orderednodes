package orderednodes

import (
	"fmt"
	"io"
	"os"
	S "strings"
)

var printTreeTo, printCssTreeTo io.Writer

func init() {
	f, e := os.Create("./contentity-tree")
	if e == nil {
		printTreeTo = f
	} else {
		println("contentity-tree:", e.Error())
	}
	f, e = os.Create("./css-tree")
	if e == nil {
		printCssTreeTo = f
	} else {
		println("css-tree:", e.Error())
	}
}

// Echo implements Markupper.
func (p *Nord) Echo() string {
	panic("recursion") // return p.Echo()
}

// LinePrefixString provides indentation and
// should start a line of display/debug.
//
// It does not end the string with (white)space.
// .
func (p Nord) LinePrefixString() string {
	if p.isRoot { // && p.Parent == nil
		return "[R]"
	} else if p.Level() == 0 && p.Parent() != nil {
		return fmt.Sprintf("[%d]", p.seqID)
	} else {
		// (spaces)[lvl:seq]"
		// func S.Repeat(s string, count int) string
		return fmt.Sprintf("%s[%02d:%02d]",
			S.Repeat("  ", p.level-1), p.level, p.seqID)
	}
}

func yn(b bool) string {
	if b {
		return "Y"
	} else {
		return "n"
	}
}

func (p *Nord) LineSummaryString() string {
	var sb S.Builder
	if p.IsRoot() {
		sb.WriteString("ROOT ")
	}
	/*
		if p.PrevKid() != nil {
			sb.WriteString("P ")
		}
		if p.Parent() == nil {
			sb.WriteString("NOPARENT ")
		}
		if p.NextKid() != nil {
			sb.WriteString("N ")
		}
		if p.HasKids() {
			sb.WriteString("kid(s) ")
		}
	*/
	if p.path == "" {
		sb.WriteString("NOPATH")
	} else {
		sb.WriteString(p.path)
	}
	return (sb.String())
}

func (p *Nord) PrintTree(w io.Writer) error {
	// println("PrintTree: could use printer fn")
	if w == nil {
		return nil
	}
	printTreeTo = w
	e := InspectTree(p, nordPrintOneLiner)
	if e != nil {
		println("nordPrintOneLiner ERR:", e.Error())
		return e
	}
	return nil
}

func (p *Nord) PrintCssTree(w io.Writer) error {
	if w == nil {
		return nil
	}
	printTreeTo = w
	e := InspectTreeWithPreAndPost(p,
		nordPrintCssOneLinerPre, nordPrintCssOneLinerPost)
	if e != nil {
		println("nordPrintCssLine ret'd ERR:", e.Error())
		return e
	}
	return nil
}

func nordPrintOneLiner(p Norder) error {
	// var F StringFunc
	// F = p.GetLineSummaryFunc()
	// fmt.Fprintf(printTreeTo, "%s %s (%T) \n", p.LinePrefixString(), F(p), p)
	// fmt.Fprintf(printTreeTo, "%s %s (%T) \n", p.LinePrefixString(), "?", p)
	fmt.Fprintf(printTreeTo, p.LineSummaryString())
	return nil
}

func nordPrintCssOneLinerPre(p Norder) error {
	// firstEntry := true
	smry := p.LineSummaryString()

	if p.IsDir() {
		// if firstEntry {
		//  } else {
		// <li><details><summary><i>Ice</i> giants</summary>
		// <ul>
		fmt.Fprintf(printTreeTo, "<li><details><summary>"+
			smry+"</summary>\nul")
		// }
	} else {
		// Do both Pre AND Post
		fmt.Fprintf(printTreeTo, "<li>"+smry+"</li>\n")
	}
	fmt.Fprintf(printTreeTo, p.LineSummaryString())
	// firstEntry = false
	return nil
}

func nordPrintCssOneLinerPost(p Norder) error {
	if p.IsDir() {
		fmt.Fprintf(printTreeTo, "</li>\n")
	} /* else {
		// Do both Pre AND Post
		fmt.Fprintf(printTreeTo, "<li>"+smry+"</li>\n")
	} */
	fmt.Fprintf(printTreeTo, p.LineSummaryString())
	return nil
}
