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

func (c *syntaxConv) decl(from Syntax) (to ast.Decl) {
	switch from := from.(type) {
	case *Const:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.CONST,
			TokPos: c.next(lenConst),
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ConstList:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.CONST,
			TokPos: c.next(lenConst),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(from.List),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *Func:
		c.markup(from.Before)
		to = &ast.FuncDecl{
			Body: blockStmt(c.stmt(from.Body)),
			Name: c.expr(from.Name).(*ast.Ident),
			Recv: c.node(from.Receiver).(*ast.FieldList),
			Type: &ast.FuncType{
				Params:  c.node(from.Parameters).(*ast.FieldList),
				Results: c.node(from.Results).(*ast.FieldList),
			},
		}
		c.markup(from.After)
	case *Import:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: c.next(lenImport),
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ImportList:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: c.next(lenImport),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(from.List),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *Type:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.TYPE,
			TokPos: c.next(lenType),
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *TypeList:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.TYPE,
			TokPos: c.next(lenType),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(from.List),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *Var:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.VAR,
			TokPos: c.next(lenVar),
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *VarList:
		c.markup(from.Before)
		to = &ast.GenDecl{
			Tok:    token.VAR,
			TokPos: c.next(lenVar),
			Lparen: c.next(lenLparen),
			Specs:  c.specs(from.List),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	}
	return to
}

func (c *syntaxConv) decls(from []Syntax) (to []ast.Decl) {
	for _, s := range from {
		to = append(to, c.decl(s))
	}
	return to
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

func (c *syntaxConv) markup(ss []Syntax) {
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
}

func (c *syntaxConv) next(n int) token.Pos {
	var p = c.end
	c.end += token.Pos(n)
	return p
}

func (c *syntaxConv) node(from Syntax) (to ast.Node) {
	switch from := from.(type) {
	case *Comment:
		to = &ast.Comment{
			Slash: c.next(len(from.Text)), // TODO: Insert new line?
			Text:  from.Text,
		}
	case *CommentGroup:
		var cs []*ast.Comment
		for _, com := range from.List {
			cs = append(cs, c.node(com).(*ast.Comment))
		}
		to = &ast.CommentGroup{
			List: cs,
		}
	case *Field:
		var tag *ast.BasicLit
		if b, ok := c.expr(from.Tag).(*ast.BasicLit); ok {
			tag = b
		}
		to = &ast.Field{
			Names: c.idents(from.Names),
			Tag:   tag,
			Type:  c.expr(from.Type),
		}
	case *FieldList:
		if from == nil {
			to = (*ast.FieldList)(nil)
		} else {
			var fs []*ast.Field
			for _, f := range from.List {
				fs = append(fs, c.node(f).(*ast.Field))
			}
			to = &ast.FieldList{
				List: fs,
			}
		}
	case *File:
		to = &ast.File{
			Name:  c.expr(from.Name).(*ast.Ident),
			Decls: c.decls(from.Decls),
		}
	case *Package:
		var fs map[string]*ast.File
		if from.Files != nil {
			fs = map[string]*ast.File{}
			for k, v := range from.Files {
				fs[k] = c.node(v).(*ast.File)
			}
		}
		to = &ast.Package{
			Files: fs,
		}
	}
	return to
}

func (c *syntaxConv) skip(n int) {
	c.end += token.Pos(n)
}

func (c *syntaxConv) spec(from Syntax) (to ast.Spec) {
	switch from := from.(type) {
	case *Const:
		c.markup(from.Before)
		to = &ast.ValueSpec{
			Names:  c.idents(from.Names),
			Type:   c.expr(from.Type),
			Values: c.exprs(from.Values),
		}
		c.markup(from.After)
	case *Import:
		c.markup(from.Before)
		to = &ast.ImportSpec{
			Name:   ident(c.expr(from.Name)),
			Path:   c.expr(from.Path).(*ast.BasicLit),
			EndPos: 0, // TODO
		}
		c.markup(from.After)
	case *Type:
		c.markup(from.Before)
		to = &ast.TypeSpec{
			Assign: from.Assign,
			Name:   c.expr(from.Name).(*ast.Ident),
			Type:   c.expr(from.Type),
		}
		c.markup(from.After)
	case *Var:
		c.markup(from.Before)
		to = &ast.ValueSpec{
			Names:  c.idents(from.Names),
			Type:   c.expr(from.Type),
			Values: c.exprs(from.Values),
		}
		c.markup(from.After)
	}
	return to
}

func (c *syntaxConv) specs(from []Syntax) (to []ast.Spec) {
	for _, f := range from {
		to = append(to, c.spec(f))
	}
	return to
}

func (c *syntaxConv) stmt(from Syntax) (to ast.Stmt) {
	switch from := from.(type) {
	case *Assign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs: c.exprs(from.Left),
			Rhs: c.exprs(from.Right),
			Tok: from.Operator,
		}
		c.markup(from.After)
	case *Block:
		if from != nil {
			c.markup(from.Before)
			to = &ast.BlockStmt{
				List: c.stmts(from.List),
			}
			c.markup(from.After)
		}
	case *Break:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			Tok:   token.BREAK,
			Label: ident(c.expr(from.Label)),
		}
		c.markup(from.After)
	case *Case:
		if from.Comm == nil {
			c.markup(from.Before)
			to = &ast.CaseClause{
				Body: c.stmts(from.Body),
				List: c.exprs(from.List),
			}
			c.markup(from.After)
		} else {
			c.markup(from.Before)
			to = &ast.CommClause{
				Body: c.stmts(from.Body),
				Comm: c.stmt(from.Comm),
			}
			c.markup(from.After)
		}
	case *Continue:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			Tok:   token.CONTINUE,
			Label: ident(c.expr(from.Label)),
		}
		c.markup(from.After)
	case *Dec:
		c.markup(from.Before)
		to = &ast.IncDecStmt{
			X:   c.expr(from.X),
			Tok: token.DEC,
		}
		c.markup(from.After)
	case *Defer:
		c.markup(from.Before)
		to = &ast.DeferStmt{
			Call: c.expr(from.Call).(*ast.CallExpr),
		}
		c.markup(from.After)
	case *Empty:
		c.markup(from.Before)
		to = &ast.EmptyStmt{}
		c.markup(from.After)
	case *Fallthrough:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			Tok: token.FALLTHROUGH,
		}
		c.markup(from.After)
	case *For:
		c.markup(from.Before)
		to = &ast.ForStmt{
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Post: c.stmt(from.Post),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *Go:
		c.markup(from.Before)
		to = &ast.GoStmt{
			Call: c.expr(from.Call).(*ast.CallExpr),
		}
		c.markup(from.After)
	case *Goto:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			Tok:   token.GOTO,
			Label: ident(c.expr(from.Label)),
		}
		c.markup(from.After)
	case *If:
		c.markup(from.Before)
		to = &ast.IfStmt{
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
			Else: c.stmt(from.Else),
		}
		c.markup(from.After)
	case *Inc:
		c.markup(from.Before)
		to = &ast.IncDecStmt{
			X:   c.expr(from.X),
			Tok: token.INC,
		}
		c.markup(from.After)
	case *Label:
		c.markup(from.Before)
		to = &ast.LabeledStmt{
			Label: c.expr(from.Label).(*ast.Ident),
			Stmt:  c.stmt(from.Stmt),
		}
		c.markup(from.After)
	case *Range:
		var t = token.DEFINE
		if from.Assign {
			t = token.ASSIGN
		}
		c.markup(from.Before)
		to = &ast.RangeStmt{
			Key:   c.expr(from.Key),
			Value: c.expr(from.Value),
			Tok:   t,
			X:     c.expr(from.X),
			Body:  c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *Return:
		c.markup(from.Before)
		to = &ast.ReturnStmt{
			Results: c.exprs(from.Results),
		}
		c.markup(from.After)
	case *Select:
		c.markup(from.Before)
		to = &ast.SelectStmt{
			Body: c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *Send:
		c.markup(from.Before)
		to = &ast.SendStmt{
			Chan:  c.expr(from.Chan),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Switch:
		if from.Type == nil {
			c.markup(from.Before)
			to = &ast.SwitchStmt{
				Body: c.stmt(from.Body).(*ast.BlockStmt),
				Init: c.stmt(from.Init),
				Tag:  c.expr(from.Value),
			}
			c.markup(from.After)
		} else {
			c.markup(from.Before)
			to = &ast.TypeSwitchStmt{
				Assign: c.stmt(from.Type),
				Body:   c.stmt(from.Body).(*ast.BlockStmt),
				Init:   c.stmt(from.Init),
			}
			c.markup(from.After)
		}
	default:
		if d := c.decl(from); d != nil {
			to = &ast.DeclStmt{
				Decl: d,
			}
		}
		if e := c.expr(from); e != nil {
			to = &ast.ExprStmt{
				X: e,
			}
		}
	}
	return to
}

func (c *syntaxConv) stmts(from []Syntax) (to []ast.Stmt) {
	for _, f := range from {
		to = append(to, c.stmt(f))
	}
	return to
}
