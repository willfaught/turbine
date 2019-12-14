package syntax

// Syntax is a simpler syntax that converts to nodes.
type Syntax interface {
	syntax()
}

type Expression interface {
	Syntax
	expression()
}

type Declaration interface {
	Syntax
	declaration()
}

type Statement interface {
	Syntax
	statement()
}

func (*Add) declaration()                {}
func (*Add) expression()                 {}
func (*Add) specification()              {}
func (*Add) statement()                  {}
func (*Add) syntax()                     {}
func (*AddAssign) declaration()          {}
func (*AddAssign) expression()           {}
func (*AddAssign) specification()        {}
func (*AddAssign) statement()            {}
func (*AddAssign) syntax()               {}
func (*And) declaration()                {}
func (*And) expression()                 {}
func (*And) specification()              {}
func (*And) statement()                  {}
func (*And) syntax()                     {}
func (*AndNot) declaration()             {}
func (*AndNot) expression()              {}
func (*AndNot) specification()           {}
func (*AndNot) statement()               {}
func (*AndNot) syntax()                  {}
func (*AndNotAssign) declaration()       {}
func (*AndNotAssign) expression()        {}
func (*AndNotAssign) specification()     {}
func (*AndNotAssign) statement()         {}
func (*AndNotAssign) syntax()            {}
func (*Array) declaration()              {}
func (*Array) expression()               {}
func (*Array) specification()            {}
func (*Array) statement()                {}
func (*Array) syntax()                   {}
func (*Assert) declaration()             {}
func (*Assert) expression()              {}
func (*Assert) specification()           {}
func (*Assert) statement()               {}
func (*Assert) syntax()                  {}
func (*Assign) declaration()             {}
func (*Assign) expression()              {}
func (*Assign) specification()           {}
func (*Assign) statement()               {}
func (*Assign) syntax()                  {}
func (*BitAnd) declaration()             {}
func (*BitAnd) expression()              {}
func (*BitAnd) specification()           {}
func (*BitAnd) statement()               {}
func (*BitAnd) syntax()                  {}
func (*BitAndAssign) declaration()       {}
func (*BitAndAssign) expression()        {}
func (*BitAndAssign) specification()     {}
func (*BitAndAssign) statement()         {}
func (*BitAndAssign) syntax()            {}
func (*BitOr) declaration()              {}
func (*BitOr) expression()               {}
func (*BitOr) specification()            {}
func (*BitOr) statement()                {}
func (*BitOr) syntax()                   {}
func (*BitOrAssign) declaration()        {}
func (*BitOrAssign) expression()         {}
func (*BitOrAssign) specification()      {}
func (*BitOrAssign) statement()          {}
func (*BitOrAssign) syntax()             {}
func (*Block) declaration()              {}
func (*Block) expression()               {}
func (*Block) specification()            {}
func (*Block) statement()                {}
func (*Block) syntax()                   {}
func (*Break) declaration()              {}
func (*Break) expression()               {}
func (*Break) specification()            {}
func (*Break) statement()                {}
func (*Break) syntax()                   {}
func (*Call) declaration()               {}
func (*Call) expression()                {}
func (*Call) specification()             {}
func (*Call) statement()                 {}
func (*Call) syntax()                    {}
func (*Case) declaration()               {}
func (*Case) expression()                {}
func (*Case) specification()             {}
func (*Case) statement()                 {}
func (*Case) syntax()                    {}
func (*Chan) declaration()               {}
func (*Chan) expression()                {}
func (*Chan) specification()             {}
func (*Chan) statement()                 {}
func (*Chan) syntax()                    {}
func (*ChanIn) declaration()             {}
func (*ChanIn) expression()              {}
func (*ChanIn) specification()           {}
func (*ChanIn) statement()               {}
func (*ChanIn) syntax()                  {}
func (*ChanOut) declaration()            {}
func (*ChanOut) expression()             {}
func (*ChanOut) specification()          {}
func (*ChanOut) statement()              {}
func (*ChanOut) syntax()                 {}
func (*Comment) declaration()            {}
func (*Comment) expression()             {}
func (*Comment) specification()          {}
func (*Comment) statement()              {}
func (*Comment) syntax()                 {}
func (*CommentGroup) declaration()       {}
func (*CommentGroup) expression()        {}
func (*CommentGroup) specification()     {}
func (*CommentGroup) statement()         {}
func (*CommentGroup) syntax()            {}
func (*Composite) declaration()          {}
func (*Composite) expression()           {}
func (*Composite) specification()        {}
func (*Composite) statement()            {}
func (*Composite) syntax()               {}
func (*Const) declaration()              {}
func (*Const) expression()               {}
func (*Const) specification()            {}
func (*Const) statement()                {}
func (*Const) syntax()                   {}
func (*ConstList) declaration()          {}
func (*ConstList) expression()           {}
func (*ConstList) specification()        {}
func (*ConstList) statement()            {}
func (*ConstList) syntax()               {}
func (*Continue) declaration()           {}
func (*Continue) expression()            {}
func (*Continue) specification()         {}
func (*Continue) statement()             {}
func (*Continue) syntax()                {}
func (*Dec) declaration()                {}
func (*Dec) expression()                 {}
func (*Dec) specification()              {}
func (*Dec) statement()                  {}
func (*Dec) syntax()                     {}
func (*Defer) declaration()              {}
func (*Defer) expression()               {}
func (*Defer) specification()            {}
func (*Defer) statement()                {}
func (*Defer) syntax()                   {}
func (*Define) declaration()             {}
func (*Define) expression()              {}
func (*Define) specification()           {}
func (*Define) statement()               {}
func (*Define) syntax()                  {}
func (*Deref) declaration()              {}
func (*Deref) expression()               {}
func (*Deref) specification()            {}
func (*Deref) statement()                {}
func (*Deref) syntax()                   {}
func (*Divide) declaration()             {}
func (*Divide) expression()              {}
func (*Divide) specification()           {}
func (*Divide) statement()               {}
func (*Divide) syntax()                  {}
func (*DivideAssign) declaration()       {}
func (*DivideAssign) expression()        {}
func (*DivideAssign) specification()     {}
func (*DivideAssign) statement()         {}
func (*DivideAssign) syntax()            {}
func (*Ellipsis) declaration()           {}
func (*Ellipsis) expression()            {}
func (*Ellipsis) specification()         {}
func (*Ellipsis) statement()             {}
func (*Ellipsis) syntax()                {}
func (*Empty) declaration()              {}
func (*Empty) expression()               {}
func (*Empty) specification()            {}
func (*Empty) statement()                {}
func (*Empty) syntax()                   {}
func (*Equal) declaration()              {}
func (*Equal) expression()               {}
func (*Equal) specification()            {}
func (*Equal) statement()                {}
func (*Equal) syntax()                   {}
func (*Fallthrough) declaration()        {}
func (*Fallthrough) expression()         {}
func (*Fallthrough) specification()      {}
func (*Fallthrough) statement()          {}
func (*Fallthrough) syntax()             {}
func (*Field) declaration()              {}
func (*Field) expression()               {}
func (*Field) specification()            {}
func (*Field) statement()                {}
func (*Field) syntax()                   {}
func (*FieldList) declaration()          {}
func (*FieldList) expression()           {}
func (*FieldList) specification()        {}
func (*FieldList) statement()            {}
func (*FieldList) syntax()               {}
func (*File) declaration()               {}
func (*File) expression()                {}
func (*File) specification()             {}
func (*File) statement()                 {}
func (*File) syntax()                    {}
func (*Float) declaration()              {}
func (*Float) expression()               {}
func (*Float) specification()            {}
func (*Float) statement()                {}
func (*Float) syntax()                   {}
func (*For) declaration()                {}
func (*For) expression()                 {}
func (*For) specification()              {}
func (*For) statement()                  {}
func (*For) syntax()                     {}
func (*Func) declaration()               {}
func (*Func) expression()                {}
func (*Func) specification()             {}
func (*Func) statement()                 {}
func (*Func) syntax()                    {}
func (*Go) declaration()                 {}
func (*Go) expression()                  {}
func (*Go) specification()               {}
func (*Go) statement()                   {}
func (*Go) syntax()                      {}
func (*Goto) declaration()               {}
func (*Goto) expression()                {}
func (*Goto) specification()             {}
func (*Goto) statement()                 {}
func (*Goto) syntax()                    {}
func (*Greater) declaration()            {}
func (*Greater) expression()             {}
func (*Greater) specification()          {}
func (*Greater) statement()              {}
func (*Greater) syntax()                 {}
func (*GreaterEqual) declaration()       {}
func (*GreaterEqual) expression()        {}
func (*GreaterEqual) specification()     {}
func (*GreaterEqual) statement()         {}
func (*GreaterEqual) syntax()            {}
func (*If) declaration()                 {}
func (*If) expression()                  {}
func (*If) specification()               {}
func (*If) statement()                   {}
func (*If) syntax()                      {}
func (*Imag) declaration()               {}
func (*Imag) expression()                {}
func (*Imag) specification()             {}
func (*Imag) statement()                 {}
func (*Imag) syntax()                    {}
func (*Import) declaration()             {}
func (*Import) expression()              {}
func (*Import) specification()           {}
func (*Import) statement()               {}
func (*Import) syntax()                  {}
func (*ImportList) declaration()         {}
func (*ImportList) expression()          {}
func (*ImportList) specification()       {}
func (*ImportList) statement()           {}
func (*ImportList) syntax()              {}
func (*Inc) declaration()                {}
func (*Inc) expression()                 {}
func (*Inc) specification()              {}
func (*Inc) statement()                  {}
func (*Inc) syntax()                     {}
func (*Index) declaration()              {}
func (*Index) expression()               {}
func (*Index) specification()            {}
func (*Index) statement()                {}
func (*Index) syntax()                   {}
func (*Int) declaration()                {}
func (*Int) expression()                 {}
func (*Int) specification()              {}
func (*Int) statement()                  {}
func (*Int) syntax()                     {}
func (*Interface) declaration()          {}
func (*Interface) expression()           {}
func (*Interface) specification()        {}
func (*Interface) statement()            {}
func (*Interface) syntax()               {}
func (*KeyValue) declaration()           {}
func (*KeyValue) expression()            {}
func (*KeyValue) specification()         {}
func (*KeyValue) statement()             {}
func (*KeyValue) syntax()                {}
func (*Label) declaration()              {}
func (*Label) expression()               {}
func (*Label) specification()            {}
func (*Label) statement()                {}
func (*Label) syntax()                   {}
func (*Less) declaration()               {}
func (*Less) expression()                {}
func (*Less) specification()             {}
func (*Less) statement()                 {}
func (*Less) syntax()                    {}
func (*LessEqual) declaration()          {}
func (*LessEqual) expression()           {}
func (*LessEqual) specification()        {}
func (*LessEqual) statement()            {}
func (*LessEqual) syntax()               {}
func (*Line) declaration()               {}
func (*Line) expression()                {}
func (*Line) specification()             {}
func (*Line) statement()                 {}
func (*Line) syntax()                    {}
func (*Map) declaration()                {}
func (*Map) expression()                 {}
func (*Map) specification()              {}
func (*Map) statement()                  {}
func (*Map) syntax()                     {}
func (*Markup) declaration()             {}
func (*Markup) expression()              {}
func (*Markup) specification()           {}
func (*Markup) statement()               {}
func (*Markup) syntax()                  {}
func (*Remainder) declaration()          {}
func (*Remainder) expression()           {}
func (*Remainder) specification()        {}
func (*Remainder) statement()            {}
func (*Remainder) syntax()               {}
func (*RemainderAssign) declaration()    {}
func (*RemainderAssign) expression()     {}
func (*RemainderAssign) specification()  {}
func (*RemainderAssign) statement()      {}
func (*RemainderAssign) syntax()         {}
func (*Multiply) declaration()           {}
func (*Multiply) expression()            {}
func (*Multiply) specification()         {}
func (*Multiply) statement()             {}
func (*Multiply) syntax()                {}
func (*MultiplyAssign) declaration()     {}
func (*MultiplyAssign) expression()      {}
func (*MultiplyAssign) specification()   {}
func (*MultiplyAssign) statement()       {}
func (*MultiplyAssign) syntax()          {}
func (*Name) declaration()               {}
func (*Name) expression()                {}
func (*Name) specification()             {}
func (*Name) statement()                 {}
func (*Name) syntax()                    {}
func (*Negate) declaration()             {}
func (*Negate) expression()              {}
func (*Negate) specification()           {}
func (*Negate) statement()               {}
func (*Negate) syntax()                  {}
func (*Not) declaration()                {}
func (*Not) expression()                 {}
func (*Not) specification()              {}
func (*Not) statement()                  {}
func (*Not) syntax()                     {}
func (*NotEqual) declaration()           {}
func (*NotEqual) expression()            {}
func (*NotEqual) specification()         {}
func (*NotEqual) statement()             {}
func (*NotEqual) syntax()                {}
func (*Or) declaration()                 {}
func (*Or) expression()                  {}
func (*Or) specification()               {}
func (*Or) statement()                   {}
func (*Or) syntax()                      {}
func (*Paren) declaration()              {}
func (*Paren) expression()               {}
func (*Paren) specification()            {}
func (*Paren) statement()                {}
func (*Paren) syntax()                   {}
func (*Pointer) declaration()            {}
func (*Pointer) expression()             {}
func (*Pointer) specification()          {}
func (*Pointer) statement()              {}
func (*Pointer) syntax()                 {}
func (*Range) declaration()              {}
func (*Range) expression()               {}
func (*Range) specification()            {}
func (*Range) statement()                {}
func (*Range) syntax()                   {}
func (*Receive) declaration()            {}
func (*Receive) expression()             {}
func (*Receive) specification()          {}
func (*Receive) statement()              {}
func (*Receive) syntax()                 {}
func (*Ref) declaration()                {}
func (*Ref) expression()                 {}
func (*Ref) specification()              {}
func (*Ref) statement()                  {}
func (*Ref) syntax()                     {}
func (*Return) declaration()             {}
func (*Return) expression()              {}
func (*Return) specification()           {}
func (*Return) statement()               {}
func (*Return) syntax()                  {}
func (*Rune) declaration()               {}
func (*Rune) expression()                {}
func (*Rune) specification()             {}
func (*Rune) statement()                 {}
func (*Rune) syntax()                    {}
func (*Select) declaration()             {}
func (*Select) expression()              {}
func (*Select) specification()           {}
func (*Select) statement()               {}
func (*Select) syntax()                  {}
func (*Selector) declaration()           {}
func (*Selector) expression()            {}
func (*Selector) specification()         {}
func (*Selector) statement()             {}
func (*Selector) syntax()                {}
func (*Send) declaration()               {}
func (*Send) expression()                {}
func (*Send) specification()             {}
func (*Send) statement()                 {}
func (*Send) syntax()                    {}
func (*ShiftLeft) declaration()          {}
func (*ShiftLeft) expression()           {}
func (*ShiftLeft) specification()        {}
func (*ShiftLeft) statement()            {}
func (*ShiftLeft) syntax()               {}
func (*ShiftLeftAssign) declaration()    {}
func (*ShiftLeftAssign) expression()     {}
func (*ShiftLeftAssign) specification()  {}
func (*ShiftLeftAssign) statement()      {}
func (*ShiftLeftAssign) syntax()         {}
func (*ShiftRight) declaration()         {}
func (*ShiftRight) expression()          {}
func (*ShiftRight) specification()       {}
func (*ShiftRight) statement()           {}
func (*ShiftRight) syntax()              {}
func (*ShiftRightAssign) declaration()   {}
func (*ShiftRightAssign) expression()    {}
func (*ShiftRightAssign) specification() {}
func (*ShiftRightAssign) statement()     {}
func (*ShiftRightAssign) syntax()        {}
func (*Slice) declaration()              {}
func (*Slice) expression()               {}
func (*Slice) specification()            {}
func (*Slice) statement()                {}
func (*Slice) syntax()                   {}
func (*Space) declaration()              {}
func (*Space) expression()               {}
func (*Space) specification()            {}
func (*Space) statement()                {}
func (*Space) syntax()                   {}
func (*String) declaration()             {}
func (*String) expression()              {}
func (*String) specification()           {}
func (*String) statement()               {}
func (*String) syntax()                  {}
func (*Struct) declaration()             {}
func (*Struct) expression()              {}
func (*Struct) specification()           {}
func (*Struct) statement()               {}
func (*Struct) syntax()                  {}
func (*Subtract) declaration()           {}
func (*Subtract) expression()            {}
func (*Subtract) specification()         {}
func (*Subtract) statement()             {}
func (*Subtract) syntax()                {}
func (*SubtractAssign) declaration()     {}
func (*SubtractAssign) expression()      {}
func (*SubtractAssign) specification()   {}
func (*SubtractAssign) statement()       {}
func (*SubtractAssign) syntax()          {}
func (*Switch) declaration()             {}
func (*Switch) expression()              {}
func (*Switch) specification()           {}
func (*Switch) statement()               {}
func (*Switch) syntax()                  {}
func (*Type) declaration()               {}
func (*Type) expression()                {}
func (*Type) specification()             {}
func (*Type) statement()                 {}
func (*Type) syntax()                    {}
func (*TypeList) declaration()           {}
func (*TypeList) expression()            {}
func (*TypeList) specification()         {}
func (*TypeList) statement()             {}
func (*TypeList) syntax()                {}
func (*Var) declaration()                {}
func (*Var) expression()                 {}
func (*Var) specification()              {}
func (*Var) statement()                  {}
func (*Var) syntax()                     {}
func (*VarList) declaration()            {}
func (*VarList) expression()             {}
func (*VarList) specification()          {}
func (*VarList) statement()              {}
func (*VarList) syntax()                 {}
func (*Xor) declaration()                {}
func (*Xor) expression()                 {}
func (*Xor) specification()              {}
func (*Xor) statement()                  {}
func (*Xor) syntax()                     {}
func (*XorAssign) declaration()          {}
func (*XorAssign) expression()           {}
func (*XorAssign) specification()        {}
func (*XorAssign) statement()            {}
func (*XorAssign) syntax()               {}

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

type Remainder struct {
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

type RemainderAssign struct {
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
	Init  Syntax
	Type  Syntax
	Value Syntax
}

type Type struct {
	Markup
	Assign bool
	Name   *Name
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
