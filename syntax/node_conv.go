package syntax

import (
	"go/ast"
	"go/token"
)

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
	comment  int
	comments ast.CommentMap
	file     *ast.File
	tokens   *token.FileSet
}

func (c *nodeConv) decls(from []ast.Decl) []Syntax {
	var to []Syntax
	for _, f := range from {
		to = append(to, c.node(f))
	}
	return to
}

func (c *nodeConv) exprs(from []ast.Expr) []Syntax {
	var to []Syntax
	for _, f := range from {
		to = append(to, c.node(f))
	}
	return to
}

func (c *nodeConv) idents(from []*ast.Ident) []*Name {
	var to []*Name
	for _, f := range from {
		to = append(to, c.node(f).(*Name))
	}
	return to
}

func (c *nodeConv) node(n ast.Node) Syntax {
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
		var ss []Syntax
		for _, spec := range n.Specs {
			var syn Syntax
			switch spec := spec.(type) {
			case *ast.ImportSpec:
				syn = &Import{
					Name: maybeName(c.node(spec.Name)),
					Path: c.node(spec.Path).(*String),
				}
			case *ast.TypeSpec:
				syn = &Type{
					Assign: spec.Assign,
					Name:   c.node(spec.Name).(*Name),
					Type:   c.node(spec.Type),
				}
			case *ast.ValueSpec:
				switch n.Tok {
				case token.CONST:
					syn = &Const{
						Names:  c.idents(spec.Names),
						Type:   c.node(spec.Type),
						Values: c.exprs(spec.Values),
					}
				case token.VAR:
					syn = &Var{
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
			return &ConstList{
				List: ss,
			}
		case token.IMPORT:
			return &ImportList{
				List: ss,
			}
		case token.TYPE:
			return &TypeList{
				List: ss,
			}
		case token.VAR:
			return &VarList{
				List: ss,
			}
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
		case token.CHAR:
			return &Rune{
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
		// doc
		// pkg name
		// imports
		// decls
		c.comments = ast.NewCommentMap(c.tokens, n, n.Comments)
		c.file = n
		return &File{
			Name:  c.node(n.Name).(*Name),
			Decls: c.decls(n.Decls),
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

func (c *nodeConv) assignComments(ln, rn ast.Node, lm, rm *Markup) {
	if c.file == nil || len(c.file.Comments) == 0 {
		return
	}
	for {
		if c.comment >= len(c.file.Comments) {
			break
		}
		var next = c.file.Comments[c.comment]
		if ln != nil && ln.Pos() >= next.Pos() {
			break
		}
		if rn != nil && next.Pos() >= rn.Pos() {
			break
		}
		var found bool
		for _, com := range c.comments[ln] {
			if com == next {
				found = true
			}
		}
		if found {

		}
	}
}

func (c *nodeConv) specs(from []ast.Spec) []Syntax {
	var to []Syntax
	for _, f := range from {
		to = append(to, c.node(f))
	}
	return to
}

func (c *nodeConv) stmts(from []ast.Stmt) []Syntax {
	var to []Syntax
	for _, f := range from {
		to = append(to, c.node(f))
	}
	return to
}
