package syntax

// Declarations

type Declaration interface {
	Syntax
	declaration()
}

func (*Const) declaration()      {}
func (*ConstList) declaration()  {}
func (*Func) declaration()       {}
func (*Import) declaration()     {}
func (*ImportList) declaration() {}
func (*Type) declaration()       {}
func (*TypeList) declaration()   {}
func (*Var) declaration()        {}
func (*VarList) declaration()    {}

// Expressions

type Expression interface {
	Syntax
	expression()
}

func (*Add) expression()          {}
func (*And) expression()          {}
func (*AndNot) expression()       {}
func (*Array) expression()        {}
func (*Assert) expression()       {}
func (*BitAnd) expression()       {}
func (*BitOr) expression()        {}
func (*Call) expression()         {}
func (*Chan) expression()         {}
func (*ChanIn) expression()       {}
func (*ChanOut) expression()      {}
func (*Composite) expression()    {}
func (*Deref) expression()        {}
func (*Divide) expression()       {}
func (*Ellipsis) expression()     {}
func (*Equal) expression()        {}
func (*Float) expression()        {}
func (*Func) expression()         {}
func (*Greater) expression()      {}
func (*GreaterEqual) expression() {}
func (*Imag) expression()         {}
func (*Index) expression()        {}
func (*Int) expression()          {}
func (*Interface) expression()    {}
func (*KeyValue) expression()     {}
func (*Less) expression()         {}
func (*LessEqual) expression()    {}
func (*Map) expression()          {}
func (*Multiply) expression()     {}
func (*Name) expression()         {}
func (*Negate) expression()       {}
func (*Not) expression()          {}
func (*NotEqual) expression()     {}
func (*Or) expression()           {}
func (*Paren) expression()        {}
func (*Pointer) expression()      {}
func (*Receive) expression()      {}
func (*Ref) expression()          {}
func (*Remainder) expression()    {}
func (*Rune) expression()         {}
func (*Selector) expression()     {}
func (*ShiftLeft) expression()    {}
func (*ShiftRight) expression()   {}
func (*Slice) expression()        {}
func (*String) expression()       {}
func (*Struct) expression()       {}
func (*Subtract) expression()     {}
func (*Xor) expression()          {}

// Statements

type Statement interface {
	Syntax
	statement()
}

func (*AddAssign) statement()        {}
func (*AndNotAssign) statement()     {}
func (*Assign) statement()           {}
func (*BitAndAssign) statement()     {}
func (*BitOrAssign) statement()      {}
func (*Block) statement()            {}
func (*Break) statement()            {}
func (*Case) statement()             {}
func (*Continue) statement()         {}
func (*Dec) statement()              {}
func (*Defer) statement()            {}
func (*Define) statement()           {}
func (*DivideAssign) statement()     {}
func (*Fallthrough) statement()      {}
func (*For) statement()              {}
func (*Go) statement()               {}
func (*Goto) statement()             {}
func (*If) statement()               {}
func (*Inc) statement()              {}
func (*Label) statement()            {}
func (*MultiplyAssign) statement()   {}
func (*Range) statement()            {}
func (*RemainderAssign) statement()  {}
func (*Return) statement()           {}
func (*Select) statement()           {}
func (*Send) statement()             {}
func (*ShiftLeftAssign) statement()  {}
func (*ShiftRightAssign) statement() {}
func (*SubtractAssign) statement()   {}
func (*Switch) statement()           {}
func (*XorAssign) statement()        {}

// Syntax

type Syntax interface {
	syntax()
}

func (*Add) syntax()              {}
func (*AddAssign) syntax()        {}
func (*And) syntax()              {}
func (*AndNot) syntax()           {}
func (*AndNotAssign) syntax()     {}
func (*Array) syntax()            {}
func (*Assert) syntax()           {}
func (*Assign) syntax()           {}
func (*BitAnd) syntax()           {}
func (*BitAndAssign) syntax()     {}
func (*BitOr) syntax()            {}
func (*BitOrAssign) syntax()      {}
func (*Block) syntax()            {}
func (*Break) syntax()            {}
func (*Call) syntax()             {}
func (*Case) syntax()             {}
func (*Chan) syntax()             {}
func (*ChanIn) syntax()           {}
func (*ChanOut) syntax()          {}
func (*Composite) syntax()        {}
func (*Const) syntax()            {}
func (*ConstList) syntax()        {}
func (*Continue) syntax()         {}
func (*Dec) syntax()              {}
func (*Defer) syntax()            {}
func (*Define) syntax()           {}
func (*Deref) syntax()            {}
func (*Divide) syntax()           {}
func (*DivideAssign) syntax()     {}
func (*Ellipsis) syntax()         {}
func (*Equal) syntax()            {}
func (*Fallthrough) syntax()      {}
func (*Field) syntax()            {}
func (*FieldList) syntax()        {}
func (*File) syntax()             {}
func (*Float) syntax()            {}
func (*For) syntax()              {}
func (*Func) syntax()             {}
func (*Go) syntax()               {}
func (*Goto) syntax()             {}
func (*Greater) syntax()          {}
func (*GreaterEqual) syntax()     {}
func (*If) syntax()               {}
func (*Imag) syntax()             {}
func (*Import) syntax()           {}
func (*ImportList) syntax()       {}
func (*Inc) syntax()              {}
func (*Index) syntax()            {}
func (*Int) syntax()              {}
func (*Interface) syntax()        {}
func (*KeyValue) syntax()         {}
func (*Label) syntax()            {}
func (*Less) syntax()             {}
func (*LessEqual) syntax()        {}
func (*Map) syntax()              {}
func (*Multiply) syntax()         {}
func (*MultiplyAssign) syntax()   {}
func (*Name) syntax()             {}
func (*Negate) syntax()           {}
func (*Not) syntax()              {}
func (*NotEqual) syntax()         {}
func (*Or) syntax()               {}
func (*Paren) syntax()            {}
func (*Pointer) syntax()          {}
func (*Range) syntax()            {}
func (*Receive) syntax()          {}
func (*Ref) syntax()              {}
func (*Remainder) syntax()        {}
func (*RemainderAssign) syntax()  {}
func (*Return) syntax()           {}
func (*Rune) syntax()             {}
func (*Select) syntax()           {}
func (*Selector) syntax()         {}
func (*Send) syntax()             {}
func (*ShiftLeft) syntax()        {}
func (*ShiftLeftAssign) syntax()  {}
func (*ShiftRight) syntax()       {}
func (*ShiftRightAssign) syntax() {}
func (*Slice) syntax()            {}
func (*String) syntax()           {}
func (*Struct) syntax()           {}
func (*Subtract) syntax()         {}
func (*SubtractAssign) syntax()   {}
func (*Switch) syntax()           {}
func (*Type) syntax()             {}
func (*TypeList) syntax()         {}
func (*Var) syntax()              {}
func (*VarList) syntax()          {}
func (*Xor) syntax()              {}
func (*XorAssign) syntax()        {}

// Context

type Context interface {
	context()
}

func (*Comment) context() {}
func (*Line) context()    {}
func (*Space) context()   {}

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
	Before     []Context
	After      []Context
	Receiver   *FieldList
	Name       *Name
	Parameters *FieldList
	Results    *FieldList
	Body       *Block
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
	Methods *FieldList
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

type Name struct {
	Before []Context
	After  []Context
	Text   string
}

type Paren struct {
	Before []Context
	After  []Context
	X      Expression
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
	Body   *Block
	Init   Statement
	Type   Statement
	Value  Expression
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
	List    []Declaration
}
