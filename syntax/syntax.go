package syntax

// Declarations

type Declaration interface {
	Syntax
	Declaration()
}

func (*Const) Declaration()      {}
func (*ConstList) Declaration()  {}
func (*Func) Declaration()       {}
func (*Import) Declaration()     {}
func (*ImportList) Declaration() {}
func (*Type) Declaration()       {}
func (*TypeList) Declaration()   {}
func (*Var) Declaration()        {}
func (*VarList) Declaration()    {}

// Expressions

type Expression interface {
	Syntax
	Expression()
}

func (*Add) Expression()          {}
func (*And) Expression()          {}
func (*AndNot) Expression()       {}
func (*Array) Expression()        {}
func (*Assert) Expression()       {}
func (*BitAnd) Expression()       {}
func (*BitOr) Expression()        {}
func (*Call) Expression()         {}
func (*Chan) Expression()         {}
func (*ChanIn) Expression()       {}
func (*ChanOut) Expression()      {}
func (*Composite) Expression()    {}
func (*Deref) Expression()        {}
func (*Divide) Expression()       {}
func (*Ellipsis) Expression()     {}
func (*Equal) Expression()        {}
func (*Float) Expression()        {}
func (*Func) Expression()         {}
func (*Greater) Expression()      {}
func (*GreaterEqual) Expression() {}
func (*Imag) Expression()         {}
func (*Index) Expression()        {}
func (*Int) Expression()          {}
func (*Interface) Expression()    {}
func (*KeyValue) Expression()     {}
func (*Less) Expression()         {}
func (*LessEqual) Expression()    {}
func (*Map) Expression()          {}
func (*Multiply) Expression()     {}
func (*Name) Expression()         {}
func (*Negate) Expression()       {}
func (*Not) Expression()          {}
func (*NotEqual) Expression()     {}
func (*Or) Expression()           {}
func (*Paren) Expression()        {}
func (*Pointer) Expression()      {}
func (*Receive) Expression()      {}
func (*Ref) Expression()          {}
func (*Remainder) Expression()    {}
func (*Rune) Expression()         {}
func (*Selector) Expression()     {}
func (*ShiftLeft) Expression()    {}
func (*ShiftRight) Expression()   {}
func (*Slice) Expression()        {}
func (*String) Expression()       {}
func (*Struct) Expression()       {}
func (*Subtract) Expression()     {}
func (*Xor) Expression()          {}

// Statements

type Statement interface {
	Syntax
	Statement()
}

// TODO: Add exprs that can compile as statements

func (*AddAssign) Statement()        {}
func (*AndNotAssign) Statement()     {}
func (*Assert) Statement()           {}
func (*Assign) Statement()           {}
func (*BitAndAssign) Statement()     {}
func (*BitOrAssign) Statement()      {}
func (*Block) Statement()            {}
func (*Break) Statement()            {}
func (*Call) Statement()             {}
func (*Case) Statement()             {}
func (*Continue) Statement()         {}
func (*Dec) Statement()              {}
func (*Defer) Statement()            {}
func (*Define) Statement()           {}
func (*DivideAssign) Statement()     {}
func (*Fallthrough) Statement()      {}
func (*For) Statement()              {}
func (*Go) Statement()               {}
func (*Goto) Statement()             {}
func (*If) Statement()               {}
func (*Inc) Statement()              {}
func (*Label) Statement()            {}
func (*MultiplyAssign) Statement()   {}
func (*Range) Statement()            {}
func (*Receive) Statement()          {}
func (*RemainderAssign) Statement()  {}
func (*Return) Statement()           {}
func (*Select) Statement()           {}
func (*Send) Statement()             {}
func (*ShiftLeftAssign) Statement()  {}
func (*ShiftRightAssign) Statement() {}
func (*SubtractAssign) Statement()   {}
func (*Switch) Statement()           {}
func (*XorAssign) Statement()        {}

// Syntax

type Syntax interface {
	Syntax()
}

func (*Add) Syntax()              {}
func (*AddAssign) Syntax()        {}
func (*And) Syntax()              {}
func (*AndNot) Syntax()           {}
func (*AndNotAssign) Syntax()     {}
func (*Array) Syntax()            {}
func (*Assert) Syntax()           {}
func (*Assign) Syntax()           {}
func (*BitAnd) Syntax()           {}
func (*BitAndAssign) Syntax()     {}
func (*BitOr) Syntax()            {}
func (*BitOrAssign) Syntax()      {}
func (*Block) Syntax()            {}
func (*Break) Syntax()            {}
func (*Call) Syntax()             {}
func (*Case) Syntax()             {}
func (*Chan) Syntax()             {}
func (*ChanIn) Syntax()           {}
func (*ChanOut) Syntax()          {}
func (*Composite) Syntax()        {}
func (*Const) Syntax()            {}
func (*ConstList) Syntax()        {}
func (*Continue) Syntax()         {}
func (*Dec) Syntax()              {}
func (*Defer) Syntax()            {}
func (*Define) Syntax()           {}
func (*Deref) Syntax()            {}
func (*Divide) Syntax()           {}
func (*DivideAssign) Syntax()     {}
func (*Ellipsis) Syntax()         {}
func (*Equal) Syntax()            {}
func (*Fallthrough) Syntax()      {}
func (*Field) Syntax()            {}
func (*FieldList) Syntax()        {}
func (*File) Syntax()             {}
func (*Float) Syntax()            {}
func (*For) Syntax()              {}
func (*Func) Syntax()             {}
func (*Go) Syntax()               {}
func (*Goto) Syntax()             {}
func (*Greater) Syntax()          {}
func (*GreaterEqual) Syntax()     {}
func (*If) Syntax()               {}
func (*Imag) Syntax()             {}
func (*Import) Syntax()           {}
func (*ImportList) Syntax()       {}
func (*Inc) Syntax()              {}
func (*Index) Syntax()            {}
func (*Int) Syntax()              {}
func (*Interface) Syntax()        {}
func (*KeyValue) Syntax()         {}
func (*Label) Syntax()            {}
func (*Less) Syntax()             {}
func (*LessEqual) Syntax()        {}
func (*Map) Syntax()              {}
func (*Method) Syntax()           {}
func (*MethodList) Syntax()       {}
func (*Multiply) Syntax()         {}
func (*MultiplyAssign) Syntax()   {}
func (*Name) Syntax()             {}
func (*Negate) Syntax()           {}
func (*Not) Syntax()              {}
func (*NotEqual) Syntax()         {}
func (*Or) Syntax()               {}
func (*Param) Syntax()            {}
func (*ParamList) Syntax()        {}
func (*Paren) Syntax()            {}
func (*Pointer) Syntax()          {}
func (*Range) Syntax()            {}
func (*Receive) Syntax()          {}
func (*Receiver) Syntax()         {}
func (*Ref) Syntax()              {}
func (*Remainder) Syntax()        {}
func (*RemainderAssign) Syntax()  {}
func (*Return) Syntax()           {}
func (*Rune) Syntax()             {}
func (*Select) Syntax()           {}
func (*Selector) Syntax()         {}
func (*Send) Syntax()             {}
func (*ShiftLeft) Syntax()        {}
func (*ShiftLeftAssign) Syntax()  {}
func (*ShiftRight) Syntax()       {}
func (*ShiftRightAssign) Syntax() {}
func (*Slice) Syntax()            {}
func (*String) Syntax()           {}
func (*Struct) Syntax()           {}
func (*Subtract) Syntax()         {}
func (*SubtractAssign) Syntax()   {}
func (*Switch) Syntax()           {}
func (*Type) Syntax()             {}
func (*TypeList) Syntax()         {}
func (*Var) Syntax()              {}
func (*VarList) Syntax()          {}
func (*Xor) Syntax()              {}
func (*XorAssign) Syntax()        {}

// Context

type Context interface {
	Context()
}

func (*Comment) Context() {}
func (*Line) Context()    {}
func (*Space) Context()   {}

type Array struct {
	Before  []Context
	After   []Context
	Length  Expression
	Element Expression
}

type Assert struct {
	Before []Context
	After  []Context
	X      Expression
	Type   Expression
}

type Assign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type Define struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type Add struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Subtract struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Multiply struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Divide struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Remainder struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type BitAnd struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type BitOr struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type And struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Or struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Xor struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type ShiftLeft struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type ShiftRight struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type AndNot struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Send struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Equal struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type NotEqual struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Less struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type LessEqual struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type Greater struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type GreaterEqual struct {
	Before []Context
	After  []Context
	X      Expression
	Y      Expression
}

type AddAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type SubtractAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type MultiplyAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type DivideAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type RemainderAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type BitAndAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type BitOrAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type XorAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type ShiftLeftAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type ShiftRightAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type AndNotAssign struct {
	Before []Context
	After  []Context
	Left   []Expression
	Right  []Expression
}

type Block struct {
	Before []Context
	After  []Context
	List   []Statement
}

type Break struct {
	Before []Context
	After  []Context
	Label  *Name
}

type Call struct {
	Before []Context
	After  []Context
	Fun    Expression
	Args   []Expression
}

type Case struct {
	Before []Context
	After  []Context
	Comm   Statement
	List   []Expression
	Body   []Statement
}

type Chan struct {
	Before []Context
	After  []Context
	Value  Expression
}

type ChanIn struct {
	Before []Context
	After  []Context
	Value  Expression
}

type ChanOut struct {
	Before []Context
	After  []Context
	Value  Expression
}

type Comment struct {
	Text string
}

// TODO
// type CommentGroup struct {
// 	List []*Comment
// }

type Composite struct {
	Before []Context
	After  []Context
	Type   Expression
	Elts   []Expression
}

type Const struct {
	Before []Context
	After  []Context
	Names  []*Name
	Type   Expression
	Values []Expression
}

type ConstList struct {
	Before  []Context
	After   []Context
	Between []Context
	List    []Declaration
}

type Continue struct {
	Before []Context
	After  []Context
	Label  *Name
}

type Dec struct {
	Before []Context
	After  []Context
	X      Expression
}

type Defer struct {
	Before []Context
	After  []Context
	Call   *Call
}

type Ellipsis struct {
	Before []Context
	After  []Context
	Elem   Expression
}

// TODO:
// type Empty struct {
// 	Before []Context
//	After  []Context
// 	// TODO: Implicit  bool
// }

type Fallthrough struct {
	Before []Context
	After  []Context
}

type Field struct {
	Before []Context
	After  []Context
	Names  []*Name
	Type   Expression
	Tag    *String
}

type FieldList struct {
	Before []Context
	After  []Context
	List   []*Field
}

type File struct {
	Before  []Context
	After   []Context
	Package *Name
	Decls   []Declaration
}

type Float struct {
	Before []Context
	After  []Context
	Text   string
}

type For struct {
	Before []Context
	After  []Context
	Init   Statement
	Cond   Expression
	Post   Statement
	Body   *Block
}

type Func struct {
	Before   []Context
	After    []Context
	Receiver *Receiver
	Name     *Name
	Params   *ParamList
	Results  *ParamList
	Body     *Block
}

type If struct {
	Before []Context
	After  []Context
	Init   Statement
	Cond   Expression
	Body   *Block
	Else   Statement
}

type Go struct {
	Before []Context
	After  []Context
	Call   *Call
}

type Goto struct {
	Before []Context
	After  []Context
	Label  *Name
}

type Imag struct {
	Before []Context
	After  []Context
	Text   string
}

type Import struct {
	Before []Context
	After  []Context
	Name   *Name
	Path   *String
}

type ImportList struct {
	Before  []Context
	After   []Context
	Between []Context
	List    []Declaration
}

type Inc struct {
	Before []Context
	After  []Context
	X      Expression
}

type Index struct {
	Before []Context
	After  []Context
	X      Expression
	Index  Expression
}

type Int struct {
	Before []Context
	After  []Context
	Text   string
}

type Interface struct {
	Before  []Context
	After   []Context
	Methods *MethodList
}

type KeyValue struct {
	Before []Context
	After  []Context
	Key    Expression
	Value  Expression
}

type Label struct {
	Before []Context
	After  []Context
	Label  *Name
	Stmt   Statement
}

type Line struct{}

type Map struct {
	Before []Context
	After  []Context
	Key    Expression
	Value  Expression
}

type Method struct {
	Before  []Context
	After   []Context
	Name    *Name
	Params  *ParamList
	Results *ParamList
}

type MethodList struct {
	Before []Context
	After  []Context
	List   []*Method
}

type Name struct {
	Before []Context
	After  []Context
	Text   string
}

type Param struct {
	Before []Context
	After  []Context
	Names  []*Name
	Type   Expression
}

type ParamList struct {
	Before []Context
	After  []Context
	List   []*Param
}

type Paren struct {
	Before []Context
	After  []Context
	X      Expression
}

type Receiver struct {
	Before []Context
	After  []Context
	Name   *Name
	Type   Expression
}

type Range struct {
	Before []Context
	After  []Context
	Assign bool
	Key    Expression
	Value  Expression
	X      Expression
	Body   *Block
}

type Return struct {
	Before  []Context
	After   []Context
	Results []Expression
}

type Rune struct {
	Before []Context
	After  []Context
	Text   string
}

type Select struct {
	Before []Context
	After  []Context
	Body   *Block
}

type Selector struct {
	Before []Context
	After  []Context
	X      Expression
	Sel    *Name
}

type Slice struct {
	Before []Context
	After  []Context
	X      Expression
	Low    Expression
	High   Expression
	Max    Expression
}

type Space struct {
	Count int
}

type String struct {
	Before []Context
	After  []Context
	Text   string
}

type Struct struct {
	Before []Context
	After  []Context
	Fields *FieldList
}

type Switch struct {
	Before []Context
	After  []Context
	Init   Statement
	Value  Expression
	Type   Statement
	Body   *Block
}

type Type struct {
	Before []Context
	After  []Context
	Assign bool
	Name   *Name
	Type   Expression
}

type TypeList struct {
	Before  []Context
	After   []Context
	Between []Context
	List    []Declaration
}

type Pointer struct {
	Before []Context
	After  []Context
	X      Expression
}

type Ref struct {
	Before []Context
	After  []Context
	X      Expression
}

type Deref struct {
	Before []Context
	After  []Context
	X      Expression
}

type Negate struct {
	Before []Context
	After  []Context
	X      Expression
}

type Receive struct {
	Before []Context
	After  []Context
	X      Expression
}

type Not struct {
	Before []Context
	After  []Context
	X      Expression
}

type Var struct {
	Before []Context
	After  []Context
	Names  []*Name
	Type   Expression
	Values []Expression
}

type VarList struct {
	Before  []Context
	After   []Context
	Between []Context
	List    []Declaration // TODO: Change to []*Var
}
