package syntax

import (
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
	p := &Package{
		Files: map[string]*File{
			"main.go": &File{
				Name: &Name{Text: "main"},
				Decls: []Syntax{
					&Var{
						Names:  []*Name{{Text: "foo"}},
						Type:   &Name{Text: "int"},
						Values: []Syntax{&Int{Text: "123"}},
					},
				},
			},
		},
	}
	n := p.Node()
	format.Node(os.Stdout, token.NewFileSet(), n)
}
