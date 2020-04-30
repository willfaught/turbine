package turbine

/*
import (
	"go/ast"
	"go/token"
)

func Convert(f *token.FileSet, n ast.Node) Syntax {
	if f == nil {
		f = token.NewFileSet()
	}
	c := nodeConv{tokens: f}
	return c.nodePos(n.Pos(), n.End(), n)
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
	tokens   *token.FileSet
}

func (c *nodeConv) decls(begin, end token.Pos, from []ast.Decl) []Syntax {
	if len(from) == 0 {
		return nil
	}
	to := make([]Syntax, len(from))
	for i, f := range from {
		var declBegin, declEnd token.Pos
		if i == 0 {
			declBegin = begin
			if len(from) > 1 {
				declEnd = c.nodeBegin(from[1])
			} else {
				declEnd = end
			}
		} else {
			declBegin = c.nodeBegin(f)
			if i == len(from)-1 {
				declEnd = end
			} else {
				declEnd = c.nodeBegin(from[i+1])
			}
		}
		to[i] = c.nodePos(declBegin, declEnd, f)
	}
	return to
}

func (c *nodeConv) exprs(from []ast.Expr) []Syntax {
	return c.exprsPos(0, 0, from)
}

func (c *nodeConv) exprsPos(begin, end token.Pos, from []ast.Expr) []Syntax {
	if len(from) == 0 {
		return nil
	}
	to := make([]Syntax, len(from))
	for i, f := range from {
		to[i] = c.node(f)
	}
	return to
}

func (c *nodeConv) idents(from []*ast.Ident) []*Name {
	if len(from) == 0 {
		return nil
	}
	to := make([]*Name, len(from))
	for i, f := range from {
		to[i] = c.node(f).(*Name)
	}
	return to
}

func (c *nodeConv) nodeBegin(n ast.Node) token.Pos {
	if cgs := c.comments[n]; len(cgs) > 0 {
		if p := cgs[0].List[0].Pos(); p < n.Pos() {
			return p
		}
	}
	return n.Pos()
}

func (c *nodeConv) nodeEnd(n ast.Node) token.Pos {
	if cgs := c.comments[n]; len(cgs) > 0 {
		cs := cgs[len(cgs)-1].List
		if p := cs[len(cs)-1].End(); p > n.End() {
			return p
		}
	}
	return n.End()
}

func (c *nodeConv) markup(begin, end token.Pos, n ast.Node) Markup {
	var after, before []Syntax
	var comment *ast.Comment
	groups := c.comments[n]
	if begin != token.NoPos { // TODO: Perhaps not needed if all of nodePos() calls markup()?
		for j, lines := 0, c.file.Line(c.nodeBegin(n))-c.file.Line(begin); j < lines; j++ {
			before = append(before, &Line{})
		}
	}
	i, l := 0, len(groups)
	for ; i < l; i++ { // TODO: but need to tell if first in list! -- for len(groups) > 0 { group := groups[0]; if group ... then break else groups = groups[1:]
		group := groups[i]
		if group.Pos() >= n.Pos() {
			break
		}
		if i > 0 {
			before = append(before, &Line{})
			for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End())-1; i < lines; i++ {
				before = append(before, &Line{})
			}
		}
		var j int
		for j, comment = range group.List {
			if j > 0 && c.file.Line(group.List[j-1].End()) != c.file.Line(comment.Pos()) {
				before = append(before, &Line{})
			}
			before = append(before, c.node(comment))
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
				after = append(after, &Line{})
				for i, lines := 0, c.file.Line(group.Pos())-c.file.Line(groups[i-1].End())-1; i < lines; i++ {
					after = append(after, &Line{})
				}
			}
			var j int
			for j, comment = range group.List {
				if j > 0 && c.file.Line(group.List[j-1].End()) != c.file.Line(comment.Pos()) {
					after = append(after, &Line{})
				}
				after = append(after, c.node(comment), &Line{})
			}
		}
	}
	if end != token.NoPos { // TODO: Perhaps not needed if all of nodePos() calls markup()?
		for j, lines := 0, c.file.Line(end)-c.file.Line(c.nodeEnd(n)-1); j < lines; j++ {
			after = append(after, &Line{})
		}
	}
	return Markup{After: after, Before: before}
}

func (c *nodeConv) node(n ast.Node) Syntax {
	return c.nodePos(token.NoPos, token.NoPos, n)
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

func (c *nodeConv) nodePos(begin, end token.Pos, n ast.Node) Syntax {
	switch n := n.(type) {
	case nil:
		return nil
	case *ast.AssignStmt:
		return &Assign{
			// Left: c.exprs(n.Lhs),
			// TODO: Operator: n.Tok,
			// Right: c.exprs(n.Rhs),
		}
	case *ast.BadStmt:
		return nil // TODO
	case *ast.BlockStmt:
		if n == nil {
			return nil
		}
		s := &Block{}
		s.Markup = c.markup(begin, end, n)
		s.List = c.stmts(n.Lbrace+1, n.Rbrace, n.List)
		return s
	case *ast.BranchStmt:
		switch n.Tok {
		case token.BREAK:
			return &Break{
				Label: maybeName(c.node(n.Label)),
			}
		case token.CONTINUE:
			return &Continue{
				Label: maybeName(c.node(n.Label)),
			}
		case token.FALLTHROUGH:
			return &Fallthrough{}
		case token.GOTO:
			return &Goto{
				Label: c.node(n.Label).(*Name),
			}
		}
	case *ast.DeclStmt:
		return c.node(n.Decl)
	case *ast.DeferStmt:
		return &Defer{
			Call: c.node(n.Call).(*Call),
		}
	case *ast.EmptyStmt:
		// TODO:
		// return &Empty{}
		panic("empty stmt unexpected")
	case *ast.ExprStmt:
		return c.nodePos(begin, end, n.X)
	case *ast.ForStmt:
		return &For{
			// TODO:
			// Init: c.node(n.Init),
			// Cond: c.node(n.Cond),
			// Post: c.node(n.Post),
			Body: c.node(n.Body).(*Block),
		}
	case *ast.GoStmt:
		return &Go{
			Call: c.node(n.Call).(*Call),
		}
	case *ast.IfStmt:
		return &If{
			Init: c.node(n.Init),
			Cond: c.node(n.Cond),
			Body: c.node(n.Body).(*Block),
			Else: c.node(n.Else),
		}
	case *ast.IncDecStmt:
		if n.Tok == token.INC {
			return &Inc{
				// TODO: X: c.node(n.X),
			}
		}
		return &Dec{
			// TODO: X: c.node(n.X),
		}
	case *ast.LabeledStmt:
		return &Label{
			Label: c.node(n.Label).(*Name),
			// TODO: Stmt:  c.node(n.Stmt),
		}
	case *ast.RangeStmt:
		return &Range{
			Assign: n.Tok == token.ASSIGN,
			// TODO:
			// Key:    c.node(n.Key),
			// Value:  c.node(n.Value),
			// X:      c.node(n.X),
			Body: c.node(n.Body).(*Block),
		}
	case *ast.ReturnStmt:
		return &Return{
			Results: c.exprs(n.Results),
		}
	case *ast.SelectStmt:
		return &Select{
			Body: c.node(n.Body).(*Block),
		}
	case *ast.SendStmt:
		// TODO:
		// return &Send{
		// 	Chan:  c.node(n.Chan),
		// 	Value: c.node(n.Value),
		// }
	case *ast.SwitchStmt:
		return &Switch{
			Body: c.node(n.Body).(*Block),
			// TODO:
			// Init:  c.node(n.Init),
			// Value: c.node(n.Tag),
		}
	case *ast.TypeSwitchStmt:
		return &Switch{
			Body: c.node(n.Body).(*Block),
			// TODO:
			// Init: c.node(n.Init),
			// Type: c.node(n.Assign),
		}
	case *ast.ImportSpec:
		return &Import{
			Name: c.node(n.Name).(*Name),
			Path: c.node(n.Path).(*String),
		}
	case *ast.TypeSpec:
		return &Type{
			// Assign: n.Assign,
			Name: c.node(n.Name).(*Name),
			// TODO: Type: c.node(n.Type),
		}
	case *ast.ValueSpec:
		return &Const{
			Names: c.idents(n.Names),
			// TODO:
			// Type:   c.node(n.Type),
			// Values: c.exprs(n.Values),
		}
	case *ast.BadDecl:
		return nil
	case *ast.FuncDecl:
		s := &Func{}
		s.Markup = c.markup(begin, end, n)
		if n.Recv != nil {
			s.Receiver = c.nodePos(add(n.Type.Func, lenFunc), c.nodeBegin(n.Name), n.Recv).(*FieldList)
			s.Name = c.nodePos(c.nodeBegin(n.Name), c.nodeBegin(n.Type), n.Name).(*Name)
		} else {
			s.Name = c.nodePos(add(n.Type.Func, lenFunc), c.nodeBegin(n.Type), n.Name).(*Name)
		}
		var p token.Pos
		if n.Type.Results == nil && n.Body == nil {
			p = c.nodeEnd(n.Type.Params)
		} else if n.Type.Results != nil {
			p = c.nodeBegin(n.Type.Results)
		} else {
			p = c.nodeBegin(n.Body)
		}
		s.Parameters = c.nodePos(c.nodeBegin(n.Type.Params), p, n.Type.Params).(*FieldList)
		if n.Type.Results != nil {
			if n.Body == nil {
				p = c.nodeEnd(n.Type.Results)
			} else {
				p = c.nodeBegin(n.Body)
			}
			s.Results = c.nodePos(n.Type.Results.Pos(), p, n.Type.Results).(*FieldList)
		}
		if n.Body != nil {
			s.Body = c.nodePos(c.nodeBegin(n.Body), c.nodeEnd(n.Body), n.Body).(*Block)
		}
		return s
	case *ast.GenDecl:
		m := c.markup(begin, end, n)
		// TODO: Capture comments and lines between n.Tok and n.Lparen
		if n.Lparen == token.NoPos {
			begin = n.TokPos + token.Pos(len(n.Tok.String()))
			end = n.End()
		} else {
			begin = n.Lparen + 1
			end = n.Rparen
		}
		var ss []Syntax
		for i, spec := range n.Specs {
			var s Syntax
			var specBegin, specEnd token.Pos
			if i == 0 {
				specBegin = begin
				if len(n.Specs) > 1 {
					specEnd = c.nodeBegin(n.Specs[1])
				} else {
					specEnd = end
				}
			} else {
				specBegin = c.nodeBegin(spec)
				if i == len(n.Specs)-1 {
					specEnd = end
				} else {
					specEnd = c.nodeBegin(n.Specs[i+1])
				}
			}
			// TODO: Capture lines between specs
			switch spec := spec.(type) {
			case *ast.ImportSpec:
				im := &Import{}
				im.Markup = c.markup(specBegin, specEnd, spec)
				if spec.Name != nil {
					im.Name = c.nodePos(c.nodeBegin(spec.Name), c.nodeBegin(spec.Path), spec.Name).(*Name)
				}
				im.Path = c.nodePos(c.nodeBegin(spec.Path), c.nodeEnd(spec.Path), spec.Path).(*String)
				s = im
			case *ast.TypeSpec:
				s = &Type{
					Markup: c.markup(specBegin, specEnd, spec),
					Name:   c.node(spec.Name).(*Name),
					// Assign: spec.Assign,
					// TODO: Type: c.node(spec.Type),
				}
			case *ast.ValueSpec:
				switch n.Tok {
				case token.CONST:
					s = &Const{
						Markup: c.markup(specBegin, specEnd, spec),
						Names:  c.idents(spec.Names),
						// TODO:
						// Type:   c.node(spec.Type),
						// Values: c.exprs(spec.Values),
					}
				case token.VAR:
					s = &Var{
						Markup: c.markup(specBegin, specEnd, spec),
						Names:  c.idents(spec.Names),
						// TODO:
						// Type:   c.node(spec.Type),
						// Values: c.exprs(spec.Values),
					}
				default:
					panic(n.Tok)
				}
			}
			ss = append(ss, s)
		}
		if n.Lparen == token.NoPos {
			return ss[0]
		}
		switch n.Tok {
		case token.CONST:
			// TODO: return &ConstList{List: ss, Markup: m}
		case token.IMPORT:
			// TODO: return &ImportList{List: ss, Markup: m}
		case token.TYPE:
			// TODO: return &TypeList{List: ss, Markup: m}
		case token.VAR:
			// TODO: return &VarList{List: ss, Markup: m}
		}
	case *ast.ArrayType:
		return &Array{
			// Element: c.node(n.Elt),
			// Length:  c.node(n.Len),
		}
	case *ast.BadExpr:
		return nil
	case *ast.BasicLit:
		if n == nil {
			return nil
		}
		switch n.Kind {
		case token.CHAR:
			return &Rune{
				Text: n.Value,
			}
		case token.FLOAT:
			return &Float{
				Text: n.Value,
			}
		case token.IMAG:
			return &Imag{
				Text: n.Value,
			}
		case token.INT:
			return &Int{
				Text: n.Value,
			}
		case token.STRING:
			return &String{
				Text: n.Value,
			}
		default:
			panic(n) // TODO
		}
	case *ast.BinaryExpr:
		// TODO:
		// return &Binary{
		// 	Operator: n.Op,
		// 	X:        c.node(n.X),
		// 	Y:        c.node(n.Y),
		// }
	case *ast.CallExpr:
		s := &Call{}
		s.Markup = c.markup(begin, end, n)
		// TODO:
		// s.Fun = c.nodePos(c.nodeBegin(n.Fun), n.Lparen, n.Fun)
		// s.Args = c.exprs(n.Args)
		s.Ellipsis = n.Ellipsis != token.NoPos
		return s
	case *ast.ChanType:
		switch n.Dir {
		case ast.RECV:
			return &ChanIn{
				// TODO: Value: c.node(n.Value),
			}
		case ast.SEND:
			return &ChanOut{
				// TODO: Value: c.node(n.Value),
			}
		default:
			return &Chan{
				// TODO: Value: c.node(n.Value),
			}
		}
	case *ast.CompositeLit:
		return &Composite{
			// TODO: Elts: c.exprs(n.Elts),
			// TODO: Type: c.node(n.Type),
		}
	case *ast.Ellipsis:
		// TODO:
		// return &Ellipsis{
		// 	Elt: c.node(n.Elt),
		// }
	case *ast.FuncLit:
		return &Func{
			Body:       c.node(n.Body).(*Block),
			Parameters: c.node(n.Type.Params).(*FieldList),
			Results:    c.node(n.Type.Results).(*FieldList),
		}
	case *ast.FuncType:
		return &Func{
			Parameters: c.node(n.Params).(*FieldList),
			Results:    c.node(n.Results).(*FieldList),
		}
	case *ast.Ident:
		if n == nil {
			return nil
		}
		return &Name{
			Markup: c.markup(begin, end, n),
			Text:   n.Name,
		}
	case *ast.IndexExpr:
		return &Index{
			// TODO:
			// Index: c.node(n.Index),
			// X:     c.node(n.X),
		}
	case *ast.InterfaceType:
		return &Interface{
			Methods: c.node(n.Methods).(*FieldList),
		}
	case *ast.KeyValueExpr:
		return &KeyValue{
			// TODO:
			// Key:   c.node(n.Key),
			// Value: c.node(n.Value),
		}
	case *ast.MapType:
		return &Map{
			// TODO:
			// Key:   c.node(n.Key),
			// Value: c.node(n.Value),
		}
	case *ast.ParenExpr:
		return &Paren{
			X: c.node(n.X),
		}
	case *ast.SelectorExpr:
		return &Selector{
			Sel: c.node(n.Sel).(*Name),
			X:   c.node(n.X),
		}
	case *ast.SliceExpr:
		return &Slice{
			High: c.node(n.High),
			Low:  c.node(n.Low),
			Max:  c.node(n.Max),
			X:    c.node(n.X),
		}
	case *ast.StarExpr:
		// TODO
		// return &Unary{
		// 	Operator: token.MUL,
		// 	X:        c.node(n.X),
		// }
	case *ast.StructType:
		return &Struct{
			Fields: c.node(n.Fields).(*FieldList),
		}
	case *ast.TypeAssertExpr:
		return &Assert{
			// Type: c.node(n.Type),
			// X:    c.node(n.X),
		}
	case *ast.UnaryExpr:
		// TODO
		// return &Unary{
		// 	Operator: n.Op,
		// 	X:        c.node(n.X),
		// }
	case *ast.CaseClause:
		return &Case{
			Body: c.stmts(token.NoPos, token.NoPos, n.Body),
			List: c.exprs(n.List),
		}
	case *ast.CommClause:
		return &Case{
			Body: c.stmts(token.NoPos, token.NoPos, n.Body),
			Comm: c.node(n.Comm),
		}
	case *ast.Comment:
		return &Comment{
			Text: n.Text,
		}
	case *ast.CommentGroup:
		var cs []*Comment
		for _, com := range n.List {
			cs = append(cs, c.node(com).(*Comment))
		}
		return &CommentGroup{
			List: cs,
		}
	case *ast.Field:
		s := &Field{}
		s.Markup = c.markup(begin, end, n)
		s.Names = c.idents(n.Names)
		s.Type = c.node(n.Type)
		if n.Tag != nil {
			s.Tag = c.node(n.Tag).(*String)
		}
		return s
	case *ast.FieldList:
		if n == nil {
			return (*FieldList)(nil)
		}
		s := &FieldList{}
		s.Markup = c.markup(begin, end, n)
		s.List = make([]*Field, len(n.List))
		for i, f := range n.List {
			var fieldBegin, fieldEnd token.Pos
			if i == 0 {
				fieldBegin = begin
				if len(n.List) > 1 {
					fieldEnd = c.nodeBegin(n.List[1])
				} else {
					fieldEnd = end
				}
			} else {
				fieldBegin = c.nodeBegin(f)
				if i == len(n.List)-1 {
					fieldEnd = end
				} else {
					fieldEnd = c.nodeBegin(n.List[i+1])
				}
			}
			s.List[i] = c.nodePos(fieldBegin, fieldEnd, f).(*Field)
		}
		return s
	case *ast.File:
		// TODO
		// doc
		// pkg name
		// imports
		// decls
		c.comments = ast.NewCommentMap(c.tokens, n, n.Comments)
		c.file = c.tokens.File(n.Pos())
		s := &File{}
		s.Markup = c.markup(c.file.Pos(0), c.file.Pos(c.file.Size()), n)
		var p token.Pos
		if len(n.Decls) > 0 {
			p = c.nodeBegin(n.Decls[0])
		} else {
			p = c.nodeEnd(n.Name)
		}
		s.Package = c.nodePos(add(n.Package, lenPackage), p, n.Name).(*Name)
		if len(n.Decls) > 0 {
			s.Decls = c.decls(p, c.nodeEnd(n.Decls[len(n.Decls)-1]), n.Decls)
		}
		return s
	case *ast.Package:
		var fs map[string]*File
		if n.Files != nil {
			fs = map[string]*File{}
			for k, v := range n.Files {
				fs[k] = c.node(v).(*File)
			}
		}
		return &Package{
			Files: fs,
		}
	}
	return nil
}

func add(p token.Pos, n int) token.Pos {
	return p + token.Pos(n)
}

func (c *nodeConv) specs(from []ast.Spec) []Syntax {
	if len(from) == 0 {
		return nil
	}
	to := make([]Syntax, len(from))
	for i, f := range from {
		to[i] = c.node(f)
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
				stmtEnd = c.nodeBegin(from[1])
			} else {
				stmtEnd = end
			}
		} else {
			stmtBegin = c.nodeBegin(f)
			if i == len(from)-1 {
				stmtEnd = end
			} else {
				stmtEnd = c.nodeBegin(from[i+1])
			}
		}
		to[i] = c.nodePos(stmtBegin, stmtEnd, f)
	}
	return to
}
*/
