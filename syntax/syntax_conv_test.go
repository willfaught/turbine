package syntax

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"testing"

	"github.com/willfaught/turbine"
)

func TestEmpty(t *testing.T) {
	p, err := turbine.Load("github.com/willfaught/turbine/syntax/testdata/empty")
	if err != nil {
		t.Fatal(err)
	}
	f := p.Nodes.Files["/Users/Will/Developer/go/src/github.com/willfaught/turbine/syntax/testdata/empty/empty.go"]
	// pretty.Println()
	format.Node(os.Stdout, p.Tokens, f)
	t.FailNow()
}

func Test(t *testing.T) {
	s := &File{
		Package: &Name{
			Markup: Markup{
				After: []Syntax{&Line{}, &Line{}},
			},
			Text: "main",
		},
		Decls: []Syntax{
			&VarList{
				Markup: Markup{
					Before: []Syntax{
						&Comment{Text: "//1"},
					},
				},
				Between: []Syntax{&Comment{Text: "/*123*/"}},
				List: []Syntax{
					&Var{
						Markup: Markup{
							Before: []Syntax{
								&Line{},
								&Comment{Text: "//2"},
							},
							After: []Syntax{
								&Comment{Text: "//3"},
							},
						},
						Names: []*Name{{Text: "x"}},
						Type:  &Name{Text: "int"},
					},
				},
			},
		},
	}
	fset, n := ConvertFile(s)
	// pretty.Println(fset, n)
	if err := format.Node(os.Stdout, fset, n); err != nil {
		t.Error(err)
	}
	// printer.Fprint(os.Stdout, fset, n)
	t.FailNow()
}

func TestBugRepro_TODO_SUBMITBUG(t *testing.T) {
	f := &ast.File{
		Package: 1,
		Name: &ast.Ident{
			NamePos: 8,
			Name:    "main",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				TokPos: 12,
				Tok:    75,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: &ast.Ident{
							NamePos: 18,
							Name:    "name1",
						},
						Path: &ast.BasicLit{ValuePos: 21, Kind: 9, Value: "path1"},
					},
				},
			},
			&ast.GenDecl{
				TokPos: 26,
				Tok:    75,
				Lparen: 32,
				Rparen: 33,
			},
		},
	}
	if err := format.Node(os.Stdout, token.NewFileSet(), f); err != nil {
		t.Error(err)
	}
}
