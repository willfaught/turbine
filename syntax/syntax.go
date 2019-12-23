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
func (*Array) expression()        {}
func (*And) expression()          {}
func (*AndNot) expression()       {}
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
func (*Remainder) expression()    {}
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
func (*RemainderAssign) statement()  {}
func (*MultiplyAssign) statement()   {}
func (*Range) statement()            {}
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
func (*Comment) syntax()          {}
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
func (*Line) syntax()             {}
func (*Map) syntax()              {}
func (*Markup) syntax()           {}
func (*Remainder) syntax()        {}
func (*RemainderAssign) syntax()  {}
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
func (*Space) syntax()            {}
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

type Array struct {
	Markup
	Length  Expression
	Element Expression
}

type Assert struct {
	Markup
	X    Expression
	Type Expression
}

type Assign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type Define struct {
	Markup
	Left  []Expression
	Right []Expression
}

type Add struct {
	Markup
	X Expression
	Y Expression
}

type Subtract struct {
	Markup
	X Expression
	Y Expression
}

type Multiply struct {
	Markup
	X Expression
	Y Expression
}

type Divide struct {
	Markup
	X Expression
	Y Expression
}

type Remainder struct {
	Markup
	X Expression
	Y Expression
}

type BitAnd struct {
	Markup
	X Expression
	Y Expression
}

type BitOr struct {
	Markup
	X Expression
	Y Expression
}

type And struct {
	Markup
	X Expression
	Y Expression
}

type Or struct {
	Markup
	X Expression
	Y Expression
}

type Xor struct {
	Markup
	X Expression
	Y Expression
}

type ShiftLeft struct {
	Markup
	X Expression
	Y Expression
}

type ShiftRight struct {
	Markup
	X Expression
	Y Expression
}

type AndNot struct {
	Markup
	X Expression
	Y Expression
}

type Send struct {
	Markup
	X Expression
	Y Expression
}

type Equal struct {
	Markup
	X Expression
	Y Expression
}

type NotEqual struct {
	Markup
	X Expression
	Y Expression
}

type Less struct {
	Markup
	X Expression
	Y Expression
}

type LessEqual struct {
	Markup
	X Expression
	Y Expression
}

type Greater struct {
	Markup
	X Expression
	Y Expression
}

type GreaterEqual struct {
	Markup
	X Expression
	Y Expression
}

type AddAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type SubtractAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type MultiplyAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type DivideAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type RemainderAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type BitAndAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type BitOrAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type XorAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type ShiftLeftAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type ShiftRightAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type AndNotAssign struct {
	Markup
	Left  []Expression
	Right []Expression
}

type Block struct {
	Markup
	List []Statement
}

type Break struct {
	Markup
	Label *Name
}

type Call struct {
	Markup
	Fun      Expression
	Args     []Expression
	Ellipsis bool
}

type Case struct {
	Markup
	Comm Statement
	List []Expression
	Body []Statement
}

type Chan struct {
	Markup
	Value Expression
}

type ChanIn struct {
	Markup
	Value Expression
}

type ChanOut struct {
	Markup
	Value Expression
}

type Comment struct {
	Text string
}

// type CommentGroup struct {
// 	List []*Comment
// }

type Composite struct {
	Markup
	Type Expression
	Elts []Expression
}

type Const struct {
	Markup
	Names  []*Name
	Type   Expression
	Values []Expression
}

type ConstList struct {
	Markup
	Between []Syntax
	List    []Declaration
}

type Continue struct {
	Markup
	Label *Name
}

type Dec struct {
	Markup
	X Expression
}

type Defer struct {
	Markup
	Call *Call
}

// Used for [...]T array type.
// type Ellipsis struct {
// 	Markup
// 	Elt Expression
// }

// type Empty struct {
// 	Markup
// 	// TODO: Implicit  bool
// }

type Fallthrough struct {
	Markup
}

type Field struct {
	Markup
	Names []*Name
	Type  Expression
	Tag   *String
}

type FieldList struct {
	Markup
	List []*Field
}

type File struct {
	Markup
	Package *Name
	Decls   []Declaration
}

type Float struct {
	Markup
	Text string
}

type For struct {
	Markup
	Init Statement
	Cond Expression
	Post Statement
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
	Init Statement
	Cond Expression
	Body *Block
	Else Statement
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
	List    []Declaration
}

type Inc struct {
	Markup
	X Expression
}

type Index struct {
	Markup
	X     Expression
	Index Expression
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
	Key   Expression
	Value Expression
}

type Label struct {
	Markup
	Label *Name
	Stmt  Statement
}

type Line struct{}

type Map struct {
	Markup
	Key   Expression
	Value Expression
}

type Markup struct {
	After  []Syntax
	Before []Syntax
}

type Name struct {
	Markup
	Text string
}

type Paren struct {
	Markup
	X Expression
}

type Range struct {
	Markup
	Assign bool
	Key    Expression
	Value  Expression
	X      Expression
	Body   *Block
}

type Return struct {
	Markup
	Results []Expression
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
	X   Expression
	Sel *Name
}

type Slice struct {
	Markup
	X    Expression
	Low  Expression
	High Expression
	Max  Expression
}

type Space struct {
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
	Init  Statement
	Type  Statement
	Value Expression
}

type Type struct {
	Markup
	Assign bool
	Name   *Name
	Type   Expression
}

type TypeList struct {
	Markup
	Between []Syntax
	List    []Declaration
}

type Pointer struct {
	Markup
	X Expression
}

type Ref struct {
	Markup
	X Expression
}

type Deref struct {
	Markup
	X Expression
}

type Negate struct {
	Markup
	X Expression
}

type Receive struct {
	Markup
	X Expression
}

type Not struct {
	Markup
	X Expression
}

type Var struct {
	Markup
	Names  []*Name
	Type   Expression
	Values []Expression
}

type VarList struct {
	Markup
	Between []Syntax
	List    []Declaration
}
