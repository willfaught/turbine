package turbine

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

func doc(c *ast.CommentGroup) ([]string, []string) {
	var lines = docLines(c)

	return lines, docRaw(lines)
}

func docLines(g *ast.CommentGroup) []string {
	var ss []string

	if g == nil {
		return ss
	}

	for _, c := range g.List {
		ss = append(ss, c.Text)
	}

	return ss
}

func docRaw(ss []string) []string {
	var raw []string

	for _, s := range ss {
		s = strings.TrimLeft(s, "/ ")

		if s != "" {
			raw = append(raw, s)
		}
	}

	return raw
}

func errFind(importpath, ident string) error {
	return fmt.Errorf("cannot find declaration: %v.%v", importpath, ident)
}

func errParse(err error) error {
	return fmt.Errorf("cannot parse package: %v", err)
}

func findDecl(p *ast.Package, ident string) ast.Decl {
	var match ast.Decl

	ast.Inspect(p, func(n ast.Node) bool {
		if n == nil || match != nil {
			return false
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Recv != nil {
				return true
			}

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

		return true
	})

	return match
}

func findPackage(fset *token.FileSet, path string) (*ast.Package, *types.Package, error) {
	var c = build.Default

	c.CgoEnabled = true

	var bp, err = c.Import(path, "", build.ImportComment)

	if err != nil {
		return nil, nil, errParse(err)
	}

	var afs []*ast.File
	var paths = map[string]*ast.File{}

	for _, bf := range append(bp.GoFiles, bp.CgoFiles...) { // TODO: Test files?
		var path = filepath.Join(bp.Dir, bf)
		var af, err = parser.ParseFile(fset, path, nil, parser.ParseComments)

		if err != nil {
			return nil, nil, errParse(err)
		}

		paths[path] = af
		afs = append(afs, af)
	}

	var ap = &ast.Package{Name: afs[0].Name.Name, Files: paths}
	tp, err := (&types.Config{FakeImportC: true, IgnoreFuncBodies: true, Importer: importer.Default()}).Check(bp.ImportPath, fset, afs, nil)

	if err != nil {
		return nil, nil, errParse(err)
	}

	return ap, tp, nil
}

func findSpec(d *ast.GenDecl, ident string) ast.Spec {
	for _, s := range d.Specs {
		switch s := s.(type) {
		case *ast.TypeSpec:
			if s.Name.Name == ident {
				return s
			}

		case *ast.ValueSpec:
			for _, n := range s.Names {
				if n.Name == ident {
					return s
				}
			}
		}
	}

	return nil
}

func findType(p *types.Package, ident string) types.Type {
	return searchScope(p.Scope(), ident)
}

func searchScope(s *types.Scope, ident string) types.Type {
	if o := s.Lookup(ident); o != nil {
		return o.Type()
	}

	for i, n := 0, s.NumChildren(); i < n; i++ {
		if t := searchScope(s.Child(i), ident); t != nil {
			return t
		}
	}

	return nil
}

// Decl is the documentation, identifier, syntax, type, and value for a package declaration.
type Decl struct {
	Doc      []string // The documentation lines.
	DocLines []string // The documentation lines stripped of line comment syntax.
	Ident    *Ident   // The identifier.
	IsConst  bool     // Whether it is a constant.
	IsFunc   bool     // Whether it is a function.
	IsType   bool     // Whether it is a type.
	IsVar    bool     // Whether it is a variable.
	Type     *Type    // The type.
	Value    string   // The value.
}

// NewDecl returns a Decl for the declaration in package importpath named ident.
func NewDecl(importpath, ident string) (*Decl, error) {
	var fset = token.NewFileSet()
	var ap, tp, err = findPackage(fset, importpath)

	if err != nil {
		return nil, err
	}

	var fd = findDecl(ap, ident)

	if fd == nil {
		return nil, errFind(importpath, ident)
	}

	var ft = findType(tp, ident)

	if ft == nil {
		return nil, errFind(importpath, ident)
	}

	var d Decl

	switch fd := fd.(type) {
	case *ast.FuncDecl:
		d.Doc, d.DocLines = doc(fd.Doc)
		d.Ident = newIdent(fd.Name.Name)
		d.IsFunc = true
		d.Type = newType(ft, fd.Type)

	case *ast.GenDecl:
		var s = findSpec(fd, ident)

		switch s := s.(type) {
		case *ast.TypeSpec:
			d.Doc, d.DocLines = doc(fd.Doc)
			d.Ident = newIdent(s.Name.Name)
			d.IsType = true
			d.Type = newType(ft, s.Type)

		case *ast.ValueSpec:
			var match = -1

			for i, n := range s.Names {
				if n.Name == ident {
					match = i

					break
				}
			}

			d.Doc, d.DocLines = doc(fd.Doc)
			d.Ident = newIdent(s.Names[match].Name)
			d.Type = newType(ft, s.Type)

			if len(s.Values) > 0 {
				d.Value = types.ExprString(s.Values[match])
			}

			switch fd.Tok {
			case token.CONST:
				d.IsConst = true

			case token.VAR:
				d.IsVar = true

			default:
				panic(fd.Tok)
			}
		}

	default:
		panic(fd)
	}

	return &d, nil
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
