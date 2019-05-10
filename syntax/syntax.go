package syntax

import (
	"go/ast"
	"go/token"
)

// Syntax is a simpler syntax that converts to nodes.
type Syntax interface {
	Node() ast.Node
}

func (a *Array) Node() ast.Node        { return convertSyntax(a) }
func (a *Assert) Node() ast.Node       { return convertSyntax(a) }
func (a *Assign) Node() ast.Node       { return convertSyntax(a) }
func (b *Binary) Node() ast.Node       { return convertSyntax(b) }
func (b *Block) Node() ast.Node        { return convertSyntax(b) }
func (b *Break) Node() ast.Node        { return convertSyntax(b) }
func (c *Call) Node() ast.Node         { return convertSyntax(c) }
func (c *Case) Node() ast.Node         { return convertSyntax(c) }
func (c *Chan) Node() ast.Node         { return convertSyntax(c) }
func (c *ChanIn) Node() ast.Node       { return convertSyntax(c) }
func (c *ChanOut) Node() ast.Node      { return convertSyntax(c) }
func (c *Comment) Node() ast.Node      { return convertSyntax(c) }
func (c *CommentGroup) Node() ast.Node { return convertSyntax(c) }
func (c *Composite) Node() ast.Node    { return convertSyntax(c) }
func (c *Const) Node() ast.Node        { return convertSyntax(c) }
func (c *ConstList) Node() ast.Node    { return convertSyntax(c) }
func (c *Continue) Node() ast.Node     { return convertSyntax(c) }
func (d *Dec) Node() ast.Node          { return convertSyntax(d) }
func (d *Defer) Node() ast.Node        { return convertSyntax(d) }
func (e *Ellipsis) Node() ast.Node     { return convertSyntax(e) }
func (e *Empty) Node() ast.Node        { return convertSyntax(e) }
func (f *Fallthrough) Node() ast.Node  { return convertSyntax(f) }
func (f *Field) Node() ast.Node        { return convertSyntax(f) }
func (f *FieldList) Node() ast.Node    { return convertSyntax(f) }
func (f *File) Node() ast.Node         { return convertSyntax(f) }
func (f *Float) Node() ast.Node        { return convertSyntax(f) }
func (f *For) Node() ast.Node          { return convertSyntax(f) }
func (f *Func) Node() ast.Node         { return convertSyntax(f) }
func (g *Go) Node() ast.Node           { return convertSyntax(g) }
func (g *Goto) Node() ast.Node         { return convertSyntax(g) }
func (i *If) Node() ast.Node           { return convertSyntax(i) }
func (i *Imag) Node() ast.Node         { return convertSyntax(i) }
func (i *Import) Node() ast.Node       { return convertSyntax(i) }
func (i *ImportList) Node() ast.Node   { return convertSyntax(i) }
func (i *Inc) Node() ast.Node          { return convertSyntax(i) }
func (i *Index) Node() ast.Node        { return convertSyntax(i) }
func (i *Int) Node() ast.Node          { return convertSyntax(i) }
func (i *Interface) Node() ast.Node    { return convertSyntax(i) }
func (k *KeyValue) Node() ast.Node     { return convertSyntax(k) }
func (l *Label) Node() ast.Node        { return convertSyntax(l) }
func (l *Line) Node() ast.Node         { return convertSyntax(l) }
func (m *Map) Node() ast.Node          { return convertSyntax(m) }
func (m *Markup) Node() ast.Node       { return convertSyntax(m) }
func (n *Name) Node() ast.Node         { return convertSyntax(n) }
func (p *Package) Node() ast.Node      { return convertSyntax(p) }
func (p *Paren) Node() ast.Node        { return convertSyntax(p) }
func (r *Range) Node() ast.Node        { return convertSyntax(r) }
func (r *Return) Node() ast.Node       { return convertSyntax(r) }
func (r *Rune) Node() ast.Node         { return convertSyntax(r) }
func (s *Select) Node() ast.Node       { return convertSyntax(s) }
func (s *Selector) Node() ast.Node     { return convertSyntax(s) }
func (s *Send) Node() ast.Node         { return convertSyntax(s) }
func (s *Slice) Node() ast.Node        { return convertSyntax(s) }
func (s *String) Node() ast.Node       { return convertSyntax(s) }
func (s *Struct) Node() ast.Node       { return convertSyntax(s) }
func (s *Switch) Node() ast.Node       { return convertSyntax(s) }
func (t *Type) Node() ast.Node         { return convertSyntax(t) }
func (t *TypeList) Node() ast.Node     { return convertSyntax(t) }
func (u *Unary) Node() ast.Node        { return convertSyntax(u) }
func (v *Var) Node() ast.Node          { return convertSyntax(v) }
func (v *VarList) Node() ast.Node      { return convertSyntax(v) }

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
	Left     []Syntax
	Operator token.Token
	Right    []Syntax
}

type Binary struct {
	Markup
	X        Syntax
	Operator token.Token
	Y        Syntax
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
	List []Syntax
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
	Name  *Name
	Decls []Syntax
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

type Import struct {
	Markup
	Name *Name
	Path *String
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

type ImportList struct {
	Markup
	List []Syntax
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

type Send struct {
	Markup
	Chan  Syntax
	Value Syntax
}

type Slice struct {
	Markup
	X    Syntax
	Low  Syntax
	High Syntax
	Max  Syntax
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
	List []Syntax
}

type Unary struct {
	Markup
	Operator token.Token
	X        Syntax
}

type Var struct {
	Markup
	Names  []*Name
	Type   Syntax
	Values []Syntax
}

type VarList struct {
	Markup
	List []Syntax
}
