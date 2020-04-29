package syntax

// TODO: Rename Space to Offset or Gap, since it can be for delimiters like {a, b, c}?

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
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
	c := newSyntaxConverter()
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

type syntaxConverter struct {
	astFile      *ast.File
	last         token.Pos
	newLineEmpty bool // new line added but nothing added to it yet. TODO: check if true at end and remove last line from token file if true.
	tokenFile    *token.File
	tokenFileSet *token.FileSet
}

func newSyntaxConverter() *syntaxConverter {
	fset := token.NewFileSet()
	return &syntaxConverter{
		astFile:      &ast.File{},
		tokenFile:    fset.AddFile("", -1, int((^uint(0))>>1)),
		tokenFileSet: fset,
	}
}

func (c *syntaxConverter) add(n int) token.Pos {
	var p = c.last + 1
	c.last += token.Pos(n)
	c.newLineEmpty = false
	return p
}

func (c *syntaxConverter) decl(from Declaration) (to ast.Decl) {
	switch from := from.(type) {
	case *Const:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenConst),
			Tok:    token.CONST,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ConstList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenConst)
		g.Tok = token.CONST
		c.markup(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.markup(from.After)
		to = g
	case *Func:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.markup(from.Before)
		funcPos := c.add(lenFunc)
		funcDecl := &ast.FuncDecl{
			Recv: c.node(from.Receiver).(*ast.FieldList),
			Name: c.expr(from.Name).(*ast.Ident),
			Type: &ast.FuncType{
				Func:    funcPos,
				Params:  c.node(from.Params).(*ast.FieldList),
				Results: c.results(from.Results),
			},
		}
		if from.Body != nil {
			funcDecl.Body = c.stmt(from.Body).(*ast.BlockStmt)
		}
		c.markup(from.After)
		to = funcDecl
	case *Import:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenImport),
			Tok:    token.IMPORT,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *ImportList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenImport)
		g.Tok = token.IMPORT
		c.markup(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.markup(from.After)
		to = g
	case *Type:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenType),
			Tok:    token.TYPE,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *TypeList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenType)
		g.Tok = token.TYPE
		c.markup(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.markup(from.After)
		to = g
	case *Var:
		c.markup(from.Before)
		to = &ast.GenDecl{
			TokPos: c.add(lenVar),
			Tok:    token.VAR,
			Specs:  []ast.Spec{c.spec(from)},
		}
		c.markup(from.After)
	case *VarList:
		c.markup(from.Before)
		g := &ast.GenDecl{}
		g.TokPos = c.add(lenVar)
		g.Tok = token.VAR
		c.markup(from.Between)
		g.Lparen = c.add(lenLparen)
		g.Specs = c.specs(from.List)
		g.Rparen = c.add(lenRparen)
		c.markup(from.After)
		to = g
	}
	return to
}

func (c *syntaxConverter) decls(from []Declaration) (to []ast.Decl) {
	to = make([]ast.Decl, len(from))
	for i, d := range from {
		to[i] = c.decl(d)
	}
	return to
}

func (c *syntaxConverter) expr(from Expression) (to ast.Expr) {
	switch from := from.(type) {
	case nil:
	case *Add:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAdd),
			Op:    token.ADD,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *And:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLand),
			Op:    token.LAND,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *AndNot:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAndNot),
			Op:    token.AND_NOT,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Array:
		c.markup(from.Before)
		a := &ast.ArrayType{
			Lbrack: c.add(lenLbrack),
			Len:    c.expr(from.Length),
		}
		c.add(lenRbrack)
		a.Elt = c.expr(from.Element)
		c.markup(from.After)
		to = a
	case *Assert:
		c.markup(from.Before)
		to = &ast.TypeAssertExpr{
			X:      c.expr(from.Value),
			Lparen: c.add(lenLparen),
			Type:   c.expr(from.Type),
			Rparen: c.add(lenRparen),
		}
		c.markup(from.After)
	case *BitAnd:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenAnd),
			Op:    token.AND,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *BitOr:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenOr),
			Op:    token.OR,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Call:
		c.markup(from.Before)
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
				c.markup(e.Before)
				call.Ellipsis = c.add(lenEllipsis)
				args[last] = c.expr(e.Elem)
				c.markup(e.After)
			} else {
				args[last] = c.expr(from.Args[last])
			}
			call.Args = args
		}
		call.Rparen = c.add(lenRparen)
		c.markup(from.After)
		to = call
	case *Chan:
		c.markup(from.Before)
		to = &ast.ChanType{
			Dir:   ast.RECV | ast.SEND,
			Begin: c.add(lenChan + 1),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *ChanIn:
		c.markup(from.Before)
		to = &ast.ChanType{
			Dir:   ast.RECV,
			Begin: c.add(lenChan),
			Arrow: c.add(lenArrow + 1),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *ChanOut:
		c.markup(from.Before)
		var p = c.add(lenChan + lenArrow + 1)
		to = &ast.ChanType{
			Dir:   ast.SEND,
			Begin: p,
			Arrow: p,
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Composite:
		c.markup(from.Before)
		to = &ast.CompositeLit{
			Type:   c.expr(from.Type),
			Lbrace: c.add(lenLbrace),
			Elts:   c.exprs(from.Elts),
			Rbrace: c.add(lenRbrace),
		}
		c.markup(from.After)
	case *Deref:
		c.markup(from.Before)
		to = &ast.StarExpr{
			Star: c.add(lenMul),
			X:    c.expr(from.Left),
		}
		c.markup(from.After)
	case *Divide:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenQuo),
			Op:    token.QUO,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Ellipsis:
		c.markup(from.Before)
		to = &ast.Ellipsis{
			Ellipsis: c.add(lenEllipsis),
			Elt:      c.expr(from.Elem),
		}
		c.markup(from.After)
	case *Equal:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenEql),
			Op:    token.EQL,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Float:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.FLOAT,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Func:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.markup(from.Before)
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
		c.markup(from.After)
	case *Greater:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenGtr),
			Op:    token.GTR,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *GreaterEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenGeq),
			Op:    token.GEQ,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Imag:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.IMAG,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Index:
		c.markup(from.Before)
		to = &ast.IndexExpr{
			X:      c.expr(from.Container),
			Lbrack: c.add(lenLbrack),
			Index:  c.expr(from.Param),
			Rbrack: c.add(lenRbrack),
		}
		c.markup(from.After)
	case *Int:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.INT,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Interface:
		if from.Methods == nil {
			from.Methods = &MethodList{}
		}
		c.markup(from.Before)
		to = &ast.InterfaceType{
			Interface: c.add(lenInterface),
			Methods:   c.node(from.Methods).(*ast.FieldList),
		}
		c.markup(from.After)
	case *KeyValue:
		c.markup(from.Before)
		to = &ast.KeyValueExpr{
			Key:   c.expr(from.Key),
			Colon: c.add(lenColon),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Less:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLss),
			Op:    token.LSS,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *LessEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLeq),
			Op:    token.LEQ,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Map:
		c.markup(from.Before)
		to = &ast.MapType{
			Map:   c.add(lenMap),
			Key:   c.expr(from.Key),
			Value: c.expr(from.Value),
		}
		c.markup(from.After)
	case *Multiply:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenMul),
			Op:    token.MUL,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Name:
		if from == nil { // Branch labels, import aliases
			to = (*ast.Ident)(nil)
		} else {
			c.markup(from.Before)
			to = &ast.Ident{
				NamePos: c.add(len(from.Text)),
				Name:    from.Text,
			}
			c.markup(from.After)
		}
	case *Negate:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenSub),
			Op:    token.SUB,
			X:     c.expr(from.Left),
		}
		c.markup(from.After)
	case *Not:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenNot),
			Op:    token.NOT,
			X:     c.expr(from.Left),
		}
		c.markup(from.After)
	case *NotEqual:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenNeq),
			Op:    token.NEQ,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Or:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenLor),
			Op:    token.LOR,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Paren:
		c.markup(from.Before)
		to = &ast.ParenExpr{
			Lparen: c.add(lenLparen),
			X:      c.expr(from.Left),
			Rparen: c.add(lenRparen),
		}
		c.markup(from.After)
	case *Pointer:
		c.markup(from.Before)
		to = &ast.StarExpr{
			Star: c.add(lenMul),
			X:    c.expr(from.Left),
		}
		c.markup(from.After)
	case *Receive:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenArrow),
			Op:    token.ARROW,
			X:     c.expr(from.Left),
		}
		c.markup(from.After)
	case *Ref:
		c.markup(from.Before)
		to = &ast.UnaryExpr{
			OpPos: c.add(lenAnd),
			Op:    token.AND,
			X:     c.expr(from.Left),
		}
		c.markup(from.After)
	case *Remainder:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenRem),
			Op:    token.REM,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Rune:
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.CHAR,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Selector:
		c.markup(from.Before)
		e := &ast.SelectorExpr{X: c.expr(from.Value)}
		c.add(lenPeriod)
		e.Sel = c.expr(from.Name).(*ast.Ident)
		c.markup(from.After)
		to = e
	case *ShiftLeft:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenShl),
			Op:    token.SHL,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *ShiftRight:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenShr),
			Op:    token.SHR,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Slice:
		c.markup(from.Before)
		to = &ast.SliceExpr{
			X:      c.expr(from.Slice),
			Lbrack: c.add(lenLbrack),
			Low:    c.expr(from.Low),
			High:   c.expr(from.High),
			Max:    c.expr(from.Max),
			Slice3: from.Max != nil,
			Rbrack: c.add(lenRbrack),
		}
		c.markup(from.After)
	case *String:
		if from == nil { // TODO: Why needed?
			return nil
		}
		c.markup(from.Before)
		to = &ast.BasicLit{
			ValuePos: c.add(len(from.Text)),
			Kind:     token.STRING,
			Value:    from.Text,
		}
		c.markup(from.After)
	case *Struct:
		if from.Fields == nil {
			from.Fields = &FieldList{}
		}
		c.markup(from.Before)
		to = &ast.StructType{
			Struct:     c.add(lenStruct),
			Fields:     c.node(from.Fields).(*ast.FieldList),
			Incomplete: false, // TODO: What is this for?
		}
		c.markup(from.After)
	case *Subtract:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenSub),
			Op:    token.SUB,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	case *Xor:
		c.markup(from.Before)
		to = &ast.BinaryExpr{
			X:     c.expr(from.Left),
			OpPos: c.add(lenXor),
			Op:    token.XOR,
			Y:     c.expr(from.Right),
		}
		c.markup(from.After)
	default:
		panic(fmt.Sprintf("invalid expression: %#v", from))
	}
	return to
}

func (c *syntaxConverter) exprs(from []Expression) (to []ast.Expr) {
	to = make([]ast.Expr, len(from))
	for i, e := range from {
		to[i] = c.expr(e)
	}
	return to
}

func (c *syntaxConverter) idents(from []*Name) (to []*ast.Ident) {
	to = make([]*ast.Ident, len(from))
	for i, n := range from {
		to[i] = c.expr(n).(*ast.Ident)
	}
	return to
}

func (c *syntaxConverter) markup(ss []Gap) {
	var cg *ast.CommentGroup
	var lastLine bool // Whether the last item was a line
	for _, s := range ss {
		switch s := s.(type) {
		case *Comment:
			if cg == nil {
				cg = &ast.CommentGroup{}
			}
			cg.List = append(cg.List, &ast.Comment{
				Slash: c.add(len(s.Text)),
				Text:  s.Text,
			})
			// TODO: Add newlines in comment
			lastLine = false
		case *Line:
			if lastLine && cg != nil {
				c.astFile.Comments = append(c.astFile.Comments, cg)
				cg = nil
			}
			c.tokenFile.AddLine(c.tokenFile.Offset(c.add(1) + 1))
			c.newLineEmpty = true
			lastLine = true
		case *Space:
			if s.Count > 0 {
				c.add(s.Count)
			} else {
				c.add(1)
			}
		default:
			panic(fmt.Sprintf("invalid context: %#v", s)) // TODO: Remove
		}
	}
	if cg != nil {
		c.astFile.Comments = append(c.astFile.Comments, cg)
	}
}

func (c *syntaxConverter) results(from *ParamList) (to *ast.FieldList) {
	if from == nil {
		return nil
	}
	parens := len(from.List) != 1 || len(from.List[0].Names) > 0
	c.markup(from.Before)
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
	c.markup(from.After)
	return to
}

func (c *syntaxConverter) node(from Syntax) (to ast.Node) {
	switch from := from.(type) {
	case nil:
	case *Field:
		c.markup(from.Before)
		field := &ast.Field{}
		field.Names = c.idents(from.Names)
		if b, ok := c.expr(from.Tag).(*ast.BasicLit); ok {
			field.Tag = b
		}
		field.Type = c.expr(from.Type)
		c.markup(from.After)
		to = field
	case *FieldList:
		// TODO: Refactor into helper for Field/Method/ParamList
		c.markup(from.Before)
		fieldList := &ast.FieldList{}
		fieldList.Opening = c.add(lenLbrace)
		for _, f := range from.List {
			fieldList.List = append(fieldList.List, c.node(f).(*ast.Field))
		}
		fieldList.Closing = c.add(lenRbrace)
		c.markup(from.After)
		to = fieldList
	case *File:
		c.markup(from.Before)
		c.astFile.Package = c.add(lenPackage)
		c.astFile.Name = c.expr(from.Package).(*ast.Ident)
		c.astFile.Decls = c.decls(from.Decls)
		c.markup(from.After)
		to = c.astFile
	case *Method:
		if from.Params == nil {
			from.Params = &ParamList{}
		}
		c.markup(from.Before)
		to = &ast.Field{
			Names: []*ast.Ident{c.expr(from.Name).(*ast.Ident)},
			Type: &ast.FuncType{
				Params:  c.node(from.Params).(*ast.FieldList),
				Results: c.results(from.Results),
			},
		}
		c.markup(from.After)
	case *MethodList:
		c.markup(from.Before)
		fieldList := &ast.FieldList{}
		fieldList.Opening = c.add(lenLbrace)
		for _, f := range from.List {
			fieldList.List = append(fieldList.List, c.node(f).(*ast.Field))
		}
		fieldList.Closing = c.add(lenRbrace)
		c.markup(from.After)
		to = fieldList
	case *Param:
		c.markup(from.Before)
		to = &ast.Field{
			Names: c.idents(from.Names),
			Type:  c.expr(from.Type),
		}
		c.markup(from.After)
	case *ParamList:
		if from == nil { // Func with no results
			to = (*ast.FieldList)(nil)
		} else {
			c.markup(from.Before)
			fieldList := &ast.FieldList{}
			fieldList.Opening = c.add(lenLparen)
			for _, f := range from.List {
				fieldList.List = append(fieldList.List, c.node(f).(*ast.Field))
			}
			fieldList.Closing = c.add(lenRparen)
			c.markup(from.After)
			to = fieldList
		}
	case *Receiver:
		if from == nil { // Func with no receiver
			to = (*ast.FieldList)(nil)
		} else {
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
		}
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

func (c *syntaxConverter) spec(from Syntax) (to ast.Spec) {
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
			Name:   c.expr(from.Name).(*ast.Ident),
			Path:   c.expr(from.Path).(*ast.BasicLit),
			EndPos: 0, // TODO: Verify this should be 0
		}
		c.markup(from.After)
	case *Type:
		c.markup(from.Before)
		s := &ast.TypeSpec{}
		s.Name = c.expr(from.Name).(*ast.Ident)
		if from.Assign {
			s.Assign = c.add(lenAssign)
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
	default:
		panic(fmt.Sprintf("invalid spec: %#v", from))
	}
	return to
}

func (c *syntaxConverter) specs(from []Declaration) (to []ast.Spec) {
	to = make([]ast.Spec, len(from))
	for i, d := range from {
		to[i] = c.spec(d)
	}
	return to
}

func (c *syntaxConverter) stmt(from Statement) (to ast.Stmt) {
	switch from := from.(type) {
	case nil:
	case *AddAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAddAssign),
			Tok:    token.ADD_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *AndNotAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAndNotAssign),
			Tok:    token.AND_NOT_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Assign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAssign),
			Tok:    token.ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *BitAndAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenAndAssign),
			Tok:    token.AND_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *BitOrAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenOrAssign),
			Tok:    token.OR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Block:
		c.markup(from.Before)
		to = &ast.BlockStmt{
			Lbrace: c.add(lenLbrace),
			List:   c.stmts(from.List),
			Rbrace: c.add(lenRbrace),
		}
		c.markup(from.After)
	case *Break:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenBreak),
			Tok:    token.BREAK,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.markup(from.After)
	case *Case:
		c.markup(from.Before)
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
		c.markup(from.After)
	case *Continue:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenContinue),
			Tok:    token.CONTINUE,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.markup(from.After)
	case *Dec:
		c.markup(from.Before)
		to = &ast.IncDecStmt{
			X:      c.expr(from.Left),
			TokPos: c.add(lenDec),
			Tok:    token.DEC,
		}
		c.markup(from.After)
	case *Defer:
		c.markup(from.Before)
		to = &ast.DeferStmt{
			Defer: c.add(lenDefer),
			Call:  c.expr(from.Call).(*ast.CallExpr),
		}
		c.markup(from.After)
	case *Define:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenDefine),
			Tok:    token.DEFINE,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *DivideAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenQuoAssign),
			Tok:    token.QUO_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	// TODO:
	// case *Empty:
	// 	c.markup(from.Before)
	// 	to = &ast.EmptyStmt{}
	// 	c.markup(from.After)
	case *Fallthrough:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenFallthrough),
			Tok:    token.FALLTHROUGH,
		}
		c.markup(from.After)
	case *For:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.markup(from.Before)
		to = &ast.ForStmt{
			For:  c.add(lenFor),
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Post: c.stmt(from.Post),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *Go:
		c.markup(from.Before)
		to = &ast.GoStmt{
			Go:   c.add(lenGo),
			Call: c.expr(from.Call).(*ast.CallExpr),
		}
		c.markup(from.After)
	case *Goto:
		c.markup(from.Before)
		to = &ast.BranchStmt{
			TokPos: c.add(lenGoto),
			Tok:    token.GOTO,
			Label:  c.expr(from.Label).(*ast.Ident),
		}
		c.markup(from.After)
	case *If:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.markup(from.Before)
		to = &ast.IfStmt{
			If:   c.add(lenIf),
			Init: c.stmt(from.Init),
			Cond: c.expr(from.Cond),
			Body: c.stmt(from.Body).(*ast.BlockStmt),
			Else: c.stmt(from.Else),
		}
		c.markup(from.After)
	case *Inc:
		c.markup(from.Before)
		to = &ast.IncDecStmt{
			X:      c.expr(from.Left),
			TokPos: c.add(lenInc),
			Tok:    token.INC,
		}
		c.markup(from.After)
	case *Label:
		c.markup(from.Before)
		to = &ast.LabeledStmt{
			Label: c.expr(from.Label).(*ast.Ident),
			Colon: c.add(lenColon),
			Stmt:  c.stmt(from.Stmt),
		}
		c.markup(from.After)
	case *MultiplyAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenMulAssign),
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
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.markup(from.Before)
		to = &ast.RangeStmt{
			For:    c.add(lenFor),
			Key:    c.expr(from.Key),
			Value:  c.expr(from.Value),
			TokPos: c.add(l), // TODO: Should not set if Key==nil
			Tok:    t,
			X:      c.expr(from.Container),
			Body:   c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *RemainderAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenRemAssign),
			Tok:    token.REM_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Return:
		c.markup(from.Before)
		to = &ast.ReturnStmt{
			Return:  c.add(lenReturn),
			Results: c.exprs(from.Results),
		}
		c.markup(from.After)
	case *Select:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.markup(from.Before)
		to = &ast.SelectStmt{
			Select: c.add(lenSelect),
			Body:   c.stmt(from.Body).(*ast.BlockStmt),
		}
		c.markup(from.After)
	case *Send:
		c.markup(from.Before)
		to = &ast.SendStmt{
			Chan:  c.expr(from.Left),
			Arrow: c.add(lenArrow),
			Value: c.expr(from.Right),
		}
		c.markup(from.After)
	case *ShiftLeftAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenShlAssign),
			Tok:    token.SHL_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *ShiftRightAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenShrAssign),
			Tok:    token.SHR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *SubtractAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenSubAssign),
			Tok:    token.SUB_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
	case *Switch:
		if from.Body == nil {
			from.Body = &Block{}
		}
		c.markup(from.Before)
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
		c.markup(from.After)
	case *XorAssign:
		c.markup(from.Before)
		to = &ast.AssignStmt{
			Lhs:    c.exprs(from.Left),
			TokPos: c.add(lenXorAssign),
			Tok:    token.XOR_ASSIGN,
			Rhs:    c.exprs(from.Right),
		}
		c.markup(from.After)
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

func (c *syntaxConverter) stmts(from []Statement) (to []ast.Stmt) {
	to = make([]ast.Stmt, len(from))
	for i, s := range from {
		to[i] = c.stmt(s)
	}
	return to
}
