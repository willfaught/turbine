package turbine

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

func NodeString(f *token.FileSet, n ast.Node) (string, error) {
	b := &bytes.Buffer{}
	if err := format.Node(b, f, n); err != nil {
		return "", fmt.Errorf("cannot format node: %v", err)
	}
	return b.String(), nil
}

func NodeSyntax(f *token.FileSet, n ast.Node) Syntax {
	if f == nil {
		f = token.NewFileSet()
	}
	c := nodeConv{fileSet: f}
	return c.node(n.Pos(), n.End(), n)
}

func maybeBlock(s Syntax) *Block {
	if s == nil {
		return nil
	}
	return s.(*Block)
}

func maybeName(s Syntax) *Name {
	if s == nil {
		return nil
	}
	return s.(*Name)
}

type nodeConv struct {
	comments ast.CommentMap
	file     *token.File
	fileSet  *token.FileSet
}

func (c *nodeConv) context2(parent, before, after, n ast.Node) Context {
	var begin, end token.Pos
	var afterGaps, beforeGaps []Gap
	var comment *ast.Comment
	groups := c.comments[n]

	if parent == nil { // File
		begin = c.file.Pos(0)
		end = c.file.Pos(c.file.Size())
	} else {
		if before == nil {
			begin = parent.Pos()
		} else {
			begin = before.End()
		}
		if after == nil {
			end = parent.End()
		} else {
			end = after.Pos()
		}
	}

	// Whitespace Gaps between (not including) left and right positions, where left <= right.
	whitespace := func(left, right token.Pos) (gaps []Gap) {
		if left > right {
			panic(fmt.Sprintf("invalid token positions: %v, %v", left, right))
		}
		// No gap
		if right-left <= 1 {
			return nil
		}
		// Same line
		leftLine, rightLine := c.file.Line(left), c.file.Line(right)
		if leftLine == rightLine {
			return []Gap{&Space{Count: int(right - left - 1)}}
		}
		// First line
		nextLine := leftLine + 1
		next := c.file.LineStart(nextLine)
		if left+2 < next {
			gaps = append(gaps, &Space{Count: int(next - left - 2)})
		}
		gaps = append(gaps, &Line{})
		// Middle lines, if any
		for ; nextLine < rightLine; nextLine++ {
			gaps = append(gaps, &Space{Count: int(c.file.LineStart(nextLine+1) - c.file.LineStart(nextLine) - 2)}, &Line{})
		}
		// Last line
		if next = c.file.LineStart(nextLine); next < right {
			gaps = append(gaps, &Space{Count: int(right - next - 1)})
		}
		return gaps
	}

	beforeGaps = append(beforeGaps, whitespace(begin, c.left(n))...)

	// for j, lines := 0, c.file.Line(c.left(n))-c.file.Line(begin); j < lines; j++ {
	// 	before = append(before, &Line{})
	// }

	i, l := 0, len(groups)
	for ; i < l; i++ { // TODO: but need to tell if first in list! -- for len(groups) > 0 { group := groups[0]; if group ... then break else groups = groups[1:]
		group := groups[i]
		if group.Pos() >= n.Pos() {
			break
		}
		if i > 0 {
			for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End()); i < lines; i++ {
				beforeGaps = append(beforeGaps, &Line{})
			}
		}
		var j int
		for j, comment = range group.List {
			if j > 0 && c.file.Line(group.List[j-1].End()) != c.file.Line(comment.Pos()) {
				beforeGaps = append(beforeGaps, &Line{})
			}
			beforeGaps = append(beforeGaps, c.gap(comment))
		}
	}
	if comment != nil {
		for j, lines := 0, c.file.Line(n.Pos())-c.file.Line(comment.End()); j < lines; j++ {
			beforeGaps = append(beforeGaps, &Line{})
		}
	}
	// After
	if len(groups) > 0 {
		for j, lines := 0, c.file.Line(groups[0].Pos())-c.file.Line(n.End()); j < lines; j++ {
			afterGaps = append(afterGaps, &Line{})
		}
		groups = groups[i:]
		for i, l = 0, len(groups); i < l; i++ {
			group := groups[i]
			if i > 0 {
				for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End()); i < lines; i++ {
					afterGaps = append(afterGaps, &Line{})
				}
			}
			for j, comment := range group.List {
				if j > 0 {
					for i, lines := 0, c.file.Line(group.List[j-1].End())-c.file.Line(comment.Pos()); i < lines; i++ {
						afterGaps = append(afterGaps, &Line{})
					}
				}
				afterGaps = append(afterGaps, c.gap(comment))
			}
		}
	}
	if end != token.NoPos { // TODO: Perhaps not needed if all of node() calls Context()?
		for j, lines := 0, c.file.Line(end)-c.file.Line(c.right(n)-1); j < lines; j++ {
			afterGaps = append(afterGaps, &Line{})
		}
	}
	return Context{After: afterGaps, Before: beforeGaps}
}

func (c *nodeConv) context(begin, end token.Pos, n ast.Node) Context {
	var after, before []Gap
	var comment *ast.Comment
	groups := c.comments[n]
	if begin == token.NoPos {
		// TODO
		panic("no pos")
	}
	if end == token.NoPos {
		// TODO
		panic("no pos")
	}

	// Whitespace Gaps between (not including) left and right positions, where left <= right.
	whitespace := func(left, right token.Pos) (gaps []Gap) {
		if left > right {
			panic(fmt.Sprintf("invalid token positions: %v, %v", left, right))
		}
		// No gap
		if right-left <= 1 {
			return nil
		}
		// Same line
		leftLine, rightLine := c.file.Line(left), c.file.Line(right)
		if leftLine == rightLine {
			return []Gap{&Space{Count: int(right - left - 1)}}
		}
		// First line
		nextLine := leftLine + 1
		next := c.file.LineStart(nextLine)
		if left+2 < next {
			gaps = append(gaps, &Space{Count: int(next - left - 2)})
		}
		gaps = append(gaps, &Line{})
		// Middle lines, if any
		for ; nextLine < rightLine; nextLine++ {
			gaps = append(gaps, &Space{Count: int(c.file.LineStart(nextLine+1) - c.file.LineStart(nextLine) - 2)}, &Line{})
		}
		// Last line
		if next = c.file.LineStart(nextLine); next < right {
			gaps = append(gaps, &Space{Count: int(right - next - 1)})
		}
		return gaps
	}

	before = append(before, whitespace(begin, c.left(n))...)

	// for j, lines := 0, c.file.Line(c.left(n))-c.file.Line(begin); j < lines; j++ {
	// 	before = append(before, &Line{})
	// }

	i, l := 0, len(groups)
	for ; i < l; i++ { // TODO: but need to tell if first in list! -- for len(groups) > 0 { group := groups[0]; if group ... then break else groups = groups[1:]
		group := groups[i]
		if group.Pos() >= n.Pos() {
			break
		}
		if i > 0 {
			for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End()); i < lines; i++ {
				before = append(before, &Line{})
			}
		}
		var j int
		for j, comment = range group.List {
			if j > 0 && c.file.Line(group.List[j-1].End()) != c.file.Line(comment.Pos()) {
				before = append(before, &Line{})
			}
			before = append(before, c.gap(comment))
		}
	}
	if comment != nil {
		for j, lines := 0, c.file.Line(n.Pos())-c.file.Line(comment.End()); j < lines; j++ {
			before = append(before, &Line{})
		}
	}
	// After
	if len(groups) > 0 {
		for j, lines := 0, c.file.Line(groups[0].Pos())-c.file.Line(n.End()); j < lines; j++ {
			after = append(after, &Line{})
		}
		groups = groups[i:]
		for i, l = 0, len(groups); i < l; i++ {
			group := groups[i]
			if i > 0 {
				for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End()); i < lines; i++ {
					after = append(after, &Line{})
				}
			}
			for j, comment := range group.List {
				if j > 0 {
					for i, lines := 0, c.file.Line(group.List[j-1].End())-c.file.Line(comment.Pos()); i < lines; i++ {
						after = append(after, &Line{})
					}
				}
				after = append(after, c.gap(comment))
			}
		}
	}
	if end != token.NoPos { // TODO: Perhaps not needed if all of node() calls Context()?
		for j, lines := 0, c.file.Line(end)-c.file.Line(c.right(n)-1); j < lines; j++ {
			after = append(after, &Line{})
		}
	}
	return Context{After: after, Before: before}
}

func (c *nodeConv) decls(l, r token.Pos, from []ast.Decl) (to []Declaration) {
	if len(from) == 0 {
		return nil
	}
	to = make([]Declaration, len(from))
	for i, d := range from {
		var dl, dr token.Pos
		if i == 0 {
			dl = l
			if len(from) > 1 {
				dr = c.left(from[1])
			} else {
				dr = r
			}
		} else {
			dl = c.left(d)
			if i == len(from)-1 {
				dr = r
			} else {
				dr = c.left(from[i+1])
			}
		}
		to[i] = c.node(dl, dr, d).(Declaration)
	}
	return to
}

func (c *nodeConv) exprs(begin, end token.Pos, from []ast.Expr) (to []Expression) {
	if len(from) == 0 {
		return nil
	}
	to = make([]Expression, len(from))
	for i, f := range from {
		to[i] = c.node(begin, end, f).(Expression)
	}
	return to
}

func (c *nodeConv) idents(l, r token.Pos, from []*ast.Ident) []*Name {
	if len(from) == 0 {
		return nil
	}
	to := make([]*Name, len(from))
	for i, f := range from {
		to[i] = c.node(l, r, f).(*Name)
	}
	return to
}

func (c *nodeConv) left(n ast.Node) token.Pos {
	if cgs := c.comments[n]; len(cgs) > 0 {
		if p := cgs[0].List[0].Pos(); p < n.Pos() {
			return p
		}
	}
	return n.Pos()
}

func (c *nodeConv) right(n ast.Node) token.Pos {
	if cgs := c.comments[n]; len(cgs) > 0 {
		cs := cgs[len(cgs)-1].List
		if p := cs[len(cs)-1].End(); p > n.End() {
			return p
		}
	}
	return n.End()
}

func (c *nodeConv) afterComments(n ast.Node) []*ast.CommentGroup {
	groups := c.comments[n]
	for i, g := range groups {
		if g.Pos() > n.Pos() {
			return groups[i:]
		}
	}
	return nil
}

func (c *nodeConv) beforeComments(n ast.Node) []*ast.CommentGroup {
	groups := c.comments[n]
	for i, g := range groups {
		if g.Pos() > n.Pos() {
			return groups[:i]
		}
	}
	return groups
}

func (c *nodeConv) gap(from ast.Node) (to Gap) {
	switch from := from.(type) {
	case *ast.Comment:
		to = &Comment{Text: from.Text}
	default:
		panic(fmt.Sprintf("invalid gap: %#v", from))
	}
	return to
}

func (c *nodeConv) node2(parent, before, after, from ast.Node) (to Syntax) {
	switch from := from.(type) {
	// case *ast.AssignStmt:
	// 	return &Assign{
	// 		// Left: c.exprs(n.Lhs),
	// 		// TODO: Operator: n.Tok,
	// 		// Right: c.exprs(n.Rhs),
	// 	}
	// case *ast.BadStmt:
	// 	return nil // TODO
	// case *ast.BlockStmt:
	// 	if n == nil {
	// 		return nil
	// 	}
	// 	s := &Block{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.List = c.stmts(n.Lbrace+1, n.Rbrace, n.List)
	// 	return s
	// case *ast.BranchStmt:
	// 	switch n.Tok {
	// 	case token.BREAK:
	// 		return &Break{
	// 			Label: maybeName(c.node(token.NoPos, token.NoPos, n.Label)),
	// 		}
	// 	case token.CONTINUE:
	// 		return &Continue{
	// 			Label: maybeName(c.node(token.NoPos, token.NoPos, n.Label)),
	// 		}
	// 	case token.FALLTHROUGH:
	// 		return &Fallthrough{}
	// 	case token.GOTO:
	// 		return &Goto{
	// 			Label: c.node(token.NoPos, token.NoPos, n.Label).(*Name),
	// 		}
	// 	}
	// case *ast.DeclStmt:
	// 	return c.node(token.NoPos, token.NoPos, n.Decl)
	// case *ast.DeferStmt:
	// 	return &Defer{
	// 		Call: c.node(token.NoPos, token.NoPos, n.Call).(*Call),
	// 	}
	// case *ast.EmptyStmt:
	// 	// TODO:
	// 	// return &Empty{}
	// 	panic("empty stmt unexpected")
	// case *ast.ExprStmt:
	// 	return c.node(begin, end, n.X)
	// case *ast.ForStmt:
	// 	return &For{
	// 		// TODO:
	// 		// Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Cond: c.node(token.NoPos, token.NoPos, n.Cond),
	// 		// Post: c.node(token.NoPos, token.NoPos, n.Post),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.GoStmt:
	// 	return &Go{
	// 		Call: c.node(token.NoPos, token.NoPos, n.Call).(*Call),
	// 	}
	// case *ast.IfStmt:
	// 	return &If{
	// 		Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		Cond: c.node(token.NoPos, token.NoPos, n.Cond),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		Else: c.node(token.NoPos, token.NoPos, n.Else),
	// 	}
	// case *ast.IncDecStmt:
	// 	if n.Tok == token.INC {
	// 		return &Inc{
	// 			// TODO: X: c.node(token.NoPos, token.NoPos, n.X),
	// 		}
	// 	}
	// 	return &Dec{
	// 		// TODO: X: c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.LabeledStmt:
	// 	return &Label{
	// 		Label: c.node(token.NoPos, token.NoPos, n.Label).(*Name),
	// 		// TODO: Stmt:  c.node(token.NoPos, token.NoPos, n.Stmt),
	// 	}
	// case *ast.RangeStmt:
	// 	return &Range{
	// 		Assign: n.Tok == token.ASSIGN,
	// 		// TODO:
	// 		// Key:    c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value:  c.node(token.NoPos, token.NoPos, n.Value),
	// 		// X:      c.node(token.NoPos, token.NoPos, n.X),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.ReturnStmt:
	// 	return &Return{
	// 		Results: c.exprs(n.Results),
	// 	}
	// case *ast.SelectStmt:
	// 	return &Select{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.SendStmt:
	// 	// TODO:
	// 	// return &Send{
	// 	// 	Chan:  c.node(token.NoPos, token.NoPos, n.Chan),
	// 	// 	Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	// }
	// case *ast.SwitchStmt:
	// 	return &Switch{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		// TODO:
	// 		// Init:  c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Tag),
	// 	}
	// case *ast.TypeSwitchStmt:
	// 	return &Switch{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		// TODO:
	// 		// Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Type: c.node(token.NoPos, token.NoPos, n.Assign),
	// 	}
	// case *ast.ImportSpec:
	// 	imp := &Import{Context: c.context2(parent, before, after, from)}
	// 	if from.Name == nil {
	// 		imp.Path = c.node(from.Pos(), from.End(), from.Path).(*String)
	// 	} else {
	// 		pathLeft := c.left(from.Path)
	// 		imp.Name = c.node(from.Pos(), pathLeft, from.Name).(*Name)
	// 		imp.Path = c.node(pathLeft, from.End(), from.Path).(*String)
	// 	}
	// 	to = imp
	// case *ast.TypeSpec:
	// 	return &Type{
	// 		// Assign: n.Assign,
	// 		Name: c.node(token.NoPos, token.NoPos, n.Name).(*Name),
	// 		// TODO: Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 	}
	// case *ast.ValueSpec:
	// 	return &Const{
	// 		Names: c.idents(n.Names),
	// 		// TODO:
	// 		// Type:   c.node(token.NoPos, token.NoPos, n.Type),
	// 		// Values: c.exprs(n.Values),
	// 	}
	// case *ast.BadDecl:
	// 	return nil
	// case *ast.FuncDecl:
	// 	s := &Func{}
	// 	s.Context = c.context(begin, end, n)
	// 	if n.Recv != nil {
	// 		s.Receiver = c.node(add(n.Type.Func, lenFunc), c.nodeBegin(n.Name), n.Recv).(*FieldList)
	// 		s.Name = c.node(c.nodeBegin(n.Name), c.nodeBegin(n.Type), n.Name).(*Name)
	// 	} else {
	// 		s.Name = c.node(add(n.Type.Func, lenFunc), c.nodeBegin(n.Type), n.Name).(*Name)
	// 	}
	// 	var p token.Pos
	// 	if n.Type.Results == nil && n.Body == nil {
	// 		p = c.nodeEnd(n.Type.Params)
	// 	} else if n.Type.Results != nil {
	// 		p = c.nodeBegin(n.Type.Results)
	// 	} else {
	// 		p = c.nodeBegin(n.Body)
	// 	}
	// 	s.Parameters = c.node(c.nodeBegin(n.Type.Params), p, n.Type.Params).(*FieldList)
	// 	if n.Type.Results != nil {
	// 		if n.Body == nil {
	// 			p = c.nodeEnd(n.Type.Results)
	// 		} else {
	// 			p = c.nodeBegin(n.Body)
	// 		}
	// 		s.Results = c.node(n.Type.Results.Pos(), p, n.Type.Results).(*FieldList)
	// 	}
	// 	if n.Body != nil {
	// 		s.Body = c.node(c.nodeBegin(n.Body), c.nodeEnd(n.Body), n.Body).(*Block)
	// 	}
	// 	return s
	// case *ast.GenDecl:
	// 	m := c.context(begin, end, n)
	// 	// TODO: Capture comments and lines between n.Tok and n.Lparen
	// 	if n.Lparen == token.NoPos {
	// 		begin = n.TokPos + token.Pos(len(n.Tok.String()))
	// 		end = n.End()
	// 	} else {
	// 		begin = n.Lparen + 1
	// 		end = n.Rparen
	// 	}
	// 	var ss []Syntax
	// 	for i, spec := range n.Specs {
	// 		var s Syntax
	// 		var specBegin, specEnd token.Pos
	// 		if i == 0 {
	// 			specBegin = begin
	// 			if len(n.Specs) > 1 {
	// 				specEnd = c.nodeBegin(n.Specs[1])
	// 			} else {
	// 				specEnd = end
	// 			}
	// 		} else {
	// 			specBegin = c.nodeBegin(spec)
	// 			if i == len(n.Specs)-1 {
	// 				specEnd = end
	// 			} else {
	// 				specEnd = c.nodeBegin(n.Specs[i+1])
	// 			}
	// 		}
	// 		// TODO: Capture lines between specs
	// 		switch spec := spec.(type) {
	// 		case *ast.ImportSpec:
	// 			im := &Import{}
	// 			im.Context = c.context(specBegin, specEnd, spec)
	// 			if spec.Name != nil {
	// 				im.Name = c.node(c.nodeBegin(spec.Name), c.nodeBegin(spec.Path), spec.Name).(*Name)
	// 			}
	// 			im.Path = c.node(c.nodeBegin(spec.Path), c.nodeEnd(spec.Path), spec.Path).(*String)
	// 			s = im
	// 		case *ast.TypeSpec:
	// 			s = &Type{
	// 				Context: c.context(specBegin, specEnd, spec),
	// 				Name:   c.node(token.NoPos, token.NoPos, spec.Name).(*Name),
	// 				// Assign: spec.Assign,
	// 				// TODO: Type: c.node(token.NoPos, token.NoPos, spec.Type),
	// 			}
	// 		case *ast.ValueSpec:
	// 			switch n.Tok {
	// 			case token.CONST:
	// 				s = &Const{
	// 					Context: c.context(specBegin, specEnd, spec),
	// 					Names:  c.idents(spec.Names),
	// 					// TODO:
	// 					// Type:   c.node(token.NoPos, token.NoPos, spec.Type),
	// 					// Values: c.exprs(spec.Values),
	// 				}
	// 			case token.VAR:
	// 				s = &Var{
	// 					Context: c.context(specBegin, specEnd, spec),
	// 					Names:  c.idents(spec.Names),
	// 					// TODO:
	// 					// Type:   c.node(token.NoPos, token.NoPos, spec.Type),
	// 					// Values: c.exprs(spec.Values),
	// 				}
	// 			default:
	// 				panic(n.Tok)
	// 			}
	// 		}
	// 		ss = append(ss, s)
	// 	}
	// 	if n.Lparen == token.NoPos {
	// 		return ss[0]
	// 	}
	// 	switch n.Tok {
	// 	case token.CONST:
	// 		// TODO: return &ConstList{List: ss, Context: m}
	// 	case token.IMPORT:
	// 		// TODO: return &ImportList{List: ss, Context: m}
	// 	case token.TYPE:
	// 		// TODO: return &TypeList{List: ss, Context: m}
	// 	case token.VAR:
	// 		// TODO: return &VarList{List: ss, Context: m}
	// 	}
	// case *ast.ArrayType:
	// 	return &Array{
	// 		// Element: c.node(token.NoPos, token.NoPos, n.Elt),
	// 		// Length:  c.node(token.NoPos, token.NoPos, n.Len),
	// 	}
	// case *ast.BadExpr:
	// 	return nil
	// case *ast.BasicLit:
	// 	if n == nil {
	// 		return nil
	// 	}
	// 	switch n.Kind {
	// 	case token.CHAR:
	// 		return &Rune{
	// 			Text: n.Value,
	// 		}
	// 	case token.FLOAT:
	// 		return &Float{
	// 			Text: n.Value,
	// 		}
	// 	case token.IMAG:
	// 		return &Imag{
	// 			Text: n.Value,
	// 		}
	// 	case token.INT:
	// 		return &Int{
	// 			Text: n.Value,
	// 		}
	// 	case token.STRING:
	// 		return &String{
	// 			Text: n.Value,
	// 		}
	// 	default:
	// 		panic(n) // TODO
	// 	}
	// case *ast.BinaryExpr:
	// 	// TODO:
	// 	// return &Binary{
	// 	// 	Operator: n.Op,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// 	Y:        c.node(token.NoPos, token.NoPos, n.Y),
	// 	// }
	// case *ast.CallExpr:
	// 	s := &Call{}
	// 	s.Context = c.context(begin, end, n)
	// 	// TODO:
	// 	// s.Fun = c.node(c.nodeBegin(n.Fun), n.Lparen, n.Fun)
	// 	// s.Args = c.exprs(n.Args)
	// 	s.Ellipsis = n.Ellipsis != token.NoPos
	// 	return s
	// case *ast.ChanType:
	// 	switch n.Dir {
	// 	case ast.RECV:
	// 		return &ChanIn{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	case ast.SEND:
	// 		return &ChanOut{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	default:
	// 		return &Chan{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	}
	// case *ast.CompositeLit:
	// 	return &Composite{
	// 		// TODO: Elts: c.exprs(n.Elts),
	// 		// TODO: Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 	}
	// case *ast.Ellipsis:
	// 	// TODO:
	// 	// return &Ellipsis{
	// 	// 	Elt: c.node(token.NoPos, token.NoPos, n.Elt),
	// 	// }
	// case *ast.FuncLit:
	// 	return &Func{
	// 		Body:       c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		Parameters: c.node(token.NoPos, token.NoPos, n.Type.Params).(*FieldList),
	// 		Results:    c.node(token.NoPos, token.NoPos, n.Type.Results).(*FieldList),
	// 	}
	// case *ast.FuncType:
	// 	return &Func{
	// 		Parameters: c.node(token.NoPos, token.NoPos, n.Params).(*FieldList),
	// 		Results:    c.node(token.NoPos, token.NoPos, n.Results).(*FieldList),
	// 	}
	case *ast.Ident:
		return &Name{
			Context: c.context2(parent, before, after, from),
			Text:    from.Name,
		}
	// case *ast.IndexExpr:
	// 	return &Index{
	// 		// TODO:
	// 		// Index: c.node(token.NoPos, token.NoPos, n.Index),
	// 		// X:     c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.InterfaceType:
	// 	return &Interface{
	// 		Methods: c.node(token.NoPos, token.NoPos, n.Methods).(*FieldList),
	// 	}
	// case *ast.KeyValueExpr:
	// 	return &KeyValue{
	// 		// TODO:
	// 		// Key:   c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	}
	// case *ast.MapType:
	// 	return &Map{
	// 		// TODO:
	// 		// Key:   c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	}
	// case *ast.ParenExpr:
	// 	return &Paren{
	// 		X: c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.SelectorExpr:
	// 	return &Selector{
	// 		Sel: c.node(token.NoPos, token.NoPos, n.Sel).(*Name),
	// 		X:   c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.SliceExpr:
	// 	return &Slice{
	// 		High: c.node(token.NoPos, token.NoPos, n.High),
	// 		Low:  c.node(token.NoPos, token.NoPos, n.Low),
	// 		Max:  c.node(token.NoPos, token.NoPos, n.Max),
	// 		X:    c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.StarExpr:
	// 	// TODO
	// 	// return &Unary{
	// 	// 	Operator: token.MUL,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// }
	// case *ast.StructType:
	// 	return &Struct{
	// 		Fields: c.node(token.NoPos, token.NoPos, n.Fields).(*FieldList),
	// 	}
	// case *ast.TypeAssertExpr:
	// 	return &Assert{
	// 		// Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 		// X:    c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.UnaryExpr:
	// 	// TODO
	// 	// return &Unary{
	// 	// 	Operator: n.Op,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// }
	// case *ast.CaseClause:
	// 	return &Case{
	// 		Body: c.stmts(token.NoPos, token.NoPos, n.Body),
	// 		List: c.exprs(n.List),
	// 	}
	// case *ast.CommClause:
	// 	return &Case{
	// 		Body: c.stmts(token.NoPos, token.NoPos, n.Body),
	// 		Comm: c.node(token.NoPos, token.NoPos, n.Comm),
	// 	}
	// case *ast.CommentGroup:
	// 	var cs []*Comment
	// 	for _, com := range n.List {
	// 		cs = append(cs, c.node(token.NoPos, token.NoPos, com).(*Comment))
	// 	}
	// 	return &CommentGroup{
	// 		List: cs,
	// 	}
	// case *ast.Field:
	// 	s := &Field{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.Names = c.idents(n.Names)
	// 	s.Type = c.node(token.NoPos, token.NoPos, n.Type)
	// 	if n.Tag != nil {
	// 		s.Tag = c.node(token.NoPos, token.NoPos, n.Tag).(*String)
	// 	}
	// 	return s
	// case *ast.FieldList:
	// 	if n == nil {
	// 		return (*FieldList)(nil)
	// 	}
	// 	s := &FieldList{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.List = make([]*Field, len(n.List))
	// 	for i, f := range n.List {
	// 		var fieldBegin, fieldEnd token.Pos
	// 		if i == 0 {
	// 			fieldBegin = begin
	// 			if len(n.List) > 1 {
	// 				fieldEnd = c.nodeBegin(n.List[1])
	// 			} else {
	// 				fieldEnd = end
	// 			}
	// 		} else {
	// 			fieldBegin = c.nodeBegin(f)
	// 			if i == len(n.List)-1 {
	// 				fieldEnd = end
	// 			} else {
	// 				fieldEnd = c.nodeBegin(n.List[i+1])
	// 			}
	// 		}
	// 		s.List[i] = c.node(fieldBegin, fieldEnd, f).(*Field)
	// 	}
	// 	return s
	case *ast.File:
		// TODO
		// doc
		// pkg name
		// imports
		// decls
		c.comments = ast.NewCommentMap(c.fileSet, from, from.Comments)
		c.file = c.fileSet.File(from.Pos())
		file := &File{}
		file.Context = c.context2(nil, nil, nil, from)
		var after ast.Node
		if len(from.Decls) > 0 {
			after = from.Decls[0]
		}
		file.Package = c.node2(from, nil, after, from.Name).(*Name)
		// if len(from.Decls) > 0 {
		// 	file.Decls = c.decls(p, from.End(), from.Decls)
		// }
		to = file
		// case *ast.Package:
		// 	var fs map[string]*File
		// 	if n.Files != nil {
		// 		fs = map[string]*File{}
		// 		for k, v := range n.Files {
		// 			fs[k] = c.node(token.NoPos, token.NoPos, v).(*File)
		// 		}
		// 	}
		// 	return &Package{
		// 		Files: fs,
		// 	}
	default:
		panic(fmt.Sprintf("invalid node: %#v", from))
	}
	return to
}

func (c *nodeConv) node(l, r token.Pos, from ast.Node) (to Syntax) {
	switch from := from.(type) {
	// case *ast.AssignStmt:
	// 	return &Assign{
	// 		// Left: c.exprs(n.Lhs),
	// 		// TODO: Operator: n.Tok,
	// 		// Right: c.exprs(n.Rhs),
	// 	}
	// case *ast.BadStmt:
	// 	return nil // TODO
	// case *ast.BlockStmt:
	// 	if n == nil {
	// 		return nil
	// 	}
	// 	s := &Block{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.List = c.stmts(n.Lbrace+1, n.Rbrace, n.List)
	// 	return s
	// case *ast.BranchStmt:
	// 	switch n.Tok {
	// 	case token.BREAK:
	// 		return &Break{
	// 			Label: maybeName(c.node(token.NoPos, token.NoPos, n.Label)),
	// 		}
	// 	case token.CONTINUE:
	// 		return &Continue{
	// 			Label: maybeName(c.node(token.NoPos, token.NoPos, n.Label)),
	// 		}
	// 	case token.FALLTHROUGH:
	// 		return &Fallthrough{}
	// 	case token.GOTO:
	// 		return &Goto{
	// 			Label: c.node(token.NoPos, token.NoPos, n.Label).(*Name),
	// 		}
	// 	}
	// case *ast.DeclStmt:
	// 	return c.node(token.NoPos, token.NoPos, n.Decl)
	// case *ast.DeferStmt:
	// 	return &Defer{
	// 		Call: c.node(token.NoPos, token.NoPos, n.Call).(*Call),
	// 	}
	// case *ast.EmptyStmt:
	// 	// TODO:
	// 	// return &Empty{}
	// 	panic("empty stmt unexpected")
	// case *ast.ExprStmt:
	// 	return c.node(begin, end, n.X)
	// case *ast.ForStmt:
	// 	return &For{
	// 		// TODO:
	// 		// Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Cond: c.node(token.NoPos, token.NoPos, n.Cond),
	// 		// Post: c.node(token.NoPos, token.NoPos, n.Post),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.GoStmt:
	// 	return &Go{
	// 		Call: c.node(token.NoPos, token.NoPos, n.Call).(*Call),
	// 	}
	// case *ast.IfStmt:
	// 	return &If{
	// 		Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		Cond: c.node(token.NoPos, token.NoPos, n.Cond),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		Else: c.node(token.NoPos, token.NoPos, n.Else),
	// 	}
	// case *ast.IncDecStmt:
	// 	if n.Tok == token.INC {
	// 		return &Inc{
	// 			// TODO: X: c.node(token.NoPos, token.NoPos, n.X),
	// 		}
	// 	}
	// 	return &Dec{
	// 		// TODO: X: c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.LabeledStmt:
	// 	return &Label{
	// 		Label: c.node(token.NoPos, token.NoPos, n.Label).(*Name),
	// 		// TODO: Stmt:  c.node(token.NoPos, token.NoPos, n.Stmt),
	// 	}
	// case *ast.RangeStmt:
	// 	return &Range{
	// 		Assign: n.Tok == token.ASSIGN,
	// 		// TODO:
	// 		// Key:    c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value:  c.node(token.NoPos, token.NoPos, n.Value),
	// 		// X:      c.node(token.NoPos, token.NoPos, n.X),
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.ReturnStmt:
	// 	return &Return{
	// 		Results: c.exprs(n.Results),
	// 	}
	// case *ast.SelectStmt:
	// 	return &Select{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 	}
	// case *ast.SendStmt:
	// 	// TODO:
	// 	// return &Send{
	// 	// 	Chan:  c.node(token.NoPos, token.NoPos, n.Chan),
	// 	// 	Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	// }
	// case *ast.SwitchStmt:
	// 	return &Switch{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		// TODO:
	// 		// Init:  c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Tag),
	// 	}
	// case *ast.TypeSwitchStmt:
	// 	return &Switch{
	// 		Body: c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		// TODO:
	// 		// Init: c.node(token.NoPos, token.NoPos, n.Init),
	// 		// Type: c.node(token.NoPos, token.NoPos, n.Assign),
	// 	}
	case *ast.ImportSpec:
		imp := &Import{Context: c.context(l, r, from)}
		if from.Name == nil {
			imp.Path = c.node(from.Pos(), from.End(), from.Path).(*String)
		} else {
			pathLeft := c.left(from.Path)
			imp.Name = c.node(from.Pos(), pathLeft, from.Name).(*Name)
			imp.Path = c.node(pathLeft, from.End(), from.Path).(*String)
		}
		to = imp
	// case *ast.TypeSpec:
	// 	return &Type{
	// 		// Assign: n.Assign,
	// 		Name: c.node(token.NoPos, token.NoPos, n.Name).(*Name),
	// 		// TODO: Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 	}
	// case *ast.ValueSpec:
	// 	return &Const{
	// 		Names: c.idents(n.Names),
	// 		// TODO:
	// 		// Type:   c.node(token.NoPos, token.NoPos, n.Type),
	// 		// Values: c.exprs(n.Values),
	// 	}
	// case *ast.BadDecl:
	// 	return nil
	// case *ast.FuncDecl:
	// 	s := &Func{}
	// 	s.Context = c.context(begin, end, n)
	// 	if n.Recv != nil {
	// 		s.Receiver = c.node(add(n.Type.Func, lenFunc), c.nodeBegin(n.Name), n.Recv).(*FieldList)
	// 		s.Name = c.node(c.nodeBegin(n.Name), c.nodeBegin(n.Type), n.Name).(*Name)
	// 	} else {
	// 		s.Name = c.node(add(n.Type.Func, lenFunc), c.nodeBegin(n.Type), n.Name).(*Name)
	// 	}
	// 	var p token.Pos
	// 	if n.Type.Results == nil && n.Body == nil {
	// 		p = c.nodeEnd(n.Type.Params)
	// 	} else if n.Type.Results != nil {
	// 		p = c.nodeBegin(n.Type.Results)
	// 	} else {
	// 		p = c.nodeBegin(n.Body)
	// 	}
	// 	s.Parameters = c.node(c.nodeBegin(n.Type.Params), p, n.Type.Params).(*FieldList)
	// 	if n.Type.Results != nil {
	// 		if n.Body == nil {
	// 			p = c.nodeEnd(n.Type.Results)
	// 		} else {
	// 			p = c.nodeBegin(n.Body)
	// 		}
	// 		s.Results = c.node(n.Type.Results.Pos(), p, n.Type.Results).(*FieldList)
	// 	}
	// 	if n.Body != nil {
	// 		s.Body = c.node(c.nodeBegin(n.Body), c.nodeEnd(n.Body), n.Body).(*Block)
	// 	}
	// 	return s
	// case *ast.GenDecl:
	// 	m := c.context(begin, end, n)
	// 	// TODO: Capture comments and lines between n.Tok and n.Lparen
	// 	if n.Lparen == token.NoPos {
	// 		begin = n.TokPos + token.Pos(len(n.Tok.String()))
	// 		end = n.End()
	// 	} else {
	// 		begin = n.Lparen + 1
	// 		end = n.Rparen
	// 	}
	// 	var ss []Syntax
	// 	for i, spec := range n.Specs {
	// 		var s Syntax
	// 		var specBegin, specEnd token.Pos
	// 		if i == 0 {
	// 			specBegin = begin
	// 			if len(n.Specs) > 1 {
	// 				specEnd = c.nodeBegin(n.Specs[1])
	// 			} else {
	// 				specEnd = end
	// 			}
	// 		} else {
	// 			specBegin = c.nodeBegin(spec)
	// 			if i == len(n.Specs)-1 {
	// 				specEnd = end
	// 			} else {
	// 				specEnd = c.nodeBegin(n.Specs[i+1])
	// 			}
	// 		}
	// 		// TODO: Capture lines between specs
	// 		switch spec := spec.(type) {
	// 		case *ast.ImportSpec:
	// 			im := &Import{}
	// 			im.Context = c.context(specBegin, specEnd, spec)
	// 			if spec.Name != nil {
	// 				im.Name = c.node(c.nodeBegin(spec.Name), c.nodeBegin(spec.Path), spec.Name).(*Name)
	// 			}
	// 			im.Path = c.node(c.nodeBegin(spec.Path), c.nodeEnd(spec.Path), spec.Path).(*String)
	// 			s = im
	// 		case *ast.TypeSpec:
	// 			s = &Type{
	// 				Context: c.context(specBegin, specEnd, spec),
	// 				Name:   c.node(token.NoPos, token.NoPos, spec.Name).(*Name),
	// 				// Assign: spec.Assign,
	// 				// TODO: Type: c.node(token.NoPos, token.NoPos, spec.Type),
	// 			}
	// 		case *ast.ValueSpec:
	// 			switch n.Tok {
	// 			case token.CONST:
	// 				s = &Const{
	// 					Context: c.context(specBegin, specEnd, spec),
	// 					Names:  c.idents(spec.Names),
	// 					// TODO:
	// 					// Type:   c.node(token.NoPos, token.NoPos, spec.Type),
	// 					// Values: c.exprs(spec.Values),
	// 				}
	// 			case token.VAR:
	// 				s = &Var{
	// 					Context: c.context(specBegin, specEnd, spec),
	// 					Names:  c.idents(spec.Names),
	// 					// TODO:
	// 					// Type:   c.node(token.NoPos, token.NoPos, spec.Type),
	// 					// Values: c.exprs(spec.Values),
	// 				}
	// 			default:
	// 				panic(n.Tok)
	// 			}
	// 		}
	// 		ss = append(ss, s)
	// 	}
	// 	if n.Lparen == token.NoPos {
	// 		return ss[0]
	// 	}
	// 	switch n.Tok {
	// 	case token.CONST:
	// 		// TODO: return &ConstList{List: ss, Context: m}
	// 	case token.IMPORT:
	// 		// TODO: return &ImportList{List: ss, Context: m}
	// 	case token.TYPE:
	// 		// TODO: return &TypeList{List: ss, Context: m}
	// 	case token.VAR:
	// 		// TODO: return &VarList{List: ss, Context: m}
	// 	}
	// case *ast.ArrayType:
	// 	return &Array{
	// 		// Element: c.node(token.NoPos, token.NoPos, n.Elt),
	// 		// Length:  c.node(token.NoPos, token.NoPos, n.Len),
	// 	}
	// case *ast.BadExpr:
	// 	return nil
	// case *ast.BasicLit:
	// 	if n == nil {
	// 		return nil
	// 	}
	// 	switch n.Kind {
	// 	case token.CHAR:
	// 		return &Rune{
	// 			Text: n.Value,
	// 		}
	// 	case token.FLOAT:
	// 		return &Float{
	// 			Text: n.Value,
	// 		}
	// 	case token.IMAG:
	// 		return &Imag{
	// 			Text: n.Value,
	// 		}
	// 	case token.INT:
	// 		return &Int{
	// 			Text: n.Value,
	// 		}
	// 	case token.STRING:
	// 		return &String{
	// 			Text: n.Value,
	// 		}
	// 	default:
	// 		panic(n) // TODO
	// 	}
	// case *ast.BinaryExpr:
	// 	// TODO:
	// 	// return &Binary{
	// 	// 	Operator: n.Op,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// 	Y:        c.node(token.NoPos, token.NoPos, n.Y),
	// 	// }
	// case *ast.CallExpr:
	// 	s := &Call{}
	// 	s.Context = c.context(begin, end, n)
	// 	// TODO:
	// 	// s.Fun = c.node(c.nodeBegin(n.Fun), n.Lparen, n.Fun)
	// 	// s.Args = c.exprs(n.Args)
	// 	s.Ellipsis = n.Ellipsis != token.NoPos
	// 	return s
	// case *ast.ChanType:
	// 	switch n.Dir {
	// 	case ast.RECV:
	// 		return &ChanIn{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	case ast.SEND:
	// 		return &ChanOut{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	default:
	// 		return &Chan{
	// 			// TODO: Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 		}
	// 	}
	// case *ast.CompositeLit:
	// 	return &Composite{
	// 		// TODO: Elts: c.exprs(n.Elts),
	// 		// TODO: Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 	}
	// case *ast.Ellipsis:
	// 	// TODO:
	// 	// return &Ellipsis{
	// 	// 	Elt: c.node(token.NoPos, token.NoPos, n.Elt),
	// 	// }
	// case *ast.FuncLit:
	// 	return &Func{
	// 		Body:       c.node(token.NoPos, token.NoPos, n.Body).(*Block),
	// 		Parameters: c.node(token.NoPos, token.NoPos, n.Type.Params).(*FieldList),
	// 		Results:    c.node(token.NoPos, token.NoPos, n.Type.Results).(*FieldList),
	// 	}
	// case *ast.FuncType:
	// 	return &Func{
	// 		Parameters: c.node(token.NoPos, token.NoPos, n.Params).(*FieldList),
	// 		Results:    c.node(token.NoPos, token.NoPos, n.Results).(*FieldList),
	// 	}
	case *ast.Ident:
		return &Name{
			Context: c.context(l, r, from),
			Text:    from.Name,
		}
	// case *ast.IndexExpr:
	// 	return &Index{
	// 		// TODO:
	// 		// Index: c.node(token.NoPos, token.NoPos, n.Index),
	// 		// X:     c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.InterfaceType:
	// 	return &Interface{
	// 		Methods: c.node(token.NoPos, token.NoPos, n.Methods).(*FieldList),
	// 	}
	// case *ast.KeyValueExpr:
	// 	return &KeyValue{
	// 		// TODO:
	// 		// Key:   c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	}
	// case *ast.MapType:
	// 	return &Map{
	// 		// TODO:
	// 		// Key:   c.node(token.NoPos, token.NoPos, n.Key),
	// 		// Value: c.node(token.NoPos, token.NoPos, n.Value),
	// 	}
	// case *ast.ParenExpr:
	// 	return &Paren{
	// 		X: c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.SelectorExpr:
	// 	return &Selector{
	// 		Sel: c.node(token.NoPos, token.NoPos, n.Sel).(*Name),
	// 		X:   c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.SliceExpr:
	// 	return &Slice{
	// 		High: c.node(token.NoPos, token.NoPos, n.High),
	// 		Low:  c.node(token.NoPos, token.NoPos, n.Low),
	// 		Max:  c.node(token.NoPos, token.NoPos, n.Max),
	// 		X:    c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.StarExpr:
	// 	// TODO
	// 	// return &Unary{
	// 	// 	Operator: token.MUL,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// }
	// case *ast.StructType:
	// 	return &Struct{
	// 		Fields: c.node(token.NoPos, token.NoPos, n.Fields).(*FieldList),
	// 	}
	// case *ast.TypeAssertExpr:
	// 	return &Assert{
	// 		// Type: c.node(token.NoPos, token.NoPos, n.Type),
	// 		// X:    c.node(token.NoPos, token.NoPos, n.X),
	// 	}
	// case *ast.UnaryExpr:
	// 	// TODO
	// 	// return &Unary{
	// 	// 	Operator: n.Op,
	// 	// 	X:        c.node(token.NoPos, token.NoPos, n.X),
	// 	// }
	// case *ast.CaseClause:
	// 	return &Case{
	// 		Body: c.stmts(token.NoPos, token.NoPos, n.Body),
	// 		List: c.exprs(n.List),
	// 	}
	// case *ast.CommClause:
	// 	return &Case{
	// 		Body: c.stmts(token.NoPos, token.NoPos, n.Body),
	// 		Comm: c.node(token.NoPos, token.NoPos, n.Comm),
	// 	}
	// case *ast.CommentGroup:
	// 	var cs []*Comment
	// 	for _, com := range n.List {
	// 		cs = append(cs, c.node(token.NoPos, token.NoPos, com).(*Comment))
	// 	}
	// 	return &CommentGroup{
	// 		List: cs,
	// 	}
	// case *ast.Field:
	// 	s := &Field{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.Names = c.idents(n.Names)
	// 	s.Type = c.node(token.NoPos, token.NoPos, n.Type)
	// 	if n.Tag != nil {
	// 		s.Tag = c.node(token.NoPos, token.NoPos, n.Tag).(*String)
	// 	}
	// 	return s
	// case *ast.FieldList:
	// 	if n == nil {
	// 		return (*FieldList)(nil)
	// 	}
	// 	s := &FieldList{}
	// 	s.Context = c.context(begin, end, n)
	// 	s.List = make([]*Field, len(n.List))
	// 	for i, f := range n.List {
	// 		var fieldBegin, fieldEnd token.Pos
	// 		if i == 0 {
	// 			fieldBegin = begin
	// 			if len(n.List) > 1 {
	// 				fieldEnd = c.nodeBegin(n.List[1])
	// 			} else {
	// 				fieldEnd = end
	// 			}
	// 		} else {
	// 			fieldBegin = c.nodeBegin(f)
	// 			if i == len(n.List)-1 {
	// 				fieldEnd = end
	// 			} else {
	// 				fieldEnd = c.nodeBegin(n.List[i+1])
	// 			}
	// 		}
	// 		s.List[i] = c.node(fieldBegin, fieldEnd, f).(*Field)
	// 	}
	// 	return s
	case *ast.File:
		// TODO
		// doc
		// pkg name
		// imports
		// decls
		c.comments = ast.NewCommentMap(c.fileSet, from, from.Comments)
		c.file = c.fileSet.File(from.Pos())
		file := &File{}
		file.Context = c.context(c.file.Pos(0), c.file.Pos(c.file.Size()), from)
		var p token.Pos
		if len(from.Decls) == 0 {
			p = from.End()
		} else {
			p = c.left(from.Decls[0])
		}
		file.Package = c.node(add(from.Package, lenPackage), p, from.Name).(*Name)
		if len(from.Decls) > 0 {
			file.Decls = c.decls(p, from.End(), from.Decls)
		}
		to = file
		// case *ast.Package:
		// 	var fs map[string]*File
		// 	if n.Files != nil {
		// 		fs = map[string]*File{}
		// 		for k, v := range n.Files {
		// 			fs[k] = c.node(token.NoPos, token.NoPos, v).(*File)
		// 		}
		// 	}
		// 	return &Package{
		// 		Files: fs,
		// 	}
	default:
		panic(fmt.Sprintf("invalid node: %#v", from))
	}
	return to
}

func add(p token.Pos, n int) token.Pos {
	return p + token.Pos(n)
}

func (c *nodeConv) specs(l, r token.Pos, from []ast.Spec) []Syntax {
	if len(from) == 0 {
		return nil
	}
	to := make([]Syntax, len(from))
	for i, f := range from {
		to[i] = c.node(token.NoPos, token.NoPos, f)
	}
	return to
}

func (c *nodeConv) stmts(begin, end token.Pos, from []ast.Stmt) []Statement {
	if len(from) == 0 {
		return nil
	}
	to := make([]Statement, len(from))
	for i, f := range from {
		var stmtBegin, stmtEnd token.Pos
		if i == 0 {
			stmtBegin = begin
			if len(from) > 1 {
				stmtEnd = c.left(from[1])
			} else {
				stmtEnd = end
			}
		} else {
			stmtBegin = c.left(f)
			if i == len(from)-1 {
				stmtEnd = end
			} else {
				stmtEnd = c.left(from[i+1])
			}
		}
		to[i] = c.node(stmtBegin, stmtEnd, f).(Statement)
	}
	return to
}
