package syntax

import (
	"go/token"
)

/*
TODO

add Syntax.Node()

change Comment to string
*/

type Syntax interface {
	syntax()
}

func (*Array) syntax()        {}
func (*Assert) syntax()       {}
func (*Assign) syntax()       {}
func (*Binary) syntax()       {}
func (*Block) syntax()        {}
func (*Break) syntax()        {}
func (*Call) syntax()         {}
func (*Case) syntax()         {}
func (*Chan) syntax()         {}
func (*ChanIn) syntax()       {}
func (*ChanOut) syntax()      {}
func (*Comment) syntax()      {}
func (*CommentGroup) syntax() {}
func (*Composite) syntax()    {}
func (*Const) syntax()        {}
func (*ConstList) syntax()    {}
func (*Continue) syntax()     {}
func (*Dec) syntax()          {}
func (*Defer) syntax()        {}
func (*Ellipsis) syntax()     {}
func (*Empty) syntax()        {}
func (*Fallthrough) syntax()  {}
func (*Field) syntax()        {}
func (*FieldList) syntax()    {}
func (*File) syntax()         {}
func (*Float) syntax()        {}
func (*For) syntax()          {}
func (*Func) syntax()         {}
func (*Go) syntax()           {}
func (*Goto) syntax()         {}
func (*If) syntax()           {}
func (*Imag) syntax()         {}
func (*Import) syntax()       {}
func (*ImportList) syntax()   {}
func (*Inc) syntax()          {}
func (*Index) syntax()        {}
func (*Int) syntax()          {}
func (*Interface) syntax()    {}
func (*KeyValue) syntax()     {}
func (*Label) syntax()        {}
func (*Line) syntax()         {}
func (*Map) syntax()          {}
func (*Markup) syntax()       {}
func (*Name) syntax()         {}
func (*Package) syntax()      {}
func (*Paren) syntax()        {}
func (*Range) syntax()        {}
func (*Return) syntax()       {}
func (*Rune) syntax()         {}
func (*Select) syntax()       {}
func (*Selector) syntax()     {}
func (*Send) syntax()         {}
func (*Slice) syntax()        {}
func (*String) syntax()       {}
func (*Struct) syntax()       {}
func (*Switch) syntax()       {}
func (*Type) syntax()         {}
func (*TypeList) syntax()     {}
func (*Unary) syntax()        {}
func (*Var) syntax()          {}
func (*VarList) syntax()      {}

type Markup struct {
	Before, After []Syntax
}

type Line struct{}

type Name struct {
	Markup
	Text string
}

type Ellipsis struct {
	Markup
	Elt Syntax
}

type Int struct {
	Markup
	Text string
}

type Float struct {
	Markup
	Text string
}

type Imag struct {
	Markup
	Text string
}

type Rune struct {
	Markup
	Text string
}

type String struct {
	Markup
	Text string
}

type Composite struct {
	Markup
	Type Syntax
	Elts []Syntax
}

type Paren struct {
	Markup
	X Syntax
}

type Selector struct {
	Markup
	X   Syntax
	Sel *Name
}

type Index struct {
	Markup
	X     Syntax
	Index Syntax
}

type Slice struct {
	Markup
	X    Syntax
	Low  Syntax
	High Syntax
	Max  Syntax
}

type Assert struct {
	Markup
	X    Syntax
	Type Syntax
}

type Call struct {
	Markup
	Fun      Syntax
	Args     []Syntax
	Ellipsis bool
}

type Unary struct {
	Markup
	Operator token.Token
	X        Syntax
}

type Binary struct {
	Markup
	X        Syntax
	Operator token.Token
	Y        Syntax
}

type KeyValue struct {
	Markup
	Key   Syntax
	Value Syntax
}

type Array struct {
	Markup
	Length  Syntax
	Element Syntax
}

type Struct struct {
	Markup
	Fields *FieldList
}

type Func struct {
	Markup
	Receiver   *FieldList
	Name       *Name
	Parameters *FieldList
	Results    *FieldList
	Body       *Block
}

type Interface struct {
	Markup
	Methods *FieldList
}

type Map struct {
	Markup
	Key   Syntax
	Value Syntax
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

type Empty struct {
	Markup
	// TODO: Implicit  bool
}

type Label struct {
	Markup
	Label *Name
	Stmt  Syntax
}

type Send struct {
	Markup
	Chan  Syntax
	Value Syntax
}

type Inc struct {
	Markup
	X Syntax
}

type Dec struct {
	Markup
	X Syntax
}

type Assign struct {
	Markup
	Left     []Syntax
	Operator token.Token
	Right    []Syntax
}

type Go struct {
	Markup
	Call *Call
}

type Defer struct {
	Markup
	Call *Call
}

type Return struct {
	Markup
	Results []Syntax
}

type Break struct {
	Markup
	Label *Name
}

type Continue struct {
	Markup
	Label *Name
}

type Goto struct {
	Markup
	Label *Name
}

type Fallthrough struct {
	Markup
}

type Block struct {
	Markup
	List []Syntax
}

type If struct {
	Markup
	Init Syntax
	Cond Syntax
	Body *Block
	Else Syntax
}

type Case struct {
	Markup
	Comm Syntax
	List []Syntax
	Body []Syntax
}

type Switch struct {
	Markup
	Body  *Block
	Init  Syntax
	Type  Syntax
	Value Syntax
}

type Select struct {
	Markup
	Body *Block
}

type For struct {
	Markup
	Init Syntax
	Cond Syntax
	Post Syntax
	Body *Block
}

type Range struct {
	Markup
	Assign     bool
	Key, Value Syntax
	X          Syntax
	Body       *Block
}

type Import struct {
	Markup
	Name *Name
	Path *String
}

type Const struct {
	Markup
	Names  []*Name
	Type   Syntax
	Values []Syntax
}

type Var struct {
	Markup
	Names  []*Name
	Type   Syntax
	Values []Syntax
}

type Type struct {
	Markup
	Name   *Name
	Assign token.Pos
	Type   Syntax
}

type VarList struct {
	Markup
	List []Syntax
}

type ConstList struct {
	Markup
	List []Syntax
}

type TypeList struct {
	Markup
	List []Syntax
}

type ImportList struct {
	Markup
	List []Syntax
}

type File struct {
	Markup
	Name  *Name
	Decls []Syntax
}

type Package struct {
	Files map[string]*File
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

type Comment struct {
	Text string
}

type CommentGroup struct {
	List []*Comment
}
