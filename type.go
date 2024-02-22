package turbine

import "go/types"

type Typ struct {
	t types.Type
}

func (t *Typ) IsBool() bool {
	return t.t == types.Typ[types.Bool]
}

func (t *Typ) IsBoolish() bool {
	return t.IsBool() || t.IsUntypedBool()
}

func (t *Typ) IsInt() bool {
	return t.t == types.Typ[types.Int]
}

func (t *Typ) IsIntish() bool {
	return t.IsInt() || t.IsUntypedInt()
}

func (t *Typ) IsInt8() bool {
	return t.t == types.Typ[types.Int8]
}

func (t *Typ) IsInt16() bool {
	return t.t == types.Typ[types.Int16]
}

func (t *Typ) IsInt32() bool {
	return t.t == types.Typ[types.Int32]
}

func (t *Typ) IsInt64() bool {
	return t.t == types.Typ[types.Int64]
}

func (t *Typ) IsSigned() bool {
	return t.IsInt() ||
		t.IsInt8() ||
		t.IsInt16() ||
		t.IsInt32() ||
		t.IsInt64()
}

func (t *Typ) IsUint() bool {
	return t.t == types.Typ[types.Uint]
}

func (t *Typ) IsUint8() bool {
	return t.t == types.Typ[types.Uint8]
}

func (t *Typ) IsUint16() bool {
	return t.t == types.Typ[types.Uint16]
}

func (t *Typ) IsUint32() bool {
	return t.t == types.Typ[types.Uint32]
}

func (t *Typ) IsUint64() bool {
	return t.t == types.Typ[types.Uint64]
}

func (t *Typ) IsUintptr() bool {
	return t.t == types.Typ[types.Uintptr]
}

func (t *Typ) IsUnsigned() bool {
	return t.IsUint() ||
		t.IsUint8() ||
		t.IsUint16() ||
		t.IsUint32() ||
		t.IsUint64() ||
		t.IsUintptr()
}

func (t *Typ) IsByte() bool {
	return t.t == types.Typ[types.Byte]
}

func (t *Typ) IsRune() bool {
	return t.t == types.Typ[types.Rune]
}

func (t *Typ) IsRuneish() bool {
	return t.IsRune() || t.IsUntypedRune()
}

func (t *Typ) IsInteger() bool {
	return t.IsInt() ||
		t.IsInt8() ||
		t.IsInt16() ||
		t.IsInt32() ||
		t.IsInt64() ||
		t.IsUint() ||
		t.IsUint8() ||
		t.IsUint16() ||
		t.IsUint32() ||
		t.IsUint64() ||
		t.IsUintptr()
}

func (t *Typ) IsFloat32() bool {
	return t.t == types.Typ[types.Float32]
}

func (t *Typ) IsFloat64() bool {
	return t.t == types.Typ[types.Float64]
}

func (t *Typ) IsFloat() bool {
	return t.IsFloat32() || t.IsFloat64()
}

func (t *Typ) IsFloatish() bool {
	return t.IsFloat32() || t.IsFloat64() || t.IsUntypedFloat()
}

func (t *Typ) IsComplex64() bool {
	return t.t == types.Typ[types.Complex64]
}

func (t *Typ) IsComplex128() bool {
	return t.t == types.Typ[types.Complex128]
}

func (t *Typ) IsComplex() bool {
	return t.IsComplex64() || t.IsComplex128()
}

func (t *Typ) IsComplexish() bool {
	return t.IsComplex64() || t.IsComplex128() || t.IsUntypedComplex()
}

func (t *Typ) IsNumber() bool {
	return t.IsInteger() || t.IsFloat() || t.IsComplex()
}

func (t *Typ) IsString() bool {
	return t.t == types.Typ[types.String]
}

func (t *Typ) IsStringish() bool {
	return t.IsString() || t.IsUntypedString()
}

func (t *Typ) IsUnsafePointer() bool {
	return t.t == types.Typ[types.UnsafePointer]
}

func (t *Typ) IsUntypedBool() bool {
	return t.t == types.Typ[types.UntypedBool]
}

func (t *Typ) IsUntypedInt() bool {
	return t.t == types.Typ[types.UntypedInt]
}

func (t *Typ) IsUntypedRune() bool {
	return t.t == types.Typ[types.UntypedRune]
}

func (t *Typ) IsUntypedFloat() bool {
	return t.t == types.Typ[types.UntypedFloat]
}

func (t *Typ) IsUntypedComplex() bool {
	return t.t == types.Typ[types.UntypedComplex]
}

func (t *Typ) IsUntypedString() bool {
	return t.t == types.Typ[types.UntypedString]
}

func (t *Typ) IsUntypedNil() bool {
	return t.t == types.Typ[types.UntypedNil]
}

func (t *Typ) IsUntyped() bool {
	return t.IsUntypedBool() ||
		t.IsUntypedInt() ||
		t.IsUntypedRune() ||
		t.IsUntypedFloat() ||
		t.IsUntypedComplex() ||
		t.IsUntypedString() ||
		t.IsUntypedNil()
}

func (t *Typ) IsArray() bool {
	_, ok := t.t.(*types.Array)
	return ok
}

func (t *Typ) Size() int64 {
	return t.t.(*types.Array).Len()
}

func (t *Typ) Key() *Typ {
	switch x := t.t.(type) {
	case *types.Array:
		return &Typ{t: types.Typ[types.Int]}
	case *types.Basic:
		switch x {
		case types.Typ[types.String], types.Typ[types.UntypedString]:
			return &Typ{t: types.Typ[types.Int]}
		default:
			panic(x)
		}
	case *types.Map:
		return &Typ{t: x.Key()}
	case *types.Slice:
		return &Typ{t: types.Typ[types.Int]}
	default:
		panic(x)
	}
}

func (t *Typ) Value() *Typ {
	switch x := t.t.(type) {
	case *types.Array:
		return &Typ{t: x.Elem()}
	case *types.Basic:
		switch x {
		case types.Typ[types.String], types.Typ[types.UntypedString]:
			return &Typ{t: types.Typ[types.Byte]}
		default:
			panic(x)
		}
	case *types.Chan:
		return &Typ{t: x.Elem()}
	case *types.Map:
		return &Typ{t: x.Elem()}
	case *types.Pointer:
		return &Typ{t: x.Elem()}
	case *types.Slice:
		return &Typ{t: x.Elem()}
	default:
		panic(x)
	}
}

func (t *Typ) IsChannel() bool {
	_, ok := t.t.(*types.Chan)
	return ok
}

func (t *Typ) HasReceive() bool {
	d := t.t.(*types.Chan).Dir()
	return d == types.SendRecv || d == types.RecvOnly
}

func (t *Typ) HasSend() bool {
	d := t.t.(*types.Chan).Dir()
	return d == types.SendRecv || d == types.SendOnly
}

func (t *Typ) IsInterface() bool {
	_, ok := t.t.(*types.Interface)
	return ok
}

func (t *Typ) Methods() []*Var {
	switch x := t.t.(type) {
	case *types.Interface:
		n := x.NumMethods()
		vs := make([]*Var, n)
		for i := 0; i < n; i++ {
			vs[i] = &Var{Type: &Typ{t: x.Method(i).Type()}}
		}
		return vs
	case *types.Named:
		n := x.NumMethods()
		vs := make([]*Var, n)
		for i := 0; i < n; i++ {
			vs[i] = &Var{Type: &Typ{t: x.Method(i).Type()}}
		}
		return vs
	}
	panic(t.t)
}

func (t *Typ) IsMap() bool {
	_, ok := t.t.(*types.Map)
	return ok
}

func (t *Typ) IsNamed() bool {
	_, ok := t.t.(*types.Named)
	return ok
}

func (t *Typ) IsAlias() bool {
	return t.t.(*types.Named).Obj().IsAlias()
}

func (t *Typ) Name() Name {
	return Name(t.t.(*types.Named).Obj().Name())
}

func (t *Typ) Origin() *Typ {
	return &Typ{t: t.t.(*types.Named).Origin()}
}

func (t *Typ) TypeArgs() []*Typ {
	args := t.t.(*types.Named).TypeArgs()
	n := args.Len()
	ts := make([]*Typ, n)
	for i := 0; i < n; i++ {
		ts[i] = &Typ{t: args.At(i)}
	}
	return ts
}

func (t *Typ) TypeParams() []*Typ {
	switch x := t.t.(type) {
	case *types.Named:
		params := x.TypeParams()
		n := params.Len()
		ts := make([]*Typ, n)
		for i := 0; i < n; i++ {
			ts[i] = &Typ{t: params.At(i)}
		}
		return ts
	case *types.Signature:
		params := t.t.(*types.Signature).TypeParams()
		n := params.Len()
		ts := make([]*Typ, n)
		for i := 0; i < n; i++ {
			ts[i] = &Typ{t: params.At(i)}
		}
		return ts
	}
	panic(t.t)
}

func (t *Typ) Constraint() *Typ {
	return &Typ{t.t.(*types.TypeParam).Constraint()}
}

func (t *Typ) PackageName() string {
	return t.t.(*types.Named).Obj().Pkg().Name()
}

func (t *Typ) PackagePath() string {
	return t.t.(*types.Named).Obj().Pkg().Path()
}

func (t *Typ) Underlying() *Typ {
	a, b := t.t, t.t.Underlying()
	for a != b {
		a = b
		b = b.Underlying()
	}
	return &Typ{t: b}
}

func (t *Typ) IsPointer() bool {
	_, ok := t.t.(*types.Pointer)
	return ok
}

func (t *Typ) IsSlice() bool {
	_, ok := t.t.(*types.Slice)
	return ok
}

func (t *Typ) IsStruct() bool {
	_, ok := t.t.(*types.Struct)
	return ok
}

func convertVar(v *types.Var) *Var {
	return &Var{Name: Name(v.Name()), Type: &Typ{t: v.Type()}}
}

func convertVars(vs []*types.Var) []*Var {
	vs2 := make([]*Var, len(vs))
	for i, v := range vs {
		vs2[i] = convertVar(v)
	}
	return vs2
}

func structVars(t *types.Struct) []*types.Var {
	n := t.NumFields()
	vs := make([]*types.Var, n)
	for i := 0; i < n; i++ {
		vs[i] = t.Field(i)
	}
	return vs
}

func (t *Typ) Fields() []*Var {
	return convertVars(structVars(t.t.(*types.Struct)))
}

func (t *Typ) IsFunc() bool {
	_, ok := t.t.(*types.Signature)
	return ok
}

func (t *Typ) Receiver() *Var {
	return convertVar(t.t.(*types.Signature).Recv())
}

func (t *Typ) ReceiverTypeParams() []*Typ {
	params := t.t.(*types.Signature).RecvTypeParams()
	n := params.Len()
	ts := make([]*Typ, n)
	for i := 0; i < n; i++ {
		ts[i] = &Typ{t: params.At(i)}
	}
	return ts
}

func tupleVars(t *types.Tuple) []*types.Var {
	n := t.Len()
	vs := make([]*types.Var, n)
	for i := 0; i < n; i++ {
		vs[i] = t.At(i)
	}
	return vs
}

func (t *Typ) Params() []*Var {
	return convertVars(tupleVars(t.t.(*types.Signature).Params()))
}

func (t *Typ) Results() []*Var {
	return convertVars(tupleVars(t.t.(*types.Signature).Results()))
}

func (t *Typ) IsVariadic() bool {
	return t.t.(*types.Signature).Variadic()
}

func (t *Typ) String() string {
	// TODO: Handle qualified identifiers.
	return types.TypeString(t.t, nil)
}
