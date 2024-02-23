package turbine

import "go/types"

type Type struct {
	t types.Type
}

func (t *Type) IsBool() bool {
	return t.t == types.Typ[types.Bool]
}

func (t *Type) IsBoolish() bool {
	return t.IsBool() || t.IsUntypedBool()
}

func (t *Type) IsInt() bool {
	return t.t == types.Typ[types.Int]
}

func (t *Type) IsIntish() bool {
	return t.IsInt() || t.IsUntypedInt()
}

func (t *Type) IsInt8() bool {
	return t.t == types.Typ[types.Int8]
}

func (t *Type) IsInt16() bool {
	return t.t == types.Typ[types.Int16]
}

func (t *Type) IsInt32() bool {
	return t.t == types.Typ[types.Int32]
}

func (t *Type) IsInt64() bool {
	return t.t == types.Typ[types.Int64]
}

func (t *Type) IsSigned() bool {
	return t.IsInt() ||
		t.IsInt8() ||
		t.IsInt16() ||
		t.IsInt32() ||
		t.IsInt64()
}

func (t *Type) IsUint() bool {
	return t.t == types.Typ[types.Uint]
}

func (t *Type) IsUint8() bool {
	return t.t == types.Typ[types.Uint8]
}

func (t *Type) IsUint16() bool {
	return t.t == types.Typ[types.Uint16]
}

func (t *Type) IsUint32() bool {
	return t.t == types.Typ[types.Uint32]
}

func (t *Type) IsUint64() bool {
	return t.t == types.Typ[types.Uint64]
}

func (t *Type) IsUintptr() bool {
	return t.t == types.Typ[types.Uintptr]
}

func (t *Type) IsUnsigned() bool {
	return t.IsUint() ||
		t.IsUint8() ||
		t.IsUint16() ||
		t.IsUint32() ||
		t.IsUint64() ||
		t.IsUintptr()
}

func (t *Type) IsByte() bool {
	return t.t == types.Typ[types.Byte]
}

func (t *Type) IsRune() bool {
	return t.t == types.Typ[types.Rune]
}

func (t *Type) IsRuneish() bool {
	return t.IsRune() || t.IsUntypedRune()
}

func (t *Type) IsInteger() bool {
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

func (t *Type) IsFloat32() bool {
	return t.t == types.Typ[types.Float32]
}

func (t *Type) IsFloat64() bool {
	return t.t == types.Typ[types.Float64]
}

func (t *Type) IsFloat() bool {
	return t.IsFloat32() || t.IsFloat64()
}

func (t *Type) IsFloatish() bool {
	return t.IsFloat32() || t.IsFloat64() || t.IsUntypedFloat()
}

func (t *Type) IsComplex64() bool {
	return t.t == types.Typ[types.Complex64]
}

func (t *Type) IsComplex128() bool {
	return t.t == types.Typ[types.Complex128]
}

func (t *Type) IsComplex() bool {
	return t.IsComplex64() || t.IsComplex128()
}

func (t *Type) IsComplexish() bool {
	return t.IsComplex64() || t.IsComplex128() || t.IsUntypedComplex()
}

func (t *Type) IsNumber() bool {
	return t.IsInteger() || t.IsFloat() || t.IsComplex()
}

func (t *Type) IsString() bool {
	return t.t == types.Typ[types.String]
}

func (t *Type) IsStringish() bool {
	return t.IsString() || t.IsUntypedString()
}

func (t *Type) IsUnsafePointer() bool {
	return t.t == types.Typ[types.UnsafePointer]
}

func (t *Type) IsUntypedBool() bool {
	return t.t == types.Typ[types.UntypedBool]
}

func (t *Type) IsUntypedInt() bool {
	return t.t == types.Typ[types.UntypedInt]
}

func (t *Type) IsUntypedRune() bool {
	return t.t == types.Typ[types.UntypedRune]
}

func (t *Type) IsUntypedFloat() bool {
	return t.t == types.Typ[types.UntypedFloat]
}

func (t *Type) IsUntypedComplex() bool {
	return t.t == types.Typ[types.UntypedComplex]
}

func (t *Type) IsUntypedString() bool {
	return t.t == types.Typ[types.UntypedString]
}

func (t *Type) IsUntypedNil() bool {
	return t.t == types.Typ[types.UntypedNil]
}

func (t *Type) IsUntyped() bool {
	return t.IsUntypedBool() ||
		t.IsUntypedInt() ||
		t.IsUntypedRune() ||
		t.IsUntypedFloat() ||
		t.IsUntypedComplex() ||
		t.IsUntypedString() ||
		t.IsUntypedNil()
}

func (t *Type) IsArray() bool {
	_, ok := t.t.(*types.Array)
	return ok
}

func (t *Type) Size() int64 {
	return t.t.(*types.Array).Len()
}

func (t *Type) Key() *Type {
	switch x := t.t.(type) {
	case *types.Array:
		return &Type{t: types.Typ[types.Int]}
	case *types.Basic:
		switch x {
		case types.Typ[types.String], types.Typ[types.UntypedString]:
			return &Type{t: types.Typ[types.Int]}
		default:
			panic(x)
		}
	case *types.Map:
		return &Type{t: x.Key()}
	case *types.Slice:
		return &Type{t: types.Typ[types.Int]}
	default:
		panic(x)
	}
}

func (t *Type) Value() *Type {
	switch x := t.t.(type) {
	case *types.Array:
		return &Type{t: x.Elem()}
	case *types.Basic:
		switch x {
		case types.Typ[types.String], types.Typ[types.UntypedString]:
			return &Type{t: types.Typ[types.Byte]}
		default:
			panic(x)
		}
	case *types.Chan:
		return &Type{t: x.Elem()}
	case *types.Map:
		return &Type{t: x.Elem()}
	case *types.Pointer:
		return &Type{t: x.Elem()}
	case *types.Slice:
		return &Type{t: x.Elem()}
	default:
		panic(x)
	}
}

func (t *Type) IsChannel() bool {
	_, ok := t.t.(*types.Chan)
	return ok
}

func (t *Type) HasReceive() bool {
	d := t.t.(*types.Chan).Dir()
	return d == types.SendRecv || d == types.RecvOnly
}

func (t *Type) HasSend() bool {
	d := t.t.(*types.Chan).Dir()
	return d == types.SendRecv || d == types.SendOnly
}

func (t *Type) IsInterface() bool {
	_, ok := t.t.(*types.Interface)
	return ok
}

func (t *Type) Methods() []*Var {
	switch x := t.t.(type) {
	case *types.Interface:
		n := x.NumMethods()
		vs := make([]*Var, n)
		for i := 0; i < n; i++ {
			vs[i] = &Var{Type: &Type{t: x.Method(i).Type()}}
		}
		return vs
	case *types.Named:
		n := x.NumMethods()
		vs := make([]*Var, n)
		for i := 0; i < n; i++ {
			vs[i] = &Var{Type: &Type{t: x.Method(i).Type()}}
		}
		return vs
	}
	panic(t.t)
}

func (t *Type) IsMap() bool {
	_, ok := t.t.(*types.Map)
	return ok
}

func (t *Type) IsNamed() bool {
	_, ok := t.t.(*types.Named)
	return ok
}

func (t *Type) IsAlias() bool {
	return t.t.(*types.Named).Obj().IsAlias()
}

func (t *Type) Name() Name {
	return Name(t.t.(*types.Named).Obj().Name())
}

func (t *Type) Origin() *Type {
	return &Type{t: t.t.(*types.Named).Origin()}
}

func (t *Type) TypeArgs() []*Type {
	args := t.t.(*types.Named).TypeArgs()
	n := args.Len()
	ts := make([]*Type, n)
	for i := 0; i < n; i++ {
		ts[i] = &Type{t: args.At(i)}
	}
	return ts
}

func (t *Type) TypeParams() []*Type {
	switch x := t.t.(type) {
	case *types.Named:
		params := x.TypeParams()
		n := params.Len()
		ts := make([]*Type, n)
		for i := 0; i < n; i++ {
			ts[i] = &Type{t: params.At(i)}
		}
		return ts
	case *types.Signature:
		params := t.t.(*types.Signature).TypeParams()
		n := params.Len()
		ts := make([]*Type, n)
		for i := 0; i < n; i++ {
			ts[i] = &Type{t: params.At(i)}
		}
		return ts
	}
	panic(t.t)
}

func (t *Type) Constraint() *Type {
	return &Type{t.t.(*types.TypeParam).Constraint()}
}

func (t *Type) PackageName() string {
	return t.t.(*types.Named).Obj().Pkg().Name()
}

func (t *Type) PackagePath() string {
	return t.t.(*types.Named).Obj().Pkg().Path()
}

func (t *Type) Underlying() *Type {
	a, b := t.t, t.t.Underlying()
	for a != b {
		a = b
		b = b.Underlying()
	}
	return &Type{t: b}
}

func (t *Type) IsPointer() bool {
	_, ok := t.t.(*types.Pointer)
	return ok
}

func (t *Type) IsSlice() bool {
	_, ok := t.t.(*types.Slice)
	return ok
}

func (t *Type) IsStruct() bool {
	_, ok := t.t.(*types.Struct)
	return ok
}

func convertVar(v *types.Var) *Var {
	return &Var{Name: Name(v.Name()), Type: &Type{t: v.Type()}}
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

func (t *Type) Fields() []*Var {
	return convertVars(structVars(t.t.(*types.Struct)))
}

func (t *Type) IsFunc() bool {
	_, ok := t.t.(*types.Signature)
	return ok
}

func (t *Type) Receiver() *Var {
	return convertVar(t.t.(*types.Signature).Recv())
}

func (t *Type) ReceiverTypeParams() []*Type {
	params := t.t.(*types.Signature).RecvTypeParams()
	n := params.Len()
	ts := make([]*Type, n)
	for i := 0; i < n; i++ {
		ts[i] = &Type{t: params.At(i)}
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

func (t *Type) Params() []*Var {
	return convertVars(tupleVars(t.t.(*types.Signature).Params()))
}

func (t *Type) Results() []*Var {
	return convertVars(tupleVars(t.t.(*types.Signature).Results()))
}

/*
TODO:
methods that fill in empty names in receiver, params, results, inputs (receiver + params)
methods on func itself (or entire type to take into account best receiver name that spans all methods)
import decls
	stores qualifiers for packages, used when referencing an imported type
	ability to add other imports and have the qualifiers updated/sorted out
methods to organize var groups, sort, regroup, etc.
*/

func (t *Type) IsVariadic() bool {
	return t.t.(*types.Signature).Variadic()
}

func (t *Type) String() string {
	// TODO: Handle qualified identifiers.
	return types.TypeString(t.t, nil)
}
