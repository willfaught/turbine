package turbine

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/willfaught/forklift"
)

func errFind(importpath, ident string) error {
	return fmt.Errorf("cannot find declaration: %v.%v", importpath, ident)
}

func findSpec(d *ast.GenDecl, ident string) ast.Spec {

	return nil
}

func searchDecl(fs []*ast.File, name string) ast.Decl {
	var match ast.Decl
	for _, f := range fs {
		ast.Inspect(f, func(node ast.Node) bool {
			if node == nil || match != nil {
				return false
			}
			switch node := node.(type) {
			case *ast.FuncDecl:
				if node.Name.Name == name {
					match = node
					return false
				}
			case *ast.GenDecl:
				for _, spec := range node.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.Name == name {
							match = node
							return false
						}
					case *ast.ValueSpec:
						for _, ident := range spec.Names {
							if ident.Name == name {
								match = node
								return false
							}
						}
					}
				}
				if match != nil {
					match = node
					return false
				}
			}
			return true
		})
		if match != nil {
			break
		}
	}
	return match
}

func searchObject(s *types.Scope, ident string) types.Object {
	if o := s.Lookup(ident); o != nil {
		return o
	}
	for i, n := 0, s.NumChildren(); i < n; i++ {
		if o := searchObject(s.Child(i), ident); o != nil {
			return o
		}
	}
	return nil
}

// Decl is the documentation, identifier, syntax, type, and value for a package declaration.
type Decl struct {
	Doc   []string
	Ident Name
	Type  *Typ
	kind  declKind
	// TODO: Value
}

type declKind int

const (
	declKindConst declKind = iota
	declKindFunc
	declKindType
	declKindVar
)

func (d Decl) IsConst() bool {
	return d.kind == declKindConst
}

func (d Decl) IsFunc() bool {
	return d.kind == declKindFunc
}

func (d Decl) IsType() bool {
	return d.kind == declKindType
}

func (d Decl) IsVar() bool {
	return d.kind == declKindVar
}

func makeDoc(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}
	ss := make([]string, 0, len(cg.List))
	for _, c := range cg.List {
		ss = append(ss, c.Text)
	}
	return ss
}

// LoadDecl returns a Decl for the declaration named name in package path.
func LoadDecl(path, name string) (*Decl, error) {
	p, err := forklift.LoadPackage(path)
	if err != nil {
		return nil, err
	}
	ad := searchDecl(p.Syntax, name)
	if ad == nil {
		return nil, errFind(path, name)
	}
	o := searchObject(p.Types.Scope(), name)
	if o == nil {
		return nil, errFind(path, name)
	}
	var d Decl
	switch ad := ad.(type) {
	case *ast.FuncDecl:
		d.Doc = makeDoc(ad.Doc)
		d.Ident = Name(ad.Name.Name)
		d.Type = &Typ{t: o.Type()}
		d.kind = declKindFunc
	case *ast.GenDecl:
		s := findSpec(ad, name)
		switch s := s.(type) {
		case *ast.TypeSpec:
			doc := s.Doc
			if doc == nil {
				doc = ad.Doc
			}
			d.Doc = makeDoc(doc)
			d.Ident = Name(s.Name.Name)
			d.Type = &Typ{t: o.Type()}
			d.kind = declKindType
		case *ast.ValueSpec:
			match := -1
			for i, n := range s.Names {
				if n.Name == name {
					match = i
					break
				}
			}
			d.Doc = makeDoc(ad.Doc)
			d.Ident = Name(s.Names[match].Name)
			d.Type = &Typ{t: o.Type()}
			switch ad.Tok {
			case token.CONST:
				d.kind = declKindConst
			case token.VAR:
				d.kind = declKindVar
			default:
				panic(ad.Tok)
			}
		default:
			panic(s)
		}
	default:
		panic(ad)
	}
	return &d, nil
}

type Package struct {
	All    []*Decl
	Doc    []string
	Consts []*Decl
	Funcs  []*Decl
	Types  []*Decl
	Vars   []*Decl
}

func LoadPackage(path string) (*Package, error) {
	fp, err := forklift.LoadPackage(path)
	if err != nil {
		return nil, err
	}
	var tp Package
	for _, f := range fp.Syntax {
		var match ast.Decl
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			switch n := n.(type) {
			case *ast.FuncDecl:

				if n.Name.Name == ident {
					match = n
					return false
				}
			case *ast.GenDecl:
				if s := findSpec(n, ident); s != nil {
					match = n
					return false
				}
			}
			return false
		})
		if match != nil {
			break
		}
	}
	return ds, nil
}

/*
func (m *Interface) receiver() (string, error) {
	var names = map[string]struct{}{}

	for _, method := range m.InterfaceMethods {
		for _, p := range method.ParamsFlat {
			names[p.Name] = struct{}{}
		}

		for _, r := range method.ResultsFlat {
			names[r.Name] = struct{}{}
		}
	}

	var words = camelcase.Split(m.ReceiverTypeName)
	var initials []string

	for i, word := range words {
		var r, _ = utf8.DecodeRuneInString(word)

		if r == utf8.RuneError {
			return "", fmt.Errorf("type %v is invalid: invalid utf8 string", m.ReceiverTypeName)
		}

		initials = append(initials, strings.ToLower(fmt.Sprintf("%c", r)))
		words[i] = strings.ToLower(word)
	}

	var tries = []string{
		initials[0],
		initials[len(initials)-1],
		strings.Join(initials, ""),
		words[0],
		words[len(words)-1],
	}

	for _, name := range tries {
		if _, ok := names[name]; !ok {
			return name, nil
		}
	}

	var name = initials[0]

	for {
		name += "_"

		if _, ok := names[name]; !ok {
			break
		}
	}

	return name, nil
}*/
