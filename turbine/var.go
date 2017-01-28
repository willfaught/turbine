package turbine

import (
	"go/ast"
	"go/types"
	"sort"
)

func varGroups(vs []*types.Var, fs []*ast.Field) (joined, ordered, original, split []*VarGroup) {
	if len(fs) == 0 {
		for _, v := range vs {
			var g = &VarGroup{Type: newType(v.Type(), nil)}

			joined = append(joined, g)
			ordered = append(ordered, g)
			original = append(original, g)
			split = append(split, g)
		}

		sort.Sort(byType(ordered))

		return
	}

	var varindex int
	var syntaxids = map[string][]*Ident{}
	var syntaxtype = map[string]*Type{}

	for _, f := range fs {
		var v = vs[varindex]
		var t = newType(v.Type(), f.Type)
		var g = &VarGroup{Type: t}
		var s = t.Syntax

		syntaxtype[s] = t

		for _, n := range f.Names {
			var id = newIdent(n.Name)

			g.Idents = append(g.Idents, id)
			split = append(split, &VarGroup{Idents: []*Ident{id}, Type: t})

			varindex++
		}

		original = append(original, g)

		if last := len(joined) - 1; len(joined) == 0 || joined[last].Type.Syntax != s {
			joined = append(joined, g)
		} else {
			joined[last].Idents = append(joined[last].Idents, g.Idents...)
		}

		syntaxids[s] = append(syntaxids[s], g.Idents...)
		syntaxtype[s] = t
	}

	for s, t := range syntaxtype {
		sort.Sort(byName(syntaxids[s]))
		ordered = append(ordered, &VarGroup{Idents: syntaxids[s], Type: t})
	}

	sort.Sort(byType(ordered))

	return
}

// TODO
type Var struct {
	// TODO
	Ident *Ident

	// TODO
	Type *Type
}

// TODO
type VarGroup struct {
	// TODO
	Idents []*Ident

	// TODO
	Type *Type
}

type byName []*Ident

func (b byName) Len() int           { return len(b) }
func (b byName) Less(i, j int) bool { return b[i].Name < b[j].Name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type byType []*VarGroup

func (b byType) Len() int           { return len(b) }
func (b byType) Less(i, j int) bool { return b[i].Type.Syntax < b[j].Type.Syntax }
func (b byType) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
