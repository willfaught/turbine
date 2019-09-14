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
	lenNewline   = 1
	lenPackage   = len(token.PACKAGE.String())
	lenPeriod    = len(token.PERIOD.String())
	lenRbrace    = len(token.RBRACE.String())
	lenRbrack    = len(token.RBRACK.String())
	lenRparen    = len(token.RPAREN.String())
	lenSpace     = 1
	lenStruct    = len(token.STRUCT.String())
	lenType      = len(token.TYPE.String())
	lenVar       = len(token.VAR.String())
)

func convertSyntax(s Syntax) ast.Node {
	var c syntaxConv
	return c.node(s)
}

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
	last         token.Pos
	newLineEmpty bool // new line added but nothing added to it yet. TODO: check if true at end and remove last line from token file if true.
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
		g := &ast.GenDecl{}
		g.Tok = token.CONST
		g.TokPos = c.next(lenConst)
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
	case *Func:
		c.markup(from.Before)
		funcPos := c.next(lenFunc)
		to = &ast.FuncDecl{
			Recv: c.node(from.Receiver).(*ast.FieldList),
			Name: c.expr(from.Name).(*ast.Ident),
			Type: &ast.FuncType{
				Func:    funcPos,
				Params:  c.node(from.Parameters).(*ast.FieldList),
				Results: c.node(from.Results).(*ast.FieldList),
			},
			Body: blockStmt(c.stmt(from.Body)),
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
		g := &ast.GenDecl{}
		g.Tok = token.IMPORT
		g.TokPos = c.next(lenImport)
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
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
		g := &ast.GenDecl{}
		g.Tok = token.TYPE
		g.TokPos = c.next(lenType)
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
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
		g := &ast.GenDecl{}
		g.Tok = token.VAR
		g.TokPos = c.next(lenVar)
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
	}
	return to
}

func (c *syntaxConv) decls(from []Syntax) (to []ast.Decl) {
	for _, s := range from {
		to = append(to, c.decl(s))
	}
	return to
}

func (c *syntaxConv) expr(from Syntax) (to ast.Expr) {
	switch from := from.(type) {
	case *Array:
		c.markup(from.Before)
		to = &ast.ArrayType{
			Lbrack: c.next(lenLbrack),
			Len:    c.expr(from.Length),
		}
		c.next(lenRbrack)
		to.(*ast.ArrayType).Elt = c.expr(from.Element)
		c.markup(from.After)
	case *Assert:
		c.markup(from.Before)
		to = &ast.TypeAssertExpr{
			X:      c.expr(from.X),
			Lparen: c.next(lenLparen),
			Type:   c.expr(from.Type),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *Binary:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(len(from.Operator.String())),
			Op:    from.Operator,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Call:
		c.markup(from.Before)
		to = &ast.CallExpr{
			Fun:      c.expr(from.Fun),
			Lparen:   c.next(lenLparen),
			Args:     c.exprs(from.Args),
			Ellipsis: 0, // TODO
			Rparen:   c.next(lenRparen),
		}
		c.markup(from.After)
	case *Chan:
		c.markup(from.Before)
		to = &ast.ChanType{
			Begin: c.next(lenChan + 1),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *ChanIn:
		c.markup(from.Before)
		to = &ast.ChanType{
			Begin: c.next(lenChan),
			Arrow: c.next(lenArrow + 1),
			Dir:   ast.RECV,
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *ChanOut:
		c.markup(from.Before)
		var p = c.next(lenChan + lenArrow + 1)
		to = &ast.ChanType{
			Begin: p,
			Arrow: p,
			Dir:   ast.SEND,
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Composite:
		c.markup(from.Before)
		to = &ast.CompositeLit{
			Type:   c.expr(from.Type),
			Lbrace: c.next(lenLbrace),
			Elts:   c.exprs(from.Elts),
			Rbrace: c.next(lenRbrace),
		}
		c.markup(from.After)
	case *Ellipsis:
		c.markup(from.Before)
		to = &ast.Ellipsis{
			Ellipsis: c.next(lenEllipsis),
			Elt:      c.expr(from.Elt),
		}
		c.markup(from.After)
	case *Float:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.next(len(from.Text)),
			Kind:     token.FLOAT,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Func:
		c.markup(from.Before)
		var t = &ast.FuncType{
			Func:    c.next(lenFunc),
			Params:  c.node(from.Parameters).(*ast.FieldList),
			Results: c.node(from.Results).(*ast.FieldList),
		}
		if from.Body == nil {
			to = t
		} else {
			to = &ast.FuncLit{
				Type: t,
				Body: c.stmt(from.Body).(*ast.BlockStmt),
			}
		}
		c.markup(from.After)
	case *Imag:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.next(len(from.Text)),
			Kind:     token.IMAG,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Index:
		c.markup(from.Before)
		to = &ast.IndexExpr{
			X:      c.expr(from.X),
			Lbrack: c.next(lenLbrack),
			Index:  c.expr(from.Index),
			Rbrack: c.next(lenRbrack),
		}
		c.markup(from.After)
	case *Int:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.next(len(from.Text)),
			Kind:     token.INT,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Interface:
		c.markup(from.Before)
		to = &ast.InterfaceType{
			Interface: c.next(lenInterface),
			Methods:   c.node(from.Methods).(*ast.FieldList),
		}
		c.markup(from.After)
	case *KeyValue:
		c.markup(from.Before)
		to = &ast.KeyValueExpr{
			Key:   c.expr(from.Key),
			Colon: c.next(lenColon),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Map:
		c.markup(from.Before)
		to = &ast.MapType{
			Map:   c.next(lenMap),
			Key:   c.expr(from.Key),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Name:
		if from == nil {
			return nil
		}
		c.markup(from.Before)
		to = &ast.Ident{
			NamePos: c.next(len(from.Text)),
			Name:    from.Text,
		}
		c.markup(from.After)
	case *Paren:
		c.markup(from.Before)
		to = &ast.ParenExpr{
			Lparen: c.next(lenLparen),
			X:      c.expr(from.X),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *Rune:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.next(len(from.Text)),
			Kind:     token.CHAR,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Selector:
		c.markup(from.Before)
		to = &ast.SelectorExpr{
			X: c.expr(from.X),
		}
		c.next(lenPeriod)
		to.(*ast.SelectorExpr).Sel = c.expr(from.Sel).(*ast.Ident)
		c.markup(from.After)
	case *Slice:
		c.markup(from.Before)
		to = &ast.SliceExpr{
			X:      c.expr(from.X),
			Lbrack: c.next(lenLbrack),
			Low:    c.expr(from.Low),
			High:   c.expr(from.High),
			Max:    c.expr(from.Max),
			Slice3: false, // TODO
			Rbrack: c.next(lenRbrack),
		}
		c.markup(from.After)
	case *String:
		if from == nil {
			return nil
		}
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.next(len(from.Text)),
			Kind:     token.STRING,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Struct:
		c.markup(from.Before)
		to = &ast.StructType{
			Struct:     c.next(lenStruct),
			Fields:     c.node(from.Fields).(*ast.FieldList),
			Incomplete: false, // TODO
		}
		c.markup(from.After)
	case *Unary:
		c.markup(from.Before)
		if from.Operator == token.MUL {
			to = &ast.StarExpr{
				Star: c.next(lenMul),
				X:    c.expr(from.X),
			}
		} else {
			to = &ast.UnaryExpr{
				OpPos: c.next(len(from.Operator.String())),
				Op:    from.Operator,
				X:     c.expr(from.X),
			}
		}
		c.markup(from.After)
	}
	return to
}

func (c *syntaxConv) exprs(from []Syntax) (to []ast.Expr) {
	for _, f := range from {
		to = append(to, c.expr(f))
	}
	return to
}

func (c *syntaxConv) idents(from []*Name) (to []*ast.Ident) {
	for _, f := range from {
		to = append(to, c.expr(f).(*ast.Ident))
	}
	return to
}

func (c *syntaxConv) markup(ss []Syntax) {
	if c.astFile == nil || c.tokenFile == nil { // TODO: Remove these, require proper setup
		return
	}
	var cg *ast.CommentGroup
	var lastLine bool // was last syntax item a *Line?
	for _, s := range ss {
		switch s := s.(type) {
		case *Comment:
			if cg == nil {
				cg = &ast.CommentGroup{}
			}
			cg.List = append(cg.List, c.node(s).(*ast.Comment)) // TODO: Add newlines in comment
			lastLine = false
		case *Line:
			if lastLine && cg != nil {
				c.astFile.Comments = append(c.astFile.Comments, cg)
				cg = nil
			}
			c.tokenFile.AddLine(c.tokenFile.Offset(c.next(1) + 1))
			c.newLineEmpty = true
			lastLine = true
		case *Space:
			c.next(1)
		case *Spaces:
			c.next(s.Count)
		default:
			panic(fmt.Sprintf("invalid markup: %#v", s)) // TODO: Remove
		}
	}
	if cg != nil {
		c.astFile.Comments = append(c.astFile.Comments, cg)
	}
}

func (c *syntaxConv) next(n int) token.Pos {
	var p = c.last + 1
	c.last += token.Pos(n)
	c.newLineEmpty = false
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
				Opening: c.next(1),
				List:    fs,
				Closing: c.next(1),
			}
		}
	case *File:
		if c.tokenFileSet == nil {
			c.tokenFileSet = token.NewFileSet()
		}
		c.tokenFile = c.tokenFileSet.AddFile("", -1, 999999999) // TODO
		c.astFile = &ast.File{}
		c.markup(from.Markup.Before)
		c.astFile.Package = c.next(lenPackage)
		c.astFile.Name = c.expr(from.Package).(*ast.Ident)
		c.astFile.Decls = c.decls(from.Decls)
		c.markup(from.Markup.After)
		to = c.astFile
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
	default:
		if d := c.decl(from); d != nil {
			to = d
		} else if e := c.expr(from); e != nil {
			to = e
		} else if s := c.spec(from); s != nil {
			to = s
		} else if s := c.stmt(from); s != nil {
			to = s
		}
	}
	return to
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
