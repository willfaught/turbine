package syntax

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

var (
	lenAdd          = len(token.ADD.String())
	lenAddAssign    = len(token.ADD_ASSIGN.String())
	lenAnd          = len(token.AND.String())
	lenAndAssign    = len(token.AND_ASSIGN.String())
	lenAndNot       = len(token.AND_NOT.String())
	lenAndNotAssign = len(token.AND_NOT_ASSIGN.String())
	lenArrow        = len(token.ARROW.String())
	lenAssign       = len(token.ASSIGN.String())
	lenBreak        = len(token.BREAK.String())
	lenCase         = len(token.CASE.String())
	lenChan         = len(token.CHAN.String())
	lenColon        = len(token.COLON.String())
	lenContinue     = len(token.CONTINUE.String())
	lenConst        = len(token.CONST.String())
	lenDec          = len(token.DEC.String())
	lenDefer        = len(token.DEFER.String())
	lenDefine       = len(token.DEFINE.String())
	lenEllipsis     = len(token.ELLIPSIS.String())
	lenEql          = len(token.EQL.String())
	lenFallthrough  = len(token.FALLTHROUGH.String())
	lenFor          = len(token.FOR.String())
	lenFunc         = len(token.FUNC.String())
	lenGeq          = len(token.GEQ.String())
	lenGtr          = len(token.GTR.String())
	lenGo           = len(token.GO.String())
	lenGoto         = len(token.GOTO.String())
	lenIf           = len(token.IF.String())
	lenImport       = len(token.IMPORT.String())
	lenInc          = len(token.INC.String())
	lenInterface    = len(token.INTERFACE.String())
	lenLand         = len(token.LAND.String())
	lenLbrace       = len(token.LBRACE.String())
	lenLbrack       = len(token.LBRACK.String())
	lenLeq          = len(token.LEQ.String())
	lenLor          = len(token.LOR.String())
	lenLparen       = len(token.LPAREN.String())
	lenLss          = len(token.LSS.String())
	lenMap          = len(token.MAP.String())
	lenMul          = len(token.MUL.String())
	lenMulAssign    = len(token.MUL_ASSIGN.String())
	lenNeq          = len(token.NEQ.String())
	lenNewline      = 1
	lenNot          = len(token.NOT.String())
	lenOr           = len(token.OR.String())
	lenOrAssign     = len(token.OR_ASSIGN.String())
	lenPackage      = len(token.PACKAGE.String())
	lenPeriod       = len(token.PERIOD.String())
	lenQuo          = len(token.QUO.String())
	lenQuoAssign    = len(token.QUO_ASSIGN.String())
	lenRbrace       = len(token.RBRACE.String())
	lenRbrack       = len(token.RBRACK.String())
	lenRem          = len(token.REM.String())
	lenRemAssign    = len(token.REM_ASSIGN.String())
	lenReturn       = len(token.RETURN.String())
	lenRparen       = len(token.RPAREN.String())
	lenSelect       = len(token.SELECT.String())
	lenShl          = len(token.SHL.String())
	lenShlAssign    = len(token.SHL_ASSIGN.String())
	lenShr          = len(token.SHR.String())
	lenShrAssign    = len(token.SHR_ASSIGN.String())
	lenSpace        = 1
	lenStruct       = len(token.STRUCT.String())
	lenSub          = len(token.SUB.String())
	lenSubAssign    = len(token.SUB_ASSIGN.String())
	lenSwitch       = len(token.SWITCH.String())
	lenType         = len(token.TYPE.String())
	lenVar          = len(token.VAR.String())
	lenXor          = len(token.XOR.String())
	lenXorAssign    = len(token.XOR_ASSIGN.String())
)

func ToNode(s Syntax) (*token.FileSet, ast.Node) {
	c := newSyntaxConv()
	return c.tokenFileSet, c.node(s)
}

func ToString(s Syntax) (string, error) {
	fset, n := ToNode(s)
	b := &bytes.Buffer{}
	if err := format.Node(b, fset, n); err != nil {
		return "", fmt.Errorf("cannot format node: %v", err)
	}
	return b.String(), nil
}

type syntaxConv struct {
	astFile      *ast.File
	eol          bool
	last         token.Pos
	tokenFile    *token.File
	tokenFileSet *token.FileSet
}

func newSyntaxConv() *syntaxConv {
	fset := token.NewFileSet()
	return &syntaxConv{
		astFile:      &ast.File{},
		tokenFile:    fset.AddFile("", -1, int((^uint(0))>>1)),
		tokenFileSet: fset,
	}
}

func (c *syntaxConv) add(n int) token.Pos {
	p := c.last + 1
	c.last += token.Pos(n)
	if c.eol {
		c.tokenFile.AddLine(c.tokenFile.Offset(p))
		c.eol = false
	}
	return p
}

func (c *syntaxConv) decl(from Declaration) (to ast.Decl) {
	switch from := from.(type) {
	case *Const:
		c.gaps(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenConst),
			Tok:    token.CONST,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.gaps(from.After)
	case *ConstList:
		c.gaps(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenConst)
		g.Tok = token.CONST
		c.gaps(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.gaps(from.After)
		to = g
	case *Func:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.gaps(from.Before)
		funcPos := c.add(lenFunc)
		funcDecl := &ast.FuncDecl{}
		if from.Receiver != nil {
			funcDecl.Recv = c.node(from.Receiver).(*ast.FieldList)
		}
		funcDecl.Name = c.expr(from.Name).(*ast.Ident)
		funcDecl.Type = &ast.FuncType{
			Func:   funcPos,
			Params: c.node(from.Params).(*ast.FieldList),
		}
		if from.Results != nil {
			funcDecl.Type.Results = c.results(from.Results)
		}
		if from.Body != nil {
			funcDecl.Body = c.stmt(from.Body).(*ast.BlockStmt)
		}
		c.gaps(from.After)
		to = funcDecl
	case *Import:
		c.gaps(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenImport),
			Tok:    token.IMPORT,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.gaps(from.After)
	case *ImportList:
		c.gaps(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenImport)
		g.Tok = token.IMPORT
		c.gaps(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.gaps(from.After)
		to = g
	case *Type:
		c.gaps(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenType),
			Tok:    token.TYPE,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.gaps(from.After)
	case *TypeList:
		c.gaps(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenType)
		g.Tok = token.TYPE
		c.gaps(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.gaps(from.After)
		to = g
	case *Var:
		c.gaps(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenVar),
			Tok:    token.VAR,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.gaps(from.After)
	case *VarList:
		c.gaps(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenVar)
		g.Tok = token.VAR
		c.gaps(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.gaps(from.After)
		to = g
	}
	return to
}

func (c *syntaxConv) decls(from []Declaration) (to []ast.Decl) {
	to = make([]ast.Decl, len(from))
	for i, d := range from {
		to[i] = c.decl(d)
	}
	return to
}

func (c *syntaxConv) expr(from Expression) (to ast.Expr) {
	switch from := from.(type) {
	case nil:
	case *Add:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAdd),
			Op:    token.ADD,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *And:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLand),
			Op:    token.LAND,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *AndNot:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAndNot),
			Op:    token.AND_NOT,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Array:
		c.gaps(from.Before)
		a := &ast.ArrayType{
			Lbrack: c.add(lenLbrack),
			Len:    c.expr(from.Length),
		}
		c.add(lenRbrack)
		a.Elt = c.expr(from.Element)
		c.gaps(from.After)
		to = a
	case *Assert:
		c.gaps(from.Before)
		to = &ast.TypeAssertExpr{
			X:      c.expr(from.Value),
			Lparen: c.add(lenLparen),
			Type:   c.expr(from.Type),
			Rparen: c.add(lenRparen),
		}
		c.gaps(from.After)
	case *BitAnd:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAnd),
			Op:    token.AND,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *BitOr:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenOr),
			Op:    token.OR,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Call:
		c.gaps(from.Before)
		call := &ast.CallExpr{
			Fun:    c.expr(from.Fun),
			Lparen: c.add(lenLparen),
		}
		if l := len(from.Args); l > 0 {
			args := make([]ast.Expr, l)
			last := l - 1
			for i, e := range from.Args[:last] {
				args[i] = c.expr(e)
			}
			if e, ok := from.Args[last].(*Ellipsis); ok {
				c.gaps(e.Before)
				call.Ellipsis = c.add(lenEllipsis)
				args[last] = c.expr(e.Elem)
				c.gaps(e.After)
			} else {
				args[last] = c.expr(from.Args[last])
			}
			call.Args = args
		}
		call.Rparen = c.add(lenRparen)
		c.gaps(from.After)
		to = call
	case *Chan:
		c.gaps(from.Before)
		to = &ast.ChanType{
			Dir:   ast.RECV | ast.SEND,
			Begin: c.add(lenChan + 1),
			Value: c.expr(from.Value),
		}
		c.gaps(from.After)
	case *ChanIn:
		c.gaps(from.Before)
		to = &ast.ChanType{
			Dir:   ast.RECV,
			Begin: c.add(lenChan),
			Arrow: c.add(lenArrow + 1),
			Value: c.expr(from.Value),
		}
		c.gaps(from.After)
	case *ChanOut:
		c.gaps(from.Before)
		var p = c.add(lenChan + lenArrow + 1)
		to = &ast.ChanType{
			Dir:   ast.SEND,
			Begin: p,
			Arrow: p,
			Value: c.expr(from.Value),
		}
		c.gaps(from.After)
	case *Composite:
		c.gaps(from.Before)
		to = &ast.CompositeLit{
			Type:   c.expr(from.Type),
			Lbrace: c.add(lenLbrace),
			Elts:   c.exprs(from.Elts),
			Rbrace: c.add(lenRbrace),
		}
		c.gaps(from.After)
	case *Deref:
		c.gaps(from.Before)
		to = &ast.StarExpr{
			Star: c.add(lenMul),
			X:    c.expr(from.Left),
		}
		c.gaps(from.After)
	case *Divide:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenQuo),
			Op:    token.QUO,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Ellipsis:
		c.gaps(from.Before)
		to = &ast.Ellipsis{
			Ellipsis: c.add(lenEllipsis),
			Elt:      c.expr(from.Elem),
		}
		c.gaps(from.After)
	case *Equal:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenEql),
			Op:    token.EQL,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Float:
		c.gaps(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.FLOAT,
			Value:    from.Text,
		}
		c.gaps(from.After)
	case *Func:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.gaps(from.Before)
		funcType := &ast.FuncType{
			Func:    c.add(lenFunc),
			Params:  c.node(from.Params).(*ast.FieldList),
			Results: c.node(from.Results).(*ast.FieldList),
		}
		if from.Body == nil {
			to = funcType
		} else {
			to = &ast.FuncLit{
				Type: funcType,
				Body: c.stmt(from.Body).(*ast.BlockStmt),
			}
		}
		c.gaps(from.After)
	case *Greater:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenGtr),
			Op:    token.GTR,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *GreaterEqual:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenGeq),
			Op:    token.GEQ,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Imag:
		c.gaps(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.IMAG,
			Value:    from.Text,
		}
		c.gaps(from.After)
	case *Index:
		c.gaps(from.Before)
		to = &ast.IndexExpr{
			X:      c.expr(from.Container),
			Lbrack: c.add(lenLbrack),
			Index:  c.expr(from.Param),
			Rbrack: c.add(lenRbrack),
		}
		c.gaps(from.After)
	case *Int:
		c.gaps(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.INT,
			Value:    from.Text,
		}
		c.gaps(from.After)
	case *Interface:
		if from.Methods == nil {
			from.Methods = &MethodList{}
		}
		c.gaps(from.Before)
		to = &ast.InterfaceType{
			Interface: c.add(lenInterface),
			Methods:   c.node(from.Methods).(*ast.FieldList),
		}
		c.gaps(from.After)
	case *KeyValue:
		c.gaps(from.Before)
		to = &ast.KeyValueExpr{
			Key:   c.expr(from.Key),
			Colon: c.add(lenColon),
			Value: c.expr(from.Value),
		}
		c.gaps(from.After)
	case *Less:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLss),
			Op:    token.LSS,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *LessEqual:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLeq),
			Op:    token.LEQ,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Map:
		c.gaps(from.Before)
		to = &ast.MapType{
			Map:   c.add(lenMap),
			Key:   c.expr(from.Key),
			Value: c.expr(from.Value),
		}
		c.gaps(from.After)
	case *Multiply:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenMul),
			Op:    token.MUL,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Name:
		if from == nil { // Branch labels, import aliases
			to = (*ast.Ident)(nil)
		} else {
			c.gaps(from.Before)
			to = &ast.Ident{
				NamePos: c.add(len(from.Text)),
				Name:    from.Text,
			}
			c.gaps(from.After)
		}
	case *Negate:
		c.gaps(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenSub),
			Op:    token.SUB,
			X:     c.expr(from.Left),
		}
		c.gaps(from.After)
	case *Not:
		c.gaps(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenNot),
			Op:    token.NOT,
			X:     c.expr(from.Left),
		}
		c.gaps(from.After)
	case *NotEqual:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenNeq),
			Op:    token.NEQ,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Or:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLor),
			Op:    token.LOR,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Paren:
		c.gaps(from.Before)
		to = &ast.ParenExpr{
			Lparen: c.add(lenLparen),
			X:      c.expr(from.Left),
			Rparen: c.add(lenRparen),
		}
		c.gaps(from.After)
	case *Pointer:
		c.gaps(from.Before)
		to = &ast.StarExpr{
			Star: c.add(lenMul),
			X:    c.expr(from.Left),
		}
		c.gaps(from.After)
	case *Receive:
		c.gaps(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenArrow),
			Op:    token.ARROW,
			X:     c.expr(from.Left),
		}
		c.gaps(from.After)
	case *Ref:
		c.gaps(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenAnd),
			Op:    token.AND,
			X:     c.expr(from.Left),
		}
		c.gaps(from.After)
	case *Remainder:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenRem),
			Op:    token.REM,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Rune:
		c.gaps(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.CHAR,
			Value:    from.Text,
		}
		c.gaps(from.After)
	case *Selector:
		c.gaps(from.Before)
		e := &ast.SelectorExpr{X: c.expr(from.Value)}
		c.add(lenPeriod)
		e.Sel = c.expr(from.Name).(*ast.Ident)
		c.gaps(from.After)
		to = e
	case *ShiftLeft:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenShl),
			Op:    token.SHL,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *ShiftRight:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenShr),
			Op:    token.SHR,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Slice:
		c.gaps(from.Before)
		to = &ast.SliceExpr{
			X:      c.expr(from.Slice),
			Lbrack: c.add(lenLbrack),
			Low:    c.expr(from.Low),
			High:   c.expr(from.High),
			Max:    c.expr(from.Max),
			Slice3: from.Max != nil,
			Rbrack: c.add(lenRbrack),
		}
		c.gaps(from.After)
	case *String:
		c.gaps(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.STRING,
			Value:    from.Text,
		}
		c.gaps(from.After)
	case *Struct:
		if from.Fields == nil {
			from.Fields = &FieldList{}
		}
		c.gaps(from.Before)
		to = &ast.StructType{
			Struct:     c.add(lenStruct),
			Fields:     c.node(from.Fields).(*ast.FieldList),
			Incomplete: false, // TODO: What is this for?
		}
		c.gaps(from.After)
	case *Subtract:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenSub),
			Op:    token.SUB,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	case *Xor:
		c.gaps(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenXor),
			Op:    token.XOR,
			Y:     c.expr(from.Right),
		}
		c.gaps(from.After)
	default:
		panic(fmt.Sprintf("invalid expression: %#v", from))
	}
	return to
}

func (c *syntaxConv) exprs(from []Expression) (to []ast.Expr) {
	to = make([]ast.Expr, len(from))
	for i, e := range from {
		to[i] = c.expr(e)
	}
	return to
}

func (c *syntaxConv) gaps(gs []Gap) {
	var cg *ast.CommentGroup
	var previousWasLine bool
	for _, g := range gs {
		switch g := g.(type) {
		case *Comment:
			if cg == nil {
				cg = &ast.CommentGroup{}
			}
			p := c.add(len(g.Text))
			cg.List = append(cg.List, &ast.Comment{
				Slash: p,
				Text:  g.Text,
			})
			text := g.Text
			for i := strings.IndexByte(text, '\n'); i != -1; text = text[i+1:] {
				c.tokenFile.AddLine(c.tokenFile.Offset(p + token.Pos(i)))
			}
			previousWasLine = false
		case *Line:
			if previousWasLine && cg != nil {
				c.astFile.Comments = append(c.astFile.Comments, cg)
				cg = nil
			}
			c.add(lenNewline)
			c.eol = true
			previousWasLine = true
		case *Space:
			n := lenSpace
			if g.Count > 0 {
				n *= g.Count
			}
			c.add(n)
		default:
			panic(fmt.Sprintf("invalid gap: %#v", g))
		}
	}
	if cg != nil {
		c.astFile.Comments = append(c.astFile.Comments, cg)
	}
}

func (c *syntaxConv) idents(from []*Name) (to []*ast.Ident) {
	to = make([]*ast.Ident, len(from))
	for i, n := range from {
		to[i] = c.expr(n).(*ast.Ident)
	}
	return to
}

func (c *syntaxConv) results(from *ParamList) (to *ast.FieldList) {
	parens := len(from.List) != 1 || len(from.List[0].Names) > 0
	c.gaps(from.Before)
	to = &ast.FieldList{}
	if parens {
		to.Opening = c.add(lenLparen)
	}
	for _, p := range from.List {
		to.List = append(to.List, c.node(p).(*ast.Field))
	}
	if parens {
		to.Closing = c.add(lenRparen)
	}
	c.gaps(from.After)
	return to
}

func (c *syntaxConv) node(from Syntax) (to ast.Node) {
	switch from := from.(type) {
	case nil:
	case *Field:
		c.gaps(from.Before)
		field := &ast.Field{}
		field.Names = c.idents(from.Names)
		if from.Tag != nil {
			field.Tag = c.expr(from.Tag).(*ast.BasicLit)
		}
		field.Type = c.expr(from.Type)
		c.gaps(from.After)
		to = field
	case *FieldList:
		c.gaps(from.Before)
		fieldList := &ast.FieldList{}
		fieldList.Opening = c.add(lenLbrace)
		for _, f := range from.List {
			fieldList.List = append(fieldList.List, c.node(f).(*ast.Field))
		}
		fieldList.Closing = c.add(lenRbrace)
		c.gaps(from.After)
		to = fieldList
	case *File:
		c.gaps(from.Before)
		c.astFile.Package = c.add(lenPackage)
		c.astFile.Name = c.expr(from.Package).(*ast.Ident)
		c.astFile.Decls = c.decls(from.Decls)
		c.gaps(from.After)
		to = c.astFile
	case *Method:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.gaps(from.Before)
		to = &ast.Field{
			Names: []*ast.Ident{c.expr(from.Name).(*ast.Ident)},
			Type: &ast.FuncType{
				Params: c.node(from.Params).(*ast.FieldList),
			},
		}
		if from.Results != nil {
			to.Type.Results = c.results(from.Results)
		}
		c.gaps(from.After)
	case *MethodList:
		c.gaps(from.Before)
		fieldList := &ast.FieldList{}
		fieldList.Opening = c.add(lenLbrace)
		for _, m := range from.List {
			fieldList.List = append(fieldList.List, c.node(m).(*ast.Field))
		}
		fieldList.Closing = c.add(lenRbrace)
		c.gaps(from.After)
		to = fieldList
	case *Param:
		c.gaps(from.Before)
		to = &ast.Field{
			Names: c.idents(from.Names),
			Type:  c.expr(from.Type),
		}
		c.gaps(from.After)
	case *ParamList:
		c.gaps(from.Before)
		fieldList := &ast.FieldList{}
		fieldList.Opening = c.add(lenLparen)
		for _, p := range from.List {
			fieldList.List = append(fieldList.List, c.node(p).(*ast.Field))
		}
		fieldList.Closing = c.add(lenRparen)
		c.gaps(from.After)
		to = fieldList
	case *Receiver:
		var names []*Name
		if from.Name != nil {
			names = []*Name{from.Name}
		}
		to = c.node(&ParamList{
			Context: Context{
				Before: from.Before,
				After:  from.After,
			},
			List: []*Param{
				&Param{
					Names: names,
					Type:  from.Type,
				},
			},
		})
	default:
		if d, ok := from.(Declaration); ok {
			to = c.decl(d)
		} else if e, ok := from.(Expression); ok {
			to = c.expr(e)
		} else if s, ok := from.(Statement); ok {
			to = c.stmt(s)
		} else {
			panic(fmt.Sprintf("invalid node: %#v", from))
		}
	}
	return to
}

func (c *syntaxConv) spec(from Syntax) (to ast.Spec) {
	switch from := from.(type) {
	case *Const:
		c.gaps(from.Before)
		to = &ast.ValueSpec{
			Names:  c.idents(from.Names),
			Type:   c.expr(from.Type),
			Values: c.exprs(from.Values),
		}
		c.gaps(from.After)
	case *Import:
		c.gaps(from.Before)
		to = &ast.ImportSpec{
			Name:   c.expr(from.Name).(*ast.Ident),
			Path:   c.expr(from.Path).(*ast.BasicLit),
			EndPos: 0, // TODO: Verify this should be 0
		}
		c.gaps(from.After)
	case *Type:
		c.gaps(from.Before)
		s := &ast.TypeSpec{}
		s.Name = c.expr(from.Name).(*ast.Ident)
		if from.Assign {
			s.Assign = c.add(lenAssign)
		}
		s.Type = c.expr(from.Type)
		c.gaps(from.After)
	case *Var:
		c.gaps(from.Before)
		to = &ast.ValueSpec{
			Names:  c.idents(from.Names),
			Type:   c.expr(from.Type),
			Values: c.exprs(from.Values),
		}
		c.gaps(from.After)
	default:
		panic(fmt.Sprintf("invalid spec: %#v", from))
	}
	return to
}

func (c *syntaxConv) specs(from []Declaration) (to []ast.Spec) {
	to = make([]ast.Spec, len(from))
	for i, d := range from {
		to[i] = c.spec(d)
	}
	return to
}

func (c *syntaxConv) stmt(from Statement) (to ast.Stmt) {
	switch from := from.(type) {
	case nil:
	case *AddAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAddAssign),
			Tok:    token.ADD_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *AndNotAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAndNotAssign),
			Tok:    token.AND_NOT_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *Assign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAssign),
			Tok:    token.ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *BitAndAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAndAssign),
			Tok:    token.AND_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *BitOrAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenOrAssign),
			Tok:    token.OR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *Block:
		c.gaps(from.Before)
		to = &ast.BlockStmt{
			Lbrace: c.add(lenLbrace),
			List:   c.stmts(from.List),
			Rbrace: c.add(lenRbrace),
		}
		c.gaps(from.After)
	case *Break:
		c.gaps(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenBreak),
			Tok:    token.BREAK,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.gaps(from.After)
	case *Case:
		c.gaps(from.Before)
		if from.Comm == nil {
			to = &ast.CaseClause{
				Case:  c.add(lenCase),
				List:  c.exprs(from.List),
				Colon: c.add(lenColon),
				Body:  c.stmts(from.Body),
			}
		} else {
			to = &ast.CommClause{
				Case:  c.add(lenCase),
				Comm:  c.stmt(from.Comm),
				Colon: c.add(lenColon),
				Body:  c.stmts(from.Body),
			}
		}
		c.gaps(from.After)
	case *Continue:
		c.gaps(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenContinue),
			Tok:    token.CONTINUE,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.gaps(from.After)
	case *Dec:
		c.gaps(from.Before)
		to = &ast.IncDecStmt{
			X:      c.expr(from.Left),
			TokPos: c.add(lenDec),
			Tok:    token.DEC,
		}
		c.gaps(from.After)
	case *Defer:
		c.gaps(from.Before)
		to = &ast.DeferStmt{
			Defer: c.add(lenDefer),
			Call:  c.expr(from.Call).(*ast.CallExpr),
		}
		c.gaps(from.After)
	case *Define:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenDefine),
			Tok:    token.DEFINE,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *DivideAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenQuoAssign),
			Tok:    token.QUO_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	// TODO:
	// case *Empty:
	// 	c.gaps(from.Before)
	// 	to = &ast.EmptyStmt{}
	// 	c.gaps(from.After)
	case *Fallthrough:
		c.gaps(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenFallthrough),
			Tok:    token.FALLTHROUGH,
		}
		c.gaps(from.After)
	case *For:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.gaps(from.Before)
		to = &ast.ForStmt{
			For:  c.add(lenFor),
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Post: c.stmt(from.Post),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.gaps(from.After)
	case *Go:
		c.gaps(from.Before)
		to = &ast.GoStmt{
			Go:   c.add(lenGo),
			Call: c.expr(from.Call).(*ast.CallExpr),
		}
		c.gaps(from.After)
	case *Goto:
		c.gaps(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenGoto),
			Tok:    token.GOTO,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.gaps(from.After)
	case *If:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.gaps(from.Before)
		to = &ast.IfStmt{
			If:   c.add(lenIf),
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
			Else: c.stmt(from.Else),
		}
		c.gaps(from.After)
	case *Inc:
		c.gaps(from.Before)
		to = &ast.IncDecStmt{
			X:      c.expr(from.Left),
			TokPos: c.add(lenInc),
			Tok:    token.INC,
		}
		c.gaps(from.After)
	case *Label:
		c.gaps(from.Before)
		to = &ast.LabeledStmt{
			Label: c.expr(from.Label).(*ast.Ident),
			Colon: c.add(lenColon),
			Stmt:  c.stmt(from.Stmt),
		}
		c.gaps(from.After)
	case *MultiplyAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenMulAssign),
			Tok:    token.MUL_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
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
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.gaps(from.Before)
		to = &ast.RangeStmt{
			For:    c.add(lenFor),
			Key:    c.expr(from.Key),
			Value:  c.expr(from.Value),
			TokPos: c.add(l), // TODO: Should not set if Key==nil
			Tok:    t,
			X:      c.expr(from.Container),
			Body:   c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.gaps(from.After)
	case *RemainderAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenRemAssign),
			Tok:    token.REM_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *Return:
		c.gaps(from.Before)
		to = &ast.ReturnStmt{
			Return:  c.add(lenReturn),
			Results: c.exprs(from.Results),
		}
		c.gaps(from.After)
	case *Select:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.gaps(from.Before)
		to = &ast.SelectStmt{
			Select: c.add(lenSelect),
			Body:   c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.gaps(from.After)
	case *Send:
		c.gaps(from.Before)
		to = &ast.SendStmt{
			Chan:  c.expr(from.Left),
			Arrow: c.add(lenArrow),
			Value: c.expr(from.Right),
		}
		c.gaps(from.After)
	case *ShiftLeftAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenShlAssign),
			Tok:    token.SHL_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *ShiftRightAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenShrAssign),
			Tok:    token.SHR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *SubtractAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenSubAssign),
			Tok:    token.SUB_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	case *Switch:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.gaps(from.Before)
		if from.Type == nil {
			to = &ast.SwitchStmt{
				Switch: c.add(lenSwitch),
				Init:   c.stmt(from.Init),
				Tag:    c.expr(from.Value),
				Body:   c.stmt(from.Body).(*ast.BlockStmt),
			}
		} else {
			to = &ast.TypeSwitchStmt{
				Switch: c.add(lenSwitch),
				Init:   c.stmt(from.Init),
				Assign: c.stmt(from.Type),
				Body:   c.stmt(from.Body).(*ast.BlockStmt),
			}
		}
		c.gaps(from.After)
	case *XorAssign:
		c.gaps(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenXorAssign),
			Tok:    token.XOR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.gaps(from.After)
	default:
		if d, ok := from.(Declaration); ok {
			to = &ast.DeclStmt{Decl: c.decl(d)}
		} else if e, ok := from.(Expression); ok {
			to = &ast.ExprStmt{X: c.expr(e)}
		} else {
			panic(fmt.Sprintf("invalid statement: %#v", from))
		}
	}
	return to
}

func (c *syntaxConv) stmts(from []Statement) (to []ast.Stmt) {
	to = make([]ast.Stmt, len(from))
	for i, s := range from {
		to[i] = c.stmt(s)
	}
	return to
}
