package turbine

import (
	"go/ast"
	"go/token"
	"sort"

	"github.com/willfaught/forklift"
	"golang.org/x/tools/go/packages"
)

type Decl struct {
	Name Name
	Type *Type
	kind declKind
}

type declKind int

const (
	declKindInvalid declKind = iota
	declKindConst
	declKindFunc
	declKindType
	declKindVar
)

func (d *Decl) IsConst() bool {
	return d.kind == declKindConst
}

func (d *Decl) IsFunc() bool {
	return d.kind == declKindFunc
}

func (d *Decl) IsType() bool {
	return d.kind == declKindType
}

func (d *Decl) IsVar() bool {
	return d.kind == declKindVar
}

type DeclList []*Decl

func (dl DeclList) Sorted() DeclList {
	dl2 := make(DeclList, len(dl))
	copy(dl2, dl)
	sort.Slice(dl2, func(i, j int) bool {
		return dl2[i].Name < dl2[j].Name
	})
	return dl2
}

type DeclGroup struct {
	Names []Name
	Type  *Type
	kind  declKind
}

func (dg *DeclGroup) IsConst() bool {
	return dg.kind == declKindConst
}

func (dg *DeclGroup) IsType() bool {
	return dg.kind == declKindType
}

func (dg *DeclGroup) IsVar() bool {
	return dg.kind == declKindVar
}

func (dg *DeclGroup) Decls() DeclList {
	dl := make(DeclList, len(dg.Names))
	for i, n := range dg.Names {
		dl[i] = &Decl{Name: n, Type: dg.Type, kind: dg.kind}
	}
	return dl
}

type DeclGroupList []*DeclGroup

func (dgl DeclGroupList) Decls() DeclList {
	var dl DeclList
	for _, dg := range dgl {
		dl = append(dl, dg.Decls()...)
	}
	return dl
}

type DeclBlock struct {
	Groups DeclGroupList
}

func (db *DeclBlock) Decls() DeclList {
	var dl DeclList
	for _, dg := range db.Groups {
		dl = append(dl, dg.Decls()...)
	}
	return dl
}

type DeclBlockList []*DeclBlock

func (dbl DeclBlockList) Decls() DeclList {
	var ds []*Decl
	for _, db := range dbl {
		ds = append(ds, db.Decls()...)
	}
	return ds
}

func (dbl DeclBlockList) DeclGroups() DeclGroupList {
	var dgl DeclGroupList
	for _, db := range dbl {
		dgl = append(dgl, db.Groups...)
	}
	return dgl
}

// LoadPackageDecl returns a [*Decl] for the declaration in the package.
// It returns nil if the declaration does not exist.
func LoadPackageDecl(p *packages.Package, name string) *Decl {
	o := p.Types.Scope().Lookup(name)
	if o == nil {
		return nil
	}
	d := &Decl{Name: Name(name), Type: &Type{t: o.Type()}}
	for _, f := range p.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil || d.kind != declKindInvalid {
				return false
			}
			switch n := n.(type) {
			case *ast.FuncDecl:
				if n.Name.Name == name {
					d.kind = declKindFunc
				}
			case *ast.GenDecl:
				for _, s := range n.Specs {
					switch s := s.(type) {
					case *ast.TypeSpec:
						if s.Name.Name == name {
							d.kind = declKindType
						}
					case *ast.ValueSpec:
						for _, i := range s.Names {
							if i.Name == name {
								switch n.Tok {
								case token.CONST:
									d.kind = declKindConst
								case token.VAR:
									d.kind = declKindVar
								default:
									panic(n.Tok)
								}
								break
							}
						}
					}
				}
			}
			return false
		})
		if d.kind != declKindInvalid {
			break
		}
	}
	if d.kind == declKindInvalid {
		return nil
	}
	return d
}

// LoadDecl returns a [*Decl] for the declaration in the package.
// It returns nil if the declaration does not exist.
// It returns an error if the package cannot be loaded.
func LoadDecl(path, name string) (*Decl, error) {
	p, err := forklift.LoadPackage(path)
	if err != nil {
		return nil, err
	}
	return LoadPackageDecl(p, name), nil
}

type Package struct {
	All    DeclList
	Consts DeclBlockList
	Funcs  DeclList
	Lookup map[string]*Decl
	Types  DeclBlockList
	Vars   DeclBlockList
}

// LoadPackage returns a [*Package] for the package.
// It returns an error if the package cannot be loaded.
func LoadPackage(path string) (*Package, error) {
	fp, err := forklift.LoadPackage(path)
	if err != nil {
		return nil, err
	}
	tp := &Package{Lookup: map[string]*Decl{}}
	for _, f := range fp.Syntax {
		ast.Inspect(f, func(an ast.Node) bool {
			if an == nil {
				return false
			}
			switch an := an.(type) {
			case *ast.FuncDecl:
				td := &Decl{
					Name: Name(an.Name.Name),
					Type: &Type{t: fp.Types.Scope().Lookup(an.Name.Name).Type()},
					kind: declKindFunc,
				}
				tp.All = append(tp.All, td)
				tp.Lookup[string(td.Name)] = td
				tp.Funcs = append(tp.Funcs, td)
			case *ast.GenDecl:
				var db DeclBlock
				for _, s := range an.Specs {
					switch s := s.(type) {
					// TODO: ImportSpec
					case *ast.TypeSpec:
						dn := Name(s.Name.Name)
						dt := &Type{t: fp.Types.Scope().Lookup(s.Name.Name).Type()}
						d := &Decl{Name: dn, Type: dt, kind: declKindType}
						dg := &DeclGroup{Names: []Name{dn}, Type: dt, kind: declKindType}
						tp.All = append(tp.All, d)
						tp.Lookup[string(d.Name)] = d
						db.Groups = append(db.Groups, dg)
					case *ast.ValueSpec:
						var dg DeclGroup
						var dt *Type
						switch an.Tok {
						case token.CONST:
							dg.kind = declKindConst
						case token.VAR:
							dg.kind = declKindVar
						default:
							panic(an.Tok)
						}
						for _, name := range s.Names {
							dn := Name(name.Name)
							if dt == nil {
								dt = &Type{t: fp.Types.Scope().Lookup(name.Name).Type()}
							}
							d := &Decl{Name: dn, Type: dt}
							dg.Names = append(dg.Names, dn)
							switch an.Tok {
							case token.CONST:
								d.kind = declKindConst
							case token.VAR:
								d.kind = declKindVar
							default:
								panic(an.Tok)
							}
							tp.All = append(tp.All, d)
							tp.Lookup[string(d.Name)] = d
						}
						db.Groups = append(db.Groups, &dg)
					default:
						panic(s)
					}
				}
				switch an.Tok {
				case token.CONST:
					tp.Consts = append(tp.Consts, &db)
				case token.TYPE:
					tp.Types = append(tp.Types, &db)
				case token.VAR:
					tp.Vars = append(tp.Vars, &db)
				default:
					panic(an.Tok)
				}
			}
			return false
		})
	}
	return tp, nil
}
