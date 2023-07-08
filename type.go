package turbine

import (
	"fmt"
	"go/ast"
	"go/types"
)

var typesyntax = map[string]*Type{}

func varsStruct(t *types.Struct) []*types.Var {
	var vs []*types.Var

	for i, n := 0, t.NumFields(); i < n; i++ {
		vs = append(vs, t.Field(i))
	}

	return vs
}

func varsTuple(t *types.Tuple) []*types.Var {
	var vs []*types.Var

	for i, n := 0, t.Len(); i < n; i++ {
		vs = append(vs, t.At(i))
	}

	return vs
}

// Type is type information. TODO.
type Type struct {
	Context        *Type
	Element        *Type
	Fields         []*VarGroup
	FieldsJoined   []*VarGroup
	FieldsOrdered  []*VarGroup
	FieldsSplit    []*VarGroup
	HasAdd         bool
	HasAnd         bool
	HasAppend      bool
	HasBitAnd      bool
	HasBitClear    bool
	HasBitLeft     bool
	HasBitOr       bool
	HasBitRight    bool
	HasBitXor      bool
	HasCall        bool
	HasCap         bool
	HasClose       bool
	HasConcat      bool
	HasCopyDst     bool
	HasCopySrc     bool
	HasDec         bool
	HasDel         bool
	HasDeref       bool
	HasDiv         bool
	HasElemRef     bool
	HasEq          bool
	HasGreater     bool
	HasGreaterEq   bool
	HasImag        bool
	HasInc         bool
	HasIndex       bool
	HasLen         bool
	HasLess        bool
	HasLessEq      bool
	HasLit         bool
	HasMake        bool
	HasMod         bool
	HasMul         bool
	HasNeg         bool
	HasNeq         bool
	HasNot         bool
	HasOr          bool
	HasRange       bool
	HasRangeElem   bool
	HasReal        bool
	HasReceive     bool
	HasSend        bool
	HasSlice       bool
	HasSub         bool
	HasVariadic    bool
	Index          *Type
	IsArray        bool
	IsBool         bool
	IsChan         bool
	IsCmp          bool
	IsComplex      bool
	IsConst        bool
	IsDecl         bool
	IsFloat        bool
	IsFunc         bool
	IsInt          bool
	IsInterface    bool
	IsMap          bool
	IsNamed        bool
	IsNil          bool
	IsNum          bool
	IsOrd          bool
	IsPtr          bool
	IsRune         bool
	IsSigned       bool
	IsSlice        bool
	IsString       bool
	IsStruct       bool
	IsUnsafe       bool
	IsUntyped      bool
	Len            int64
	Methods        []*Decl
	Name           string
	Package        string
	Params         []*VarGroup
	ParamsJoined   []*VarGroup
	ParamsOrdered  []*VarGroup
	ParamsSplit    []*VarGroup
	Path           string
	Qualified      string
	Qualifier      string
	Receiver       *Var
	Results        []*VarGroup
	ResultsJoined  []*VarGroup
	ResultsOrdered []*VarGroup
	ResultsSplit   []*VarGroup
	Syntax         string
	Underlying     *Type
	// TODO: More bitwise ops?
	// TODO: Zero value string
}

func newType(t types.Type, n ast.Node) *Type {
	if t == nil {
		return nil
	}

	if _, ok := n.(*ast.SelectorExpr); ok {
		return newType(t, nil)
	}

	var s = types.TypeString(t, nil)

	if t, ok := typesyntax[s]; ok {
		return t
	}

	var x = &Type{Syntax: s}

	typesyntax[s] = x

	switch t := t.(type) {
	case *types.Array:
		var e ast.Expr

		if n != nil {
			e = n.(*ast.ArrayType).Elt
		}

		x.Element = newType(t.Elem(), e)
		x.Element.Context = x
		x.Index = newType(types.Typ[types.Int], nil)
		x.Index.Context = x
		x.IsArray = true
		x.Len = t.Len()

	case *types.Basic:
		switch t {
		case types.Typ[types.Bool], types.Typ[types.UntypedBool]:
			x.HasAnd = true
			x.HasNot = true
			x.HasOr = true
			x.IsBool = true

		case types.Typ[types.Complex64], types.Typ[types.Complex128], types.Typ[types.UntypedComplex]:
			x.HasImag = true
			x.HasReal = true
			x.IsComplex = true

		case types.Typ[types.Float32], types.Typ[types.Float64], types.Typ[types.UntypedFloat]:
			x.IsFloat = true

		case types.Typ[types.Byte], types.Typ[types.Int], types.Typ[types.Int8], types.Typ[types.Int16], types.Typ[types.Int32], types.Typ[types.Int64], types.Typ[types.Rune], types.Typ[types.Uint], types.Typ[types.Uint8], types.Typ[types.Uint16], types.Typ[types.Uint32], types.Typ[types.Uint64], types.Typ[types.Uintptr], types.Typ[types.UntypedInt], types.Typ[types.UntypedRune]:
			x.HasBitAnd = true
			x.HasBitClear = true
			x.HasBitLeft = true
			x.HasBitOr = true
			x.HasBitRight = true
			x.HasBitXor = true
			x.HasMod = true
			x.IsInt = true

			switch t {
			case types.Typ[types.Rune], types.Typ[types.UntypedRune]:
				x.IsRune = true
			}

		case types.Typ[types.String], types.Typ[types.UntypedString]:
			x.Element = newType(types.Typ[types.Byte], nil)
			x.Element.Context = x
			x.HasConcat = true
			x.Index = newType(types.Typ[types.Int], nil)
			x.Index.Context = x
			x.IsString = true

		case types.Typ[types.UnsafePointer]:
			x.IsUnsafe = true

		case types.Typ[types.UntypedNil]:
			x.IsNil = true
		}

	case *types.Chan:
		var e ast.Expr

		if n != nil {
			e = n.(*ast.ChanType).Value
		}

		x.Element = newType(t.Elem(), e)
		x.Element.Context = x
		x.IsChan = true

		switch t.Dir() {
		case types.SendRecv:
			x.HasClose = true
			x.HasMake = true
			x.HasReceive = true
			x.HasSend = true

		case types.RecvOnly:
			x.HasReceive = true

		case types.SendOnly:
			x.HasClose = true
			x.HasSend = true
		}

	case *types.Interface:
		x.IsInterface = true

		for i, c := 0, t.NumMethods(); i < c; i++ {
			var mt = t.Method(i)
			var mn ast.Node
			var d = Decl{Ident: newIdent(mt.Name()), IsFunc: true}

			if n != nil {
				mn = n.(*ast.InterfaceType).Methods.List[i].Type
				d.Doc, d.DocLines = doc(n.(*ast.InterfaceType).Methods.List[i].Doc)
			}

			d.Type = newType(mt.Type(), mn)
			x.Methods = append(x.Methods, &d)
		}

	case *types.Map:
		var k, v ast.Expr

		if n != nil {
			k = n.(*ast.MapType).Key
			v = n.(*ast.MapType).Value
		}

		x.Element = newType(t.Elem(), v)
		x.Element.Context = x
		x.HasDel = true
		x.Index = newType(t.Key(), k)
		x.Index.Context = x
		x.IsMap = true

	case *types.Named:
		var n = t.Obj()

		x.IsNamed = true
		x.Name = n.Name()

		if p := n.Pkg(); p == nil {
			x.Qualified = x.Name
			x.Syntax = types.TypeString(t, nil)
		} else {
			x.Package = p.Name()
			x.Path = p.Path()
			x.Qualified = fmt.Sprintf("%v.%v", x.Package, x.Name)
			x.Qualifier = x.Package
			x.Syntax = x.Qualified
		}

	case *types.Pointer:
		var e ast.Expr

		if n != nil {
			e = n.(*ast.StarExpr).X
		}

		x.Element = newType(t.Elem(), e)
		x.Element.Context = x
		x.HasDeref = true
		x.IsPtr = true

	case *types.Signature:
		x.HasCall = true
		x.HasVariadic = t.Variadic()
		x.IsFunc = true

		var pfs, rfs []*ast.Field

		if n != nil {
			var f = n.(*ast.FuncType)

			if f.Params != nil {
				pfs = f.Params.List
			}

			if f.Results != nil {
				rfs = f.Results.List
			}
		}

		x.ParamsJoined, x.ParamsOrdered, x.Params, x.ParamsSplit = varGroups(varsTuple(t.Params()), pfs)
		x.ResultsJoined, x.ResultsOrdered, x.Results, x.ResultsSplit = varGroups(varsTuple(t.Results()), rfs)

		if r := t.Recv(); r != nil {
			var i *Ident
			var e ast.Expr

			if n != nil {
				var f = n.(*ast.FuncDecl).Recv.List[0]

				e = f.Type

				if len(f.Names) > 0 {
					i = newIdent(f.Names[0].Name)
				}
			}

			x.Receiver = &Var{Ident: i, Type: newType(r.Type(), e)}
		}

	case *types.Slice:
		var e ast.Expr

		if n != nil {
			if a, ok := n.(*ast.ArrayType); ok {
				e = a.Elt
			} else {
				e = n.(*ast.Ellipsis).Elt
			}

		}

		x.Element = newType(t.Elem(), e)
		x.Element.Context = x
		x.HasAppend = true
		x.Index = newType(types.Typ[types.Int], nil)
		x.Index.Context = x
		x.IsSlice = true

	case *types.Struct:
		var fs []*ast.Field

		if n != nil {
			fs = n.(*ast.StructType).Fields.List
		}

		x.FieldsJoined, x.FieldsOrdered, x.Fields, x.FieldsSplit = varGroups(varsStruct(t), fs)
		x.IsStruct = true

	default:
		panic(t)
	}

	switch {
	case x.IsComplex, x.IsFloat, x.IsInt:
		x.HasAdd = true
		x.HasDec = true
		x.HasDiv = true
		x.HasInc = true
		x.HasMul = true
		x.HasSub = true
		x.IsNum = true

		switch t {
		case types.Typ[types.Uint], types.Typ[types.Uint8], types.Typ[types.Uint16], types.Typ[types.Uint32], types.Typ[types.Uint64], types.Typ[types.Uintptr]:

		default:
			x.HasNeg = true
			x.IsSigned = true
		}
	}

	switch {
	case x.IsFunc, x.IsMap, x.IsSlice:

	default:
		x.HasEq = true
		x.HasNeq = true
		x.IsCmp = true
	}

	switch {
	case x.IsBool, x.IsNum, x.IsRune, x.IsString:
		x.IsConst = true
	}

	switch {
	case x.IsInt, x.IsFloat, x.IsString:
		x.HasGreater = true
		x.HasGreaterEq = true
		x.HasLess = true
		x.HasLessEq = true
		x.IsOrd = true
	}

	switch t {
	case types.Typ[types.UntypedBool], types.Typ[types.UntypedComplex], types.Typ[types.UntypedFloat], types.Typ[types.UntypedInt], types.Typ[types.UntypedNil], types.Typ[types.UntypedRune], types.Typ[types.UntypedString]:
		x.IsUntyped = true
	}

	switch {
	case x.IsArray, x.IsChan, x.IsMap, x.IsSlice:
		x.HasCap = true
	}

	switch {
	case x.IsArray, x.IsSlice:
		x.HasCopyDst = true
	}

	switch {
	case x.IsArray, x.IsSlice, x.IsString:
		x.HasCopySrc = true
	}

	switch {
	case x.IsArray, x.IsSlice:
		x.HasElemRef = true
	}

	switch {
	case x.IsArray, x.IsMap, x.IsSlice, x.IsString:
		x.HasIndex = true
	}

	switch {
	case x.IsArray, x.IsChan, x.IsMap, x.IsSlice, x.IsString:
		x.HasLen = true
	}

	switch {
	case x.IsArray, x.IsMap, x.IsSlice, x.IsStruct:
		x.HasLit = true
	}

	switch {
	case x.IsMap, x.IsSlice:
		x.HasMake = true
	}

	switch {
	case x.IsArray, x.IsSlice, x.IsString:
		x.HasSlice = true
	}

	switch t {
	case types.Typ[types.Bool], types.Typ[types.UntypedBool]:
		x.IsNamed = true
		x.Name = "bool"
		x.Syntax = "bool"

	case types.Typ[types.Complex64]:
		x.IsNamed = true
		x.Name = "complex64"
		x.Syntax = "complex64"

	case types.Typ[types.Complex128], types.Typ[types.UntypedComplex]:
		x.IsNamed = true
		x.Name = "complex128"
		x.Syntax = "complex128"

	case types.Typ[types.Float32]:
		x.IsNamed = true
		x.Name = "float32"
		x.Syntax = "float32"

	case types.Typ[types.Float64], types.Typ[types.UntypedFloat]:
		x.IsNamed = true
		x.Name = "float64"
		x.Syntax = "float64"

	case types.Typ[types.Byte]:
		x.IsNamed = true
		x.Name = "byte"
		x.Syntax = "byte"

	case types.Typ[types.Int], types.Typ[types.UntypedInt]:
		x.IsNamed = true
		x.Name = "int"
		x.Syntax = "int"

	case types.Typ[types.Int8]:
		x.IsNamed = true
		x.Name = "int8"
		x.Syntax = "int8"

	case types.Typ[types.Int16]:
		x.IsNamed = true
		x.Name = "int16"
		x.Syntax = "int16"

	case types.Typ[types.Int32]:
		x.IsNamed = true
		x.Name = "int32"
		x.Syntax = "int32"

	case types.Typ[types.Int64]:
		x.IsNamed = true
		x.Name = "int64"
		x.Syntax = "int64"

	case types.Typ[types.Rune]:
		x.IsNamed = true
		x.Name = "rune"
		x.Syntax = "rune"

	case types.Typ[types.Uint]:
		x.IsNamed = true
		x.Name = "uint"
		x.Syntax = "uint"

	case types.Typ[types.Uint8]:
		x.IsNamed = true
		x.Name = "uint8"
		x.Syntax = "uint8"

	case types.Typ[types.Uint16]:
		x.IsNamed = true
		x.Name = "uint16"
		x.Syntax = "uint16"

	case types.Typ[types.Uint32]:
		x.IsNamed = true
		x.Name = "uint32"
		x.Syntax = "uint32"

	case types.Typ[types.Uint64]:
		x.IsNamed = true
		x.Name = "uint64"
		x.Syntax = "uint64"

	case types.Typ[types.Uintptr]:
		x.IsNamed = true
		x.Name = "uintptr"
		x.Syntax = "uintptr"

	case types.Typ[types.Rune], types.Typ[types.UntypedRune]:
		x.IsNamed = true
		x.Name = "rune"
		x.Syntax = "rune"
	}

	switch t.(type) {
	case *types.Named:

	default:
		x.Qualified = x.Name
	}

	if u := t.Underlying(); u != t {
		x.Underlying = newType(u, nil)
		x.Underlying.Context = x
		x.has(x.Underlying)
	}

	return x
}

// TODO
func (t *Type) String() string {
	return t.Syntax
}

func (t *Type) has(from *Type) {
	t.Fields = from.Fields
	t.FieldsJoined = from.FieldsJoined
	t.FieldsOrdered = from.FieldsOrdered
	t.FieldsSplit = from.FieldsSplit
	t.HasAdd = from.HasAdd
	t.HasAnd = from.HasAnd
	t.HasAppend = from.HasAppend
	t.HasBitAnd = from.HasBitAnd
	t.HasBitClear = from.HasBitClear
	t.HasBitLeft = from.HasBitLeft
	t.HasBitOr = from.HasBitOr
	t.HasBitRight = from.HasBitRight
	t.HasBitXor = from.HasBitXor
	t.HasCall = from.HasCall
	t.HasCap = from.HasCap
	t.HasClose = from.HasClose
	t.HasConcat = from.HasConcat
	t.HasCopyDst = from.HasCopyDst
	t.HasCopySrc = from.HasCopySrc
	t.HasDec = from.HasDec
	t.HasDel = from.HasDel
	t.HasDeref = from.HasDeref
	t.HasDiv = from.HasDiv
	t.HasElemRef = from.HasElemRef
	t.HasEq = from.HasEq
	t.HasGreater = from.HasGreater
	t.HasGreaterEq = from.HasGreaterEq
	t.HasImag = from.HasImag
	t.HasInc = from.HasInc
	t.HasIndex = from.HasIndex
	t.HasLen = from.HasLen
	t.HasLess = from.HasLess
	t.HasLessEq = from.HasLessEq
	t.HasLit = from.HasLit
	t.HasMake = from.HasMake
	t.HasMod = from.HasMod
	t.HasMul = from.HasMul
	t.HasNeg = from.HasNeg
	t.HasNeq = from.HasNeq
	t.HasNot = from.HasNot
	t.HasOr = from.HasOr
	t.HasRange = from.HasRange
	t.HasRangeElem = from.HasRangeElem
	t.HasReal = from.HasReal
	t.HasReceive = from.HasReceive
	t.HasSend = from.HasSend
	t.HasSlice = from.HasSlice
	t.HasSub = from.HasSub
	t.HasVariadic = from.HasVariadic
}
