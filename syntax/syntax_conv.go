package syntax

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

var (
	lenAdd       = len(token.ADD.String())
	lenAnd       = len(token.AND.String())
	lenAndNot    = len(token.AND_NOT.String())
	lenArrow     = len(token.ARROW.String())
	lenAssign    = len(token.ASSIGN.String())
	lenChan      = len(token.CHAN.String())
	lenColon     = len(token.COLON.String())
	lenConst     = len(token.CONST.String())
	lenDefine    = len(token.DEFINE.String())
	lenQuo       = len(token.QUO.String())
	lenEllipsis  = len(token.ELLIPSIS.String())
	lenEql       = len(token.EQL.String())
	lenFunc      = len(token.FUNC.String())
	lenFor       = len(token.FOR.String())
	lenGeq       = len(token.GEQ.String())
	lenGtr       = len(token.GTR.String())
	lenImport    = len(token.IMPORT.String())
	lenInterface = len(token.INTERFACE.String())
	lenLand      = len(token.LAND.String())
	lenLbrace    = len(token.LBRACE.String())
	lenLbrack    = len(token.LBRACK.String())
	lenLeq       = len(token.LEQ.String())
	lenLss       = len(token.LSS.String())
	lenLor       = len(token.LOR.String())
	lenLparen    = len(token.LPAREN.String())
	lenMap       = len(token.MAP.String())
	lenMul       = len(token.MUL.String())
	lenNeq       = len(token.NEQ.String())
	lenNewline   = 1
	lenNot       = len(token.NOT.String())
	lenOr        = len(token.OR.String())
	lenPackage   = len(token.PACKAGE.String())
	lenPeriod    = len(token.PERIOD.String())
	lenRbrace    = len(token.RBRACE.String())
	lenRbrack    = len(token.RBRACK.String())
	lenRem       = len(token.REM.String())
	lenRparen    = len(token.RPAREN.String())
	lenShl       = len(token.SHL.String())
	lenShr       = len(token.SHR.String())
	lenSpace     = 1
	lenStruct    = len(token.STRUCT.String())
	lenSub       = len(token.SUB.String())
	lenSwitch    = len(token.SWITCH.String())
	lenType      = len(token.TYPE.String())
	lenVar       = len(token.VAR.String())
	lenXor       = len(token.XOR.String())
)

func ConvertFile(f *File) (*token.FileSet, *ast.File) {
	var c syntaxConv
	n := c.node(f)
	return c.tokenFileSet, n.(*ast.File)
}

func MustFileString(f *File) string {
	fset, n := ConvertFile(f)
	b := &bytes.Buffer{}
	if err := format.Node(b, fset, n); err != nil {
		panic(err)
	}
	return b.String()
}

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
			TokPos: c.next(lenConst),
			Tok:    token.CONST,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ConstList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.next(lenConst)
		g.Tok = token.CONST
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
	case *Func:
		c.markup(from.Before)
		funcPos := c.next(lenFunc)
		params := from.Parameters
		if params == nil {
			params = &FieldList{}
		}
		to = &ast.FuncDecl{
			Recv: c.node(from.Receiver).(*ast.FieldList),
			Name: c.expr(from.Name).(*ast.Ident),
			Type: &ast.FuncType{
				Func:    funcPos,
				Params:  c.node(params).(*ast.FieldList),
				Results: c.node(from.Results).(*ast.FieldList),
			},
			Body: blockStmt(c.stmt(from.Body)),
		}
		c.markup(from.After)
	case *Import:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.next(lenImport),
			Tok:    token.IMPORT,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ImportList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.next(lenImport)
		g.Tok = token.IMPORT
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
	case *Type:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.next(lenType),
			Tok:    token.TYPE,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *TypeList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.next(lenType)
		g.Tok = token.TYPE
		c.markup(from.Between)
		g.Lparen = c.next(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.next(lenRparen)
		c.markup(from.After)
		to = g
	case *Var:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.next(lenVar),
			Tok:    token.VAR,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *VarList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.next(lenVar)
		g.Tok = token.VAR
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
	case *Add:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenAdd),
			Op:    token.ADD,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *And:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenLand),
			Op:    token.LAND,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *AndNot:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenAndNot),
			Op:    token.AND_NOT,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Array:
		c.markup(from.Before)
		a := &ast.ArrayType{
			Lbrack: c.next(lenLbrack),
			Len:    c.expr(from.Length),
		}
		c.next(lenRbrack)
		a.Elt = c.expr(from.Element)
		c.markup(from.After)
		to = a
	case *Assert:
		c.markup(from.Before)
		to = &ast.TypeAssertExpr{
			X:      c.expr(from.X),
			Lparen: c.next(lenLparen),
			Type:   c.expr(from.Type),
			Rparen: c.next(lenRparen),
		}
		c.markup(from.After)
	case *BitAnd:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenAnd),
			Op:    token.AND,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *BitOr:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenOr),
			Op:    token.OR,
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
	case *Deref:
		c.markup(from.Before)
		to = &ast.StarExpr{
			Star: c.next(lenMul),
			X:    c.expr(from.X),
		}
		c.markup(from.After)
	case *Divide:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenQuo),
			Op:    token.QUO,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Ellipsis: // TODO: Where is this used?
		c.markup(from.Before)
		to = &ast.Ellipsis{
			Ellipsis: c.next(lenEllipsis),
			Elt:      c.expr(from.Elt),
		}
		c.markup(from.After)
	case *Equal:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenEql),
			Op:    token.EQL,
			Y:     c.expr(from.Y),
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
		f := &ast.FuncType{
			Func:    c.next(lenFunc),
			Params:  c.node(from.Parameters).(*ast.FieldList),
			Results: c.node(from.Results).(*ast.FieldList),
		}
		if from.Body == nil {
			to = f
		} else {
			to = &ast.FuncLit{
				Type: f,
				Body: c.stmt(from.Body).(*ast.BlockStmt),
			}
		}
		c.markup(from.After)
	case *Greater:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenGtr),
			Op:    token.GTR,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *GreaterEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenGeq),
			Op:    token.GEQ,
			Y:     c.expr(from.Y),
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
	case *Less:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenLss),
			Op:    token.LSS,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *LessEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenLeq),
			Op:    token.LEQ,
			Y:     c.expr(from.Y),
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
	case *Multiply:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenMul),
			Op:    token.MUL,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Name:
		if from == nil { // TODO: Why is this needed?
			return nil
		}
		c.markup(from.Before)
		to = &ast.Ident{
			NamePos: c.next(len(from.Text)),
			Name:    from.Text,
		}
		c.markup(from.After)
	case *Negate:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.next(lenSub),
			Op:    token.SUB,
			X:     c.expr(from.X),
		}
		c.markup(from.After)
	case *Not:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.next(lenNot),
			Op:    token.NOT,
			X:     c.expr(from.X),
		}
		c.markup(from.After)
	case *NotEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenNeq),
			Op:    token.NEQ,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Or:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenLor),
			Op:    token.LOR,
			Y:     c.expr(from.Y),
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
	case *Pointer:
		c.markup(from.Before)
		to = &ast.StarExpr{
			Star: c.next(lenMul),
			X:    c.expr(from.X),
		}
		c.markup(from.After)
	case *Receive:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.next(lenArrow),
			Op:    token.ARROW,
			X:     c.expr(from.X),
		}
		c.markup(from.After)
	case *Ref:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.next(lenAnd),
			Op:    token.AND,
			X:     c.expr(from.X),
		}
		c.markup(from.After)
	case *Remainder:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenRem),
			Op:    token.REM,
			Y:     c.expr(from.Y),
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
		e := &ast.SelectorExpr{X: c.expr(from.X)}
		c.next(lenPeriod)
		e.Sel = c.expr(from.Sel).(*ast.Ident)
		c.markup(from.After)
		to = e
	case *Send:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenArrow),
			Op:    token.ARROW,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *ShiftLeft:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenShl),
			Op:    token.SHL,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *ShiftRight:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenShr),
			Op:    token.SHR,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Slice:
		c.markup(from.Before)
		to = &ast.SliceExpr{
			X:      c.expr(from.X),
			Lbrack: c.next(lenLbrack),
			Low:    c.expr(from.Low),
			High:   c.expr(from.High),
			Max:    c.expr(from.Max),
			Slice3: from.Max != nil, // TODO: Test
			Rbrack: c.next(lenRbrack),
		}
		c.markup(from.After)
	case *String:
		if from == nil { // TODO: Why needed?
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
			Incomplete: false, // TODO: What is this for?
		}
		c.markup(from.After)
	case *Subtract:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenSub),
			Op:    token.SUB,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	case *Xor:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.X),
			OpPos: c.next(lenXor),
			Op:    token.XOR,
			Y:     c.expr(from.Y),
		}
		c.markup(from.After)
	default:
		panic(from) // TODO: Return error
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
			if s.Count > 0 {
				c.next(s.Count)
			} else {
				c.next(1)
			}
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
			Slash: c.next(len(from.Text)),
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
		c.markup(from.Before)
		f := &ast.Field{}
		f.Names = c.idents(from.Names)
		if b, ok := c.expr(from.Tag).(*ast.BasicLit); ok {
			f.Tag = b
		}
		f.Type = c.expr(from.Type)
		to = f
		c.markup(from.After)
	case *FieldList:
		if from == nil {
			to = (*ast.FieldList)(nil) // TODO: Why?
		} else {
			c.markup(from.Before)
			n := &ast.FieldList{}
			n.Opening = c.next(1) // Either parens or curly braces
			for _, f := range from.List {
				n.List = append(n.List, c.node(f).(*ast.Field))
			}
			n.Closing = c.next(1) // Either parens or curly braces
			c.markup(from.After)
			to = n
		}
	case *File:
		if c.tokenFileSet == nil {
			c.tokenFileSet = token.NewFileSet()
		}
		c.tokenFile = c.tokenFileSet.AddFile("", -1, 99999999999999) // TODO
		c.astFile = &ast.File{}
		c.markup(from.Markup.Before)
		c.astFile.Package = c.next(lenPackage)
		c.astFile.Name = c.expr(from.Package).(*ast.Ident)
		c.astFile.Decls = c.decls(from.Decls)
		c.markup(from.Markup.After)
		to = c.astFile
	default:
		if d := c.decl(from); d != nil {
			to = d
		} else if e := c.expr(from); e != nil {
			to = e
		} else if s := c.spec(from); s != nil {
			to = s
		} else if s := c.stmt(from); s != nil {
			to = s
		} else {
			panic(from) // TODO: Return error
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
			EndPos: 0, // TODO: Verify this should be 0
		}
		c.markup(from.After)
	case *Type:
		c.markup(from.Before)
		s := &ast.TypeSpec{}
		s.Name = c.expr(from.Name).(*ast.Ident)
		if from.Assign {
			s.Assign = c.next(lenAssign)
		}
		s.Type = c.expr(from.Type)
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
	case *AddAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.ADD_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *AndNotAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.AND_NOT_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Assign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *BitAndAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.AND_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *BitOrAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.OR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Block:
		if from != nil {
			c.markup(from.Before)
			to = &ast.BlockStmt{
				Lbrace: c.next(lenLbrace),
				List:   c.stmts(from.List),
				Rbrace: c.next(lenRbrace),
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
		c.markup(from.Before)
		if from.Comm == nil {
			to = &ast.CaseClause{
				Body: c.stmts(from.Body),
				List: c.exprs(from.List),
			}
		} else {
			to = &ast.CommClause{
				Body: c.stmts(from.Body),
				Comm: c.stmt(from.Comm),
			}
		}
		c.markup(from.After)
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
	case *Define:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenDefine),
			Tok:    token.DEFINE,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *DivideAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.QUO_ASSIGN,
			Rhs:    c.exprs(from.Right),
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
	case *MultiplyAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.MUL_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Range:
		var t token.Token
		var l int
		if from.Assign {
			t = token.ASSIGN
			l = lenAssign
		} else {
			t = token.DEFINE
			l = lenDefine
		}
		c.markup(from.Before)
		to = &ast.RangeStmt{
			For:    c.next(lenFor),
			Key:    c.expr(from.Key),
			Value:  c.expr(from.Value),
			TokPos: c.next(l), // TODO: Should not set if Key==nil
			Tok:    t,
			X:      c.expr(from.X),
			Body:   c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *RemainderAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.REM_ASSIGN,
			Rhs:    c.exprs(from.Right),
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
			Chan:  c.expr(from.X),
			Arrow: c.next(lenArrow),
			Value: c.expr(from.Y),
		}
		c.markup(from.After)
	case *ShiftLeftAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.SHL_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *ShiftRightAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.SHR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *SubtractAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.SUB_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Switch:
		c.markup(from.Before)
		if from.Type == nil {
			to = &ast.SwitchStmt{
				Switch: c.next(lenSwitch),
				Init:   c.stmt(from.Init),
				Tag:    c.expr(from.Value),
				Body:   c.stmt(from.Body).(*ast.BlockStmt),
			}
		} else {
			to = &ast.TypeSwitchStmt{
				Switch: c.next(lenSwitch),
				Init:   c.stmt(from.Init),
				Assign: c.stmt(from.Type),
				Body:   c.stmt(from.Body).(*ast.BlockStmt),
			}
		}
		c.markup(from.After)
	case *XorAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.next(lenAssign),
			Tok:    token.XOR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	default:
		if d := c.decl(from); d != nil {
			to = &ast.DeclStmt{Decl: d}
		} else if e := c.expr(from); e != nil {
			to = &ast.ExprStmt{X: e}
		} else {
			panic(from) // TODO: Return error
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
