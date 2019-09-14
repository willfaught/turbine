package syntax

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"sort"
	"testing"

	"github.com/kr/pretty"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/willfaught/turbine"
)

func loadTestDataFile(name string) *ast.File {
	p, err := turbine.Load("github.com/willfaught/turbine/syntax/testdata/" + name)
	if err != nil {
		panic(err)
	}
	for _, n := range p.Nodes.Files {
		return n
	}
	panic("no files")
}

func parseFile(content string) *ast.File {
	f, err := parser.ParseFile(token.NewFileSet(), "test.go", content, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}

func parseFileLines(content string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", content, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return fset, f
}

func TestConvEmpty(t *testing.T) {
	f := parseFile("package empty")
	a := Convert(nil, f)
	e := &File{Package: &Name{Text: "empty"}}
	pretty.Println(a)
	if !reflect.DeepEqual(a, e) {
		pretty.Ldiff(t, a, e)
	}
	t.FailNow()
}

func TestLines(t *testing.T) {
	fset, _ := parseFileLines(`//1
package p
`)
	f := fset.File(1)
	t.Log(f.LineCount())
	for i := 0; i < f.LineCount(); i++ {
		t.Log(f.LineStart(i + 1))
	}
	t.FailNow()
}

func TestCommentMap(t *testing.T) {
	fset, f := parseFileLines(`package p

var /*1*/ (


	x int
)
`)
	cm := ast.NewCommentMap(fset, f, f.Comments)
	pretty.Println(cm)
	t.FailNow()
}

func TestPrint(t *testing.T) {
	x := &ast.ForStmt{
		Init: &ast.EmptyStmt{Implicit: false},
		Cond: nil,
		Post: &ast.EmptyStmt{Implicit: false},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.BranchStmt{Tok: token.BREAK},
				&ast.EmptyStmt{Implicit: true},
				&ast.BranchStmt{Tok: token.BREAK},
			},
		},
	}
	b := &bytes.Buffer{}
	if err := format.Node(b, token.NewFileSet(), x); err != nil {
		t.Fatal(err)
	}
	t.Log(b.String())
	t.FailNow()
}

func TestParse(t *testing.T) {
	f := parseFile(`package p
func f() {
	for ; ; {
	}
}
`)
	pretty.Println(f)
	t.FailNow()
}

func TestMarkup(t *testing.T) {
	fset, f := parseFileLines(`//1
package p
`)
	a := Convert(fset, f)
	pretty.Println(a)
	t.FailNow()
}

func TestConvertNodeValueExpr(t *testing.T) {
	for _, test := range []struct {
		code   string
		syntax Syntax
	}{
		// Identifier
		{`a`, &Name{Text: "a"}},
		{`/*z*/a/*y*/`, &Name{
			Markup: Markup{
				After:  []Syntax{&CommentGroup{List: []*Comment{{Text: "y"}}}},
				Before: []Syntax{&CommentGroup{List: []*Comment{{Text: "z"}}}},
			},
			Text: "a",
		}},

		// Literals
		{`0.0`, &Float{Text: "0.0"}},
		{`0i`, &Imag{Text: "0i"}},
		{`0`, &Int{Text: "0"}},
		{`'a'`, &Rune{Text: "'a'"}},
		{`"a"`, &String{Text: `"a"`}},

		// Unary
		{`!a`, &Unary{Operator: token.NOT, X: &Name{Text: "a"}}},
		{`*a`, &Unary{Operator: token.MUL, X: &Name{Text: "a"}}},
		{`<-a`, &Unary{Operator: token.ARROW, X: &Name{Text: "a"}}},

		// Binary
		{`1 + 1`, &Binary{Operator: token.ADD, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 - 1`, &Binary{Operator: token.SUB, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 * 1`, &Binary{Operator: token.MUL, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 / 1`, &Binary{Operator: token.QUO, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 % 1`, &Binary{Operator: token.REM, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 & 1`, &Binary{Operator: token.AND, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 | 1`, &Binary{Operator: token.OR, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 ^ 1`, &Binary{Operator: token.XOR, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 << 1`, &Binary{Operator: token.SHL, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 >> 1`, &Binary{Operator: token.SHR, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 &^ 1`, &Binary{Operator: token.AND_NOT, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 && 1`, &Binary{Operator: token.LAND, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 || 1`, &Binary{Operator: token.LOR, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 == 1`, &Binary{Operator: token.EQL, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 < 1`, &Binary{Operator: token.LSS, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 <= 1`, &Binary{Operator: token.LEQ, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 > 1`, &Binary{Operator: token.GTR, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 >= 1`, &Binary{Operator: token.GEQ, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},
		{`1 != 1`, &Binary{Operator: token.NEQ, X: &Int{Text: "1"}, Y: &Int{Text: "1"}}},

		{`(1)`, &Paren{X: &Int{Text: "1"}}},
		{`a.b`, &Selector{X: &Name{Text: "a"}, Sel: &Name{Text: "b"}}},
		{`a[:]`, &Slice{X: &Name{Text: "a"}}},
	} {
		f := parseFile(fmt.Sprintf("package p\nvar _ = %s", test.code))
		a := Convert(nil, f)
		// TODO: Add newlines when supported
		e := &File{
			Package: &Name{Text: "p"},
			Decls: []Syntax{&Var{
				Names:  []*Name{{Text: "_"}},
				Values: []Syntax{test.syntax}},
			},
		}
		if !reflect.DeepEqual(a, e) {
			t.Errorf("Code: %s", test.code)
			pretty.Ldiff(t, a, e)
		}
	}
}

func TestNodeConv(t *testing.T) {
	var standardPaths []string
	for path := range turbine.Standard {
		standardPaths = append(standardPaths, path)
	}
	sort.Strings(standardPaths)
	for _, path := range standardPaths[:1] {
		t.Run(path, func(t *testing.T) {
			var loaderPkg, err = turbine.BuildLoader.Load(path)
			if err != nil {
				t.Fatal(err)
			}
			var nc = &nodeConv{tokens: loaderPkg.Tokens}
			var synPkg = nc.node(loaderPkg.Nodes).(*Package)
			var sc = &syntaxConv{tokenFileSet: token.NewFileSet()}
			var nodePkg = sc.node(synPkg).(*ast.Package)
			if len(loaderPkg.Nodes.Files) != len(nodePkg.Files) {
				t.Fatal(loaderPkg.Nodes.Files, nodePkg.Files)
			}
			for k := range loaderPkg.Nodes.Files {
				if _, ok := nodePkg.Files[k]; !ok {
					t.Fatal(k, loaderPkg.Nodes.Files, nodePkg.Files[k])
				}
			}
			for fileName, loaderFile := range loaderPkg.Nodes.Files {
				nodeFile, ok := nodePkg.Files[fileName]
				if !ok {
					t.Fatal(nodePkg.Files, fileName)
				}
				var loaderBuf = &bytes.Buffer{}
				if err := format.Node(loaderBuf, loaderPkg.Tokens, loaderFile); err != nil {
					t.Fatal(err)
				}
				var nodeBuf = &bytes.Buffer{}
				if err := format.Node(nodeBuf, sc.tokenFileSet, nodeFile); err != nil {
					t.Fatal(err)
				}
				loaderStr := loaderBuf.String()
				nodeStr := nodeBuf.String()
				if nodeStr != loaderStr {
					dmp := diffmatchpatch.New()
					diffs := dmp.DiffMain(loaderStr, nodeStr, false)
					t.Fatal(dmp.DiffPrettyText(diffs))
				}
			}
		})
	}
}
