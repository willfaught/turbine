package syntax

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
		for j, lines := 0, c.file.Line(end-1)-c.file.Line(c.nodeEnd(n)-1); j < lines; j++ { // end-1 because end is a go/ast.Node.End
			after = append(after, &Line{})
		}
	}
	return Markup{After: after, Before: before}
}

func (c *nodeConv) node(n ast.Node) Syntax {
	return c.nodePos(token.NoPos, token.NoPos, n)
}

func (c *nodeConv) nodePos(begin, end token.Pos, n ast.Node) Syntax {
	switch n := n.(type) {
	case nil:
		return nil
	case *ast.AssignStmt:
		return &Assign{
			Left:     c.exprs(n.Lhs),
			Operator: n.Tok,
			Right:    c.exprs(n.Rhs),
		}
	case *ast.BadStmt:
		return nil // TODO
	case *ast.BlockStmt:
		if n == nil {
			return nil
		}
		return &Block{
			List: c.stmts(n.List),
		}
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
		return &Empty{}
	case *ast.ExprStmt:
		return c.node(n.X)
	case *ast.ForStmt:
		return &For{
			Init: c.node(n.Init),
			Cond: c.node(n.Cond),
			Post: c.node(n.Post),
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
				X: c.node(n.X),
			}
		}
		return &Dec{
			X: c.node(n.X),
		}
	case *ast.LabeledStmt:
		return &Label{
			Label: c.node(n.Label).(*Name),
			Stmt:  c.node(n.Stmt),
		}
	case *ast.RangeStmt:
		return &Range{
			Assign: n.Tok == token.ASSIGN,
			Key:    c.node(n.Key),
			Value:  c.node(n.Value),
			X:      c.node(n.X),
			Body:   c.node(n.Body).(*Block),
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
		return &Send{
			Chan:  c.node(n.Chan),
			Value: c.node(n.Value),
		}
	case *ast.SwitchStmt:
		return &Switch{
			Body:  c.node(n.Body).(*Block),
			Init:  c.node(n.Init),
			Value: c.node(n.Tag),
		}
	case *ast.TypeSwitchStmt:
		return &Switch{
			Body: c.node(n.Body).(*Block),
			Init: c.node(n.Init),
			Type: c.node(n.Assign),
		}
	case *ast.ImportSpec:
		return &Import{
			Name: c.node(n.Name).(*Name),
			Path: c.node(n.Path).(*String),
		}
	case *ast.TypeSpec:
		return &Type{
			Assign: n.Assign,
			Name:   c.node(n.Name).(*Name),
			Type:   c.node(n.Type),
		}
	case *ast.ValueSpec:
		return &Const{
			Names:  c.idents(n.Names),
			Type:   c.node(n.Type),
			Values: c.exprs(n.Values),
		}
	case *ast.BadDecl:
		return nil
	case *ast.FuncDecl:
		return &Func{
			Body:       maybeBlock(c.node(n.Body)),
			Name:       c.node(n.Name).(*Name),
			Parameters: c.node(n.Type.Params).(*FieldList),
			Receiver:   c.node(n.Recv).(*FieldList),
			Results:    c.node(n.Type.Results).(*FieldList),
		}
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
			var syn Syntax
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
					im.Name = c.nodePos(spec.Name.Pos(), spec.Path.Pos(), spec.Name).(*Name)
				}
				im.Path = c.nodePos(c.nodeBegin(spec.Path), c.nodeEnd(spec.Path), spec.Path).(*String)
				syn = im
			case *ast.TypeSpec:
				syn = &Type{
					Assign: spec.Assign,
					Markup: c.markup(specBegin, specEnd, spec),
					Name:   c.node(spec.Name).(*Name),
					Type:   c.node(spec.Type),
				}
			case *ast.ValueSpec:
				switch n.Tok {
				case token.CONST:
					syn = &Const{
						Markup: c.markup(specBegin, specEnd, spec),
						Names:  c.idents(spec.Names),
						Type:   c.node(spec.Type),
						Values: c.exprs(spec.Values),
					}
				case token.VAR:
					syn = &Var{
						Markup: c.markup(specBegin, specEnd, spec),
						Names:  c.idents(spec.Names),
						Type:   c.node(spec.Type),
						Values: c.exprs(spec.Values),
					}
				}
			}
			ss = append(ss, syn)
		}
		if n.Lparen == token.NoPos {
			return ss[0]
		}
		switch n.Tok {
		case token.CONST:
			return &ConstList{List: ss, Markup: m}
		case token.IMPORT:
			return &ImportList{List: ss, Markup: m}
		case token.TYPE:
			return &TypeList{List: ss, Markup: m}
		case token.VAR:
			return &VarList{List: ss, Markup: m}
		}
	case *ast.ArrayType:
		return &Array{
			Element: c.node(n.Elt),
			Length:  c.node(n.Len),
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
		return &Binary{
			Operator: n.Op,
			X:        c.node(n.X),
			Y:        c.node(n.Y),
		}
	case *ast.CallExpr:
		return &Call{
			Args:     c.exprs(n.Args),
			Ellipsis: n.Ellipsis != token.NoPos,
			Fun:      c.node(n.Fun),
		}
	case *ast.ChanType:
		switch n.Dir {
		case ast.RECV:
			return &ChanIn{
				Value: c.node(n.Value),
			}
		case ast.SEND:
			return &ChanOut{
				Value: c.node(n.Value),
			}
		default:
			return &Chan{
				Value: c.node(n.Value),
			}
		}
	case *ast.CompositeLit:
		return &Composite{
			Elts: c.exprs(n.Elts),
			Type: c.node(n.Type),
		}
	case *ast.Ellipsis:
		return &Ellipsis{
			Elt: c.node(n.Elt),
		}
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
			Text: n.Name,
		}
	case *ast.IndexExpr:
		return &Index{
			Index: c.node(n.Index),
			X:     c.node(n.X),
		}
	case *ast.InterfaceType:
		return &Interface{
			Methods: c.node(n.Methods).(*FieldList),
		}
	case *ast.KeyValueExpr:
		return &KeyValue{
			Key:   c.node(n.Key),
			Value: c.node(n.Value),
		}
	case *ast.MapType:
		return &Map{
			Key:   c.node(n.Key),
			Value: c.node(n.Value),
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
		return &Unary{
			Operator: token.MUL,
			X:        c.node(n.X),
		}
	case *ast.StructType:
		return &Struct{
			Fields: c.node(n.Fields).(*FieldList),
		}
	case *ast.TypeAssertExpr:
		return &Assert{
			Type: c.node(n.Type),
			X:    c.node(n.X),
		}
	case *ast.UnaryExpr:
		return &Unary{
			Operator: n.Op,
			X:        c.node(n.X),
		}
	case *ast.CaseClause:
		return &Case{
			Body: c.stmts(n.Body),
			List: c.exprs(n.List),
		}
	case *ast.CommClause:
		return &Case{
			Body: c.stmts(n.Body),
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
		var tag *String
		if b, ok := c.node(n.Tag).(*String); ok {
			tag = b
		}
		return &Field{
			Names: c.idents(n.Names),
			Tag:   tag,
			Type:  c.node(n.Type),
		}
	case *ast.FieldList:
		if n == nil {
			return (*FieldList)(nil)
		}
		var fs []*Field
		for _, f := range n.List {
			fs = append(fs, c.node(f).(*Field))
		}
		return &FieldList{
			List: fs,
		}
	case *ast.File:
		// TODO
		// doc
		// pkg name
		// imports
		// decls
		c.comments = ast.NewCommentMap(c.tokens, n, n.Comments)
		c.file = c.tokens.File(n.Pos())
		var x token.Pos
		if l := len(n.Decls); l > 0 {
			x = c.nodeBegin(n.Decls[0])
		} else {
			x = n.End()
		}
		return &File{
			Markup:  c.markup(c.file.Pos(0), c.file.Pos(c.file.Size()), n),
			Package: c.nodePos(add(n.Package, lenPackage), x, n.Name).(*Name),
			Decls:   c.decls(x, n.End(), n.Decls),
		}
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

func (c *nodeConv) stmts(from []ast.Stmt) []Syntax {
	if len(from) == 0 {
		return nil
	}
	to := make([]Syntax, len(from))
	for i, f := range from {
		to[i] = c.node(f)
	}
	return to
}
