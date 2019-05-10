package syntax

import (
	"fmt"
	"go/ast"
	"go/token"
)

var (
	lenArrow     = len(token.ARROW.String())
	lenChan      = len(token.CHAN.String())
	lenColon     = len(token.COLON.String())
	lenConst     = len(token.CONST.String())
	lenEllipsis  = len(token.ELLIPSIS.String())
	lenFunc      = len(token.FUNC.String())
	lenImport    = len(token.IMPORT.String())
	lenInterface = len(token.INTERFACE.String())
	lenLbrace    = len(token.LBRACE.String())
	lenLbrack    = len(token.LBRACK.String())
	lenLparen    = len(token.LPAREN.String())
	lenMap       = len(token.MAP.String())
	lenMul       = len(token.MUL.String())
	lenPeriod    = len(token.PERIOD.String())
	lenRbrace    = len(token.RBRACE.String())
	lenRbrack    = len(token.RBRACK.String())
	lenRparen    = len(token.RPAREN.String())
	lenStruct    = len(token.STRUCT.String())
	lenType      = len(token.TYPE.String())
	lenVar       = len(token.VAR.String())
)

func blockStmt(s ast.Stmt) *ast.BlockStmt {
	if s == nil {
		return nil
	}
	return s.(*ast.BlockStmt)
}

func ident(e ast.Expr) *ast.Ident {
	if e == nil {
		return nil
	}
	return e.(*ast.Ident)
}

type syntaxConv struct {
	astFile      *ast.File
	end          token.Pos
	tokenFile    *token.File
	tokenFileSet *token.FileSet
}

func (c *syntaxConv) decl(s Syntax) ast.Decl {
	switch s := s.(type) {
	case nil:
		return nil
	case *Const:
		return &ast.GenDecl{
			Tok:    token.CONST,
			TokPos: c.next(lenConst),
			Specs:  []ast.Spec{c.spec(s)},
		}
	case *ConstList:
		return &ast.GenDecl{
			Tok:    token.CONST,
			TokPos: c.next(lenConst),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(s.List),
			Rparen: c.next(lenRparen),
		}
	case *Func:
		return &ast.FuncDecl{
			Body: blockStmt(c.stmt(s.Body)),
			Name: c.expr(s.Name).(*ast.Ident),
			Recv: c.node(s.Receiver).(*ast.FieldList),
			Type: &ast.FuncType{
				Params:  c.node(s.Parameters).(*ast.FieldList),
				Results: c.node(s.Results).(*ast.FieldList),
			},
		}
	case *Import:
		return &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: c.next(lenImport),
			Specs:  []ast.Spec{c.spec(s)},
		}
	case *ImportList:
		return &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: c.next(lenImport),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(s.List),
			Rparen: c.next(lenRparen),
		}
	case *Type:
		return &ast.GenDecl{
			Tok:    token.TYPE,
			TokPos: c.next(lenType),
			Specs:  []ast.Spec{c.spec(s)},
		}
	case *TypeList:
		return &ast.GenDecl{
			Tok:    token.TYPE,
			TokPos: c.next(lenType),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(s.List),
			Rparen: c.next(lenRparen),
		}
	case *Var:
		return &ast.GenDecl{
			Tok:    token.VAR,
			TokPos: c.next(lenVar),
			Specs:  []ast.Spec{c.spec(s)},
		}
	case *VarList:
		return &ast.GenDecl{
			Tok:    token.VAR,
			TokPos: c.next(lenVar),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(s.List),
			Rparen: c.next(lenRparen),
		}
	default:
		return nil
	}
}

func (c *syntaxConv) decls(ss []Syntax) []ast.Decl {
	var ds []ast.Decl
	for _, s := range ss {
		ds = append(ds, c.decl(s))
	}
	return ds
}

func (c *syntaxConv) expr(s Syntax) (e ast.Expr) {
	switch s := s.(type) {
	case nil:
		return nil
	case *Array:
		c.markup(s.Before)
		e = &ast.ArrayType{
			Lbrack: c.next(lenLbrack),
			Len:    c.expr(s.Length),
		}
		c.skip(lenRbrack)
		e.(*ast.ArrayType).Elt = c.expr(s.Element)
		c.markup(s.After)
	case *Float:
		c.markup(s.Before)
		e = &ast.BasicLit{
			ValuePos: c.next(len(s.Text)),
			Kind:     token.FLOAT,
			Value:    s.Text,
		}
		c.markup(s.After)
	case *Imag:
		c.markup(s.Before)
		e = &ast.BasicLit{
			ValuePos: c.next(len(s.Text)),
			Kind:     token.IMAG,
			Value:    s.Text,
		}
		c.markup(s.After)
	case *Int:
		c.markup(s.Before)
		e = &ast.BasicLit{
			ValuePos: c.next(len(s.Text)),
			Kind:     token.INT,
			Value:    s.Text,
		}
		c.markup(s.After)
	case *Rune:
		c.markup(s.Before)
		e = &ast.BasicLit{
			ValuePos: c.next(len(s.Text)),
			Kind:     token.CHAR,
			Value:    s.Text,
		}
		c.markup(s.After)
	case *String:
		if s == nil {
			return nil
		}
		c.markup(s.Before)
		e = &ast.BasicLit{
			ValuePos: c.next(len(s.Text)),
			Kind:     token.STRING,
			Value:    s.Text,
		}
		c.markup(s.After)
	case *Binary:
		c.markup(s.Before)
		e = &ast.BinaryExpr{
			X:     c.expr(s.X),
			OpPos: c.next(len(s.Operator.String())),
			Op:    s.Operator,
			Y:     c.expr(s.Y),
		}
		c.markup(s.After)
	case *Call:
		c.markup(s.Before)
		e = &ast.CallExpr{
			Fun:      c.expr(s.Fun),
			Lparen:   c.next(lenLparen),
			Args:     c.exprs(s.Args),
			Ellipsis: 0, // TODO
			Rparen:   c.next(lenRparen),
		}
		c.markup(s.After)
	case *Chan:
		c.markup(s.Before)
		e = &ast.ChanType{
			Begin: c.next(lenChan + 1),
			Value: c.expr(s.Value),
		}
		c.markup(s.After)
	case *ChanIn:
		c.markup(s.Before)
		e = &ast.ChanType{
			Begin: c.next(lenChan),
			Arrow: c.next(lenArrow + 1),
			Dir:   ast.RECV,
			Value: c.expr(s.Value),
		}
		c.markup(s.After)
	case *ChanOut:
		c.markup(s.Before)
		var p = c.next(lenChan + lenArrow + 1)
		e = &ast.ChanType{
			Begin: p,
			Arrow: p,
			Dir:   ast.SEND,
			Value: c.expr(s.Value),
		}
		c.markup(s.After)
	case *Composite:
		c.markup(s.Before)
		e = &ast.CompositeLit{
			Type:   c.expr(s.Type),
			Lbrace: c.next(lenLbrace),
			Elts:   c.exprs(s.Elts),
			Rbrace: c.next(lenRbrace),
		}
		c.markup(s.After)
	case *Ellipsis:
		c.markup(s.Before)
		e = &ast.Ellipsis{
			Ellipsis: c.next(lenEllipsis),
			Elt:      c.expr(s.Elt),
		}
		c.markup(s.After)
	case *Func:
		c.markup(s.Before)
		var t = &ast.FuncType{
			Func:    c.next(lenFunc),
			Params:  c.node(s.Parameters).(*ast.FieldList),
			Results: c.node(s.Results).(*ast.FieldList),
		}
		if s.Body == nil {
			e = t
		} else {
			e = &ast.FuncLit{
				Type: t,
				Body: c.stmt(s.Body).(*ast.BlockStmt),
			}
		}
		c.markup(s.After)
	case *Name:
		if s == nil {
			return nil
		}
		c.markup(s.Before)
		e = &ast.Ident{
			NamePos: c.next(len(s.Text)),
			Name:    s.Text,
		}
		c.markup(s.After)
	case *Index:
		c.markup(s.Before)
		e = &ast.IndexExpr{
			X:      c.expr(s.X),
			Lbrack: c.next(lenLbrack),
			Index:  c.expr(s.Index),
			Rbrack: c.next(lenRbrack),
		}
		c.markup(s.After)
	case *Interface:
		c.markup(s.Before)
		e = &ast.InterfaceType{
			Interface: c.next(lenInterface),
			Methods:   c.node(s.Methods).(*ast.FieldList),
		}
		c.markup(s.After)
	case *KeyValue:
		c.markup(s.Before)
		e = &ast.KeyValueExpr{
			Key:   c.expr(s.Key),
			Colon: c.next(lenColon),
			Value: c.expr(s.Value),
		}
		c.markup(s.After)
	case *Map:
		c.markup(s.Before)
		e = &ast.MapType{
			Map:   c.next(lenMap),
			Key:   c.expr(s.Key),
			Value: c.expr(s.Value),
		}
		c.markup(s.After)
	case *Paren:
		c.markup(s.Before)
		e = &ast.ParenExpr{
			Lparen: c.next(lenLparen),
			X:      c.expr(s.X),
			Rparen: c.next(lenRparen),
		}
		c.markup(s.After)
	case *Selector:
		c.markup(s.Before)
		e = &ast.SelectorExpr{
			X: c.expr(s.X),
		}
		c.skip(lenPeriod)
		e.(*ast.SelectorExpr).Sel = c.expr(s.Sel).(*ast.Ident)
		c.markup(s.After)
	case *Slice:
		c.markup(s.Before)
		e = &ast.SliceExpr{
			X:      c.expr(s.X),
			Lbrack: c.next(lenLbrack),
			Low:    c.expr(s.Low),
			High:   c.expr(s.High),
			Max:    c.expr(s.Max),
			Slice3: false, // TODO
			Rbrack: c.next(lenRbrack),
		}
		c.markup(s.After)
	case *Struct:
		c.markup(s.Before)
		e = &ast.StructType{
			Struct:     c.next(lenStruct),
			Fields:     c.node(s.Fields).(*ast.FieldList),
			Incomplete: false, // TODO
		}
		c.markup(s.After)
	case *Assert:
		c.markup(s.Before)
		e = &ast.TypeAssertExpr{
			X:      c.expr(s.X),
			Lparen: c.next(lenLparen),
			Type:   c.expr(s.Type),
			Rparen: c.next(lenRparen),
		}
		c.markup(s.After)
	case *Unary:
		c.markup(s.Before)
		if s.Operator == token.MUL {
			e = &ast.StarExpr{
				Star: c.next(lenMul),
				X:    c.expr(s.X),
			}
		} else {
			e = &ast.UnaryExpr{
				OpPos: c.next(len(s.Operator.String())),
				Op:    s.Operator,
				X:     c.expr(s.X),
			}
		}
		c.markup(s.After)
	default:
		panic(fmt.Sprintf("invalid expression: %#v", s)) // TODO: Remove
	}
	return e
}

func (c *syntaxConv) exprs(from []Syntax) []ast.Expr {
	var to []ast.Expr
	for _, f := range from {
		to = append(to, c.expr(f))
	}
	return to
}

func (c *syntaxConv) idents(from []*Name) []*ast.Ident {
	var to []*ast.Ident
	for _, f := range from {
		to = append(to, c.expr(f).(*ast.Ident))
	}
	return to
}

func (c *syntaxConv) markup(ss []Syntax) *ast.CommentGroup {
	var cg *ast.CommentGroup
	var lastLine bool
	for _, s := range ss {
		switch s.(type) {
		case *Comment:
			if cg == nil {
				cg = &ast.CommentGroup{}
			}
			cg.List = append(cg.List, c.node(s).(*ast.Comment))
			lastLine = false
		case *Line:
			if lastLine && cg != nil {
				c.astFile.Comments = append(c.astFile.Comments, cg)
				cg = nil
			}
			c.tokenFile.AddLine(c.tokenFile.Offset(c.next(1)))
			lastLine = true
		default:
			panic(fmt.Sprintf("invalid markup: %#v", s)) // TODO: Remove
		}
	}
	if cg != nil {
		c.astFile.Comments = append(c.astFile.Comments, cg)
	}
	return cg
}

func (c *syntaxConv) next(n int) token.Pos {
	var p = c.end
	c.end += token.Pos(n)
	return p
}

func (c *syntaxConv) node(s Syntax) ast.Node {
	switch s := s.(type) {
	case nil:
		return nil
	case *Comment:
		return &ast.Comment{
			Text: s.Text,
		}
	case *CommentGroup:
		var cs []*ast.Comment
		for _, com := range s.List {
			cs = append(cs, c.node(com).(*ast.Comment))
		}
		return &ast.CommentGroup{
			List: cs,
		}
	case *Field:
		var tag *ast.BasicLit
		if b, ok := c.expr(s.Tag).(*ast.BasicLit); ok {
			tag = b
		}
		return &ast.Field{
			Names: c.idents(s.Names),
			Tag:   tag,
			Type:  c.expr(s.Type),
		}
	case *FieldList:
		if s == nil {
			return (*ast.FieldList)(nil)
		}
		var fs []*ast.Field
		for _, f := range s.List {
			fs = append(fs, c.node(f).(*ast.Field))
		}
		return &ast.FieldList{
			List: fs,
		}
	case *File:
		return &ast.File{
			Name:  c.expr(s.Name).(*ast.Ident),
			Decls: c.decls(s.Decls),
		}
	case *Package:
		var fs map[string]*ast.File
		if s.Files != nil {
			fs = map[string]*ast.File{}
			for k, v := range s.Files {
				fs[k] = c.node(v).(*ast.File)
			}
		}
		return &ast.Package{
			Files: fs,
		}
	default:
		panic(fmt.Sprintf("invalid node: %#v", s)) // TODO: Remove
	}
}

func (c *syntaxConv) skip(n int) {
	c.end += token.Pos(n)
}

func (c *syntaxConv) spec(s Syntax) (spec ast.Spec) {
	switch s := s.(type) {
	case nil:
		return nil
	case *Const:
		c.markup(s.Before)
		spec = &ast.ValueSpec{
			Names:  c.idents(s.Names),
			Type:   c.expr(s.Type),
			Values: c.exprs(s.Values),
		}
		c.markup(s.After)
	case *Import:
		c.markup(s.Before)
		spec = &ast.ImportSpec{
			Name: ident(c.expr(s.Name)),
			Path: c.expr(s.Path).(*ast.BasicLit),
		}
		c.markup(s.After)
	case *Type:
		c.markup(s.Before)
		spec = &ast.TypeSpec{
			Assign: s.Assign,
			Name:   c.expr(s.Name).(*ast.Ident),
			Type:   c.expr(s.Type),
		}
		c.markup(s.After)
	case *Var:
		c.markup(s.Before)
		spec = &ast.ValueSpec{
			Names:  c.idents(s.Names),
			Type:   c.expr(s.Type),
			Values: c.exprs(s.Values),
		}
		c.markup(s.After)
	default:
		panic(fmt.Sprintf("invalid specification: %#v", s)) // TODO: Remove
	}
	return spec
}

func (c *syntaxConv) specs(from []Syntax) []ast.Spec {
	var to []ast.Spec
	for _, f := range from {
		to = append(to, c.spec(f))
	}
	return to
}

func (c *syntaxConv) stmt(s Syntax) ast.Stmt {
	switch s := s.(type) {
	case nil:
		return nil
	case *Assign:
		return &ast.AssignStmt{
			Lhs: c.exprs(s.Left),
			Rhs: c.exprs(s.Right),
			Tok: s.Operator,
		}
	case *Block:
		if s == nil {
			return nil
		}
		return &ast.BlockStmt{
			List: c.stmts(s.List),
		}
	case *Break:
		return &ast.BranchStmt{
			Tok:   token.BREAK,
			Label: ident(c.expr(s.Label)),
		}
	case *Case:
		if s.Comm == nil {
			return &ast.CaseClause{
				Body: c.stmts(s.Body),
				List: c.exprs(s.List),
			}
		}
		return &ast.CommClause{
			Body: c.stmts(s.Body),
			Comm: c.stmt(s.Comm),
		}
	case *Continue:
		return &ast.BranchStmt{
			Tok:   token.CONTINUE,
			Label: ident(c.expr(s.Label)),
		}
	case *Fallthrough:
		return &ast.BranchStmt{
			Tok: token.FALLTHROUGH,
		}
	case *Goto:
		return &ast.BranchStmt{
			Tok:   token.GOTO,
			Label: ident(c.expr(s.Label)),
		}
	case *Defer:
		return &ast.DeferStmt{
			Call: c.expr(s.Call).(*ast.CallExpr),
		}
	case *Empty:
		return &ast.EmptyStmt{}
	case *For:
		return &ast.ForStmt{
			Init: c.stmt(s.Init),
			Cond: c.expr(s.Cond),
			Post: c.stmt(s.Post),
			Body: c.stmt(s.Body).(*ast.BlockStmt),
		}
	case *Go:
		return &ast.GoStmt{
			Call: c.expr(s.Call).(*ast.CallExpr),
		}
	case *If:
		return &ast.IfStmt{
			Init: c.stmt(s.Init),
			Cond: c.expr(s.Cond),
			Body: c.stmt(s.Body).(*ast.BlockStmt),
			Else: c.stmt(s.Else),
		}
	case *Inc:
		return &ast.IncDecStmt{
			X:   c.expr(s.X),
			Tok: token.INC,
		}
	case *Dec:
		return &ast.IncDecStmt{
			X:   c.expr(s.X),
			Tok: token.DEC,
		}
	case *Label:
		return &ast.LabeledStmt{
			Label: c.expr(s.Label).(*ast.Ident),
			Stmt:  c.stmt(s.Stmt),
		}
	case *Range:
		var t = token.DEFINE
		if s.Assign {
			t = token.ASSIGN
		}
		return &ast.RangeStmt{
			Key:   c.expr(s.Key),
			Value: c.expr(s.Value),
			Tok:   t,
			X:     c.expr(s.X),
			Body:  c.stmt(s.Body).(*ast.BlockStmt),
		}
	case *Return:
		return &ast.ReturnStmt{
			Results: c.exprs(s.Results),
		}
	case *Select:
		return &ast.SelectStmt{
			Body: c.stmt(s.Body).(*ast.BlockStmt),
		}
	case *Send:
		return &ast.SendStmt{
			Chan:  c.expr(s.Chan),
			Value: c.expr(s.Value),
		}
	case *Switch:
		if s.Type == nil {
			return &ast.SwitchStmt{
				Body: c.stmt(s.Body).(*ast.BlockStmt),
				Init: c.stmt(s.Init),
				Tag:  c.expr(s.Value),
			}
		}
		return &ast.TypeSwitchStmt{
			Assign: c.stmt(s.Type),
			Body:   c.stmt(s.Body).(*ast.BlockStmt),
			Init:   c.stmt(s.Init),
		}
	default:
		if d := c.decl(s); d != nil {
			return &ast.DeclStmt{
				Decl: d,
			}
		}
		if e := c.expr(s); e != nil {
			return &ast.ExprStmt{
				X: e,
			}
		}
		panic(fmt.Sprintf("invalid statement: %#v", s)) // TODO: Remove
	}
}

func (c *syntaxConv) stmts(from []Syntax) []ast.Stmt {
	var to []ast.Stmt
	for _, f := range from {
		to = append(to, c.stmt(f))
	}
	return to
}
