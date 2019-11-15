package syntax

import "go/token"

// Syntax is a simpler syntax that converts to nodes.
type Syntax interface{}

type Array struct {
	Markup
	Length  Syntax
	Element Syntax
}

type Assert struct {
	Markup
	X    Syntax
	Type Syntax
}

type Assign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type Define struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type Add struct {
	Markup
	X Syntax
	Y Syntax
}

type Subtract struct {
	Markup
	X Syntax
	Y Syntax
}

type Multiply struct {
	Markup
	X Syntax
	Y Syntax
}

type Divide struct {
	Markup
	X Syntax
	Y Syntax
}

type Modulo struct {
	Markup
	X Syntax
	Y Syntax
}

type BitAnd struct {
	Markup
	X Syntax
	Y Syntax
}

type BitOr struct {
	Markup
	X Syntax
	Y Syntax
}

type And struct {
	Markup
	X Syntax
	Y Syntax
}

type Or struct {
	Markup
	X Syntax
	Y Syntax
}

type Xor struct {
	Markup
	X Syntax
	Y Syntax
}

type ShiftLeft struct {
	Markup
	X Syntax
	Y Syntax
}

type ShiftRight struct {
	Markup
	X Syntax
	Y Syntax
}

type AndNot struct {
	Markup
	X Syntax
	Y Syntax
}

type Send struct {
	Markup
	X Syntax
	Y Syntax
}

type Equal struct {
	Markup
	X Syntax
	Y Syntax
}

type NotEqual struct {
	Markup
	X Syntax
	Y Syntax
}

type Less struct {
	Markup
	X Syntax
	Y Syntax
}

type LessEqual struct {
	Markup
	X Syntax
	Y Syntax
}

type Greater struct {
	Markup
	X Syntax
	Y Syntax
}

type GreaterEqual struct {
	Markup
	X Syntax
	Y Syntax
}

type AddAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type SubtractAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type MultiplyAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type DivideAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type ModuloAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type BitAndAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type BitOrAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type XorAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type ShiftLeftAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type ShiftRightAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type AndNotAssign struct {
	Markup
	Left  []Syntax
	Right []Syntax
}

type Block struct {
	Markup
	List []Syntax
}

type Break struct {
	Markup
	Label *Name
}

type Call struct {
	Markup
	Fun      Syntax
	Args     []Syntax
	Ellipsis bool
}

type Case struct {
	Markup
	Comm Syntax
	List []Syntax
	Body []Syntax
}

type Chan struct {
	Markup
	Value Syntax
}

type ChanIn struct {
	Markup
	Value Syntax
}

type ChanOut struct {
	Markup
	Value Syntax
}

type Comment struct {
	Text string
}

type CommentGroup struct {
	List []*Comment
}

type Composite struct {
	Markup
	Type Syntax
	Elts []Syntax
}

type Const struct {
	Markup
	Names  []*Name
	Type   Syntax
	Values []Syntax
}

type ConstList struct {
	Markup
	Between []Syntax
	List    []Syntax
}

type Continue struct {
	Markup
	Label *Name
}

type Dec struct {
	Markup
	X Syntax
}

type Defer struct {
	Markup
	Call *Call
}

type Ellipsis struct {
	Markup
	Elt Syntax
}

type Empty struct {
	Markup
	// TODO: Implicit  bool
}

type Fallthrough struct {
	Markup
}

type Field struct {
	Markup
	Names []*Name
	Type  Syntax
	Tag   *String
}

type FieldList struct {
	Markup
	List []*Field
}

type File struct {
	Markup
	Package *Name
	Decls   []Syntax
}

type Float struct {
	Markup
	Text string
}

type For struct {
	Markup
	Init Syntax
	Cond Syntax
	Post Syntax
	Body *Block
}

type Func struct {
	Markup
	Receiver   *FieldList
	Name       *Name
	Parameters *FieldList
	Results    *FieldList
	Body       *Block
}

type If struct {
	Markup
	Init Syntax
	Cond Syntax
	Body *Block
	Else Syntax
}

type Go struct {
	Markup
	Call *Call
}

type Goto struct {
	Markup
	Label *Name
}

type Imag struct {
	Markup
	Text string
}

type Import struct {
	Markup
	Name *Name
	Path *String
}

type ImportList struct {
	Markup
	Between []Syntax
	List    []Syntax
}

type Inc struct {
	Markup
	X Syntax
}

type Index struct {
	Markup
	X     Syntax
	Index Syntax
}

type Int struct {
	Markup
	Text string
}

type Interface struct {
	Markup
	Methods *FieldList
}

type KeyValue struct {
	Markup
	Key   Syntax
	Value Syntax
}

type Label struct {
	Markup
	Label *Name
	Stmt  Syntax
}

type Line struct{}

type Map struct {
	Markup
	Key   Syntax
	Value Syntax
}

type Markup struct {
	Before, After []Syntax
}

type Name struct {
	Markup
	Text string
}

type Package struct {
	Files map[string]*File
}

type Paren struct {
	Markup
	X Syntax
}

type Range struct {
	Markup
	Assign     bool
	Key, Value Syntax
	X          Syntax
	Body       *Block
}

type Return struct {
	Markup
	Results []Syntax
}

type Rune struct {
	Markup
	Text string
}

type Select struct {
	Markup
	Body *Block
}

type Selector struct {
	Markup
	X   Syntax
	Sel *Name
}

type Slice struct {
	Markup
	X    Syntax
	Low  Syntax
	High Syntax
	Max  Syntax
}

type Space struct{}

type Spaces struct {
	Count int
}

type String struct {
	Markup
	Text string
}

type Struct struct {
	Markup
	Fields *FieldList
}

type Switch struct {
	Markup
	Body  *Block
	Init  Syntax
	Type  Syntax
	Value Syntax
}

type Type struct {
	Markup
	Name   *Name
	Assign token.Pos
	Type   Syntax
}

type TypeList struct {
	Markup
	Between []Syntax
	List    []Syntax
}

type Pointer struct {
	Markup
	X Syntax
}

type Ref struct {
	Markup
	X Syntax
}

type Deref struct {
	Markup
	X Syntax
}

type Negate struct {
	Markup
	X Syntax
}

type Receive struct {
	Markup
	X Syntax
}

type Not struct {
	Markup
	X Syntax
}

type Var struct {
	Markup
	Names  []*Name
	Type   Syntax
	Values []Syntax
}

type VarList struct {
	Markup
	Between []Syntax
	List    []Syntax
}
