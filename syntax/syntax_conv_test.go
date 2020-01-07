package syntax

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestContext(t *testing.T) {
	// TODO
	// a := []Context{&Comment{Text: "/*a*/"}}
	// b := []Context{&Comment{Text: "/*b*/"}}
}

func TestExpressions(t *testing.T) {
	t.Parallel()
	for exp, str := range expressionString {
		func(exp Expression, str string) {
			t.Run(fmt.Sprintf("%#v", exp), func(t *testing.T) {
				t.Parallel()
				if s, err := ToString(&File{
					Package: &Name{Text: "p"},
					Decls: []Declaration{
						&Var{
							Names:  []*Name{{Text: "_"}},
							Values: []Expression{exp},
						},
					},
				}); err != nil {
					t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(exp), err)
				} else if e := fmt.Sprintf("package p\n\nvar _ = %s\n", str); s != e {
					t.Errorf("Syntax: %s\nString: %s", pretty.Sprint(exp), s)
				}

				// Context
				// a := []Context{&Comment{Text: "/*a*/"}}
				// b := []Context{&Comment{Text: "/*b*/"}}
				if s, err := ToString(&File{
					Package: &Name{Text: "p"},
					Decls: []Declaration{
						&Var{
							Names:  []*Name{{Text: "_"}},
							Values: []Expression{exp},
						},
					},
				}); err != nil {
					t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(exp), err)
				} else if e := fmt.Sprintf("package p\n\nvar _ = %s\n", str); s != e {
					t.Errorf("Syntax: %s\nString: %s", pretty.Sprint(exp), s)
				}
			})
		}(exp, str)
	}
}

var expressionString = map[Expression]string{
	&Name{Text: "z"}:     "z",
	&String{Text: `"z"`}: `"z"`,
	&String{Text: "`z`"}: "`z`",
}

/*
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

func TestSyntax(t *testing.T) {
	s := &File{
		Package: &Name{
			Text: "main",
		},
		Decls: []Declaration{
			&Func{
				Name: &Name{Text: "f"},
				Body: &Block{
					List: []Statement{},
				},
			},
		},
	}
	t.Log(mustToString(s))
	t.FailNow()
}

func TestXAssignSyntax(t *testing.T) {
	s := &File{
		Package: &Name{
			Text: "main",
		},
		Decls: []Syntax{
			&Func{
				Name: &Name{Text: "f"},
				Body: &Block{
					List: []Syntax{
						&Assign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&AddAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&SubtractAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&MultiplyAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&DivideAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&ModuloAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&XorAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&BitAndAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&BitOrAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&ShiftLeftAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&ShiftRightAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
						&AndNotAssign{Left: []Expression{&Name{Text: "x"}}, Right: []Expression{&Name{Text: "y"}}},
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
	t.FailNow()
}

func TestBinarySyntax(t *testing.T) {
	s := &File{
		Package: &Name{
			Text: "main",
		},
		Decls: []Syntax{
			&Func{
				Name: &Name{Text: "f"},
				Body: &Block{
					List: []Syntax{
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Add{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Subtract{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Multiply{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Divide{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Modulo{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&BitAnd{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&BitOr{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&And{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Or{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Xor{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&ShiftLeft{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&ShiftRight{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&AndNot{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Send{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Equal{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&NotEqual{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Less{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&LessEqual{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&Greater{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
						&Assign{Left: []Expression{&Name{Text: "_"}}, Right: []Expression{&GreaterEqual{X: &Name{Text: "x"}, Y: &Name{Text: "y"}}}},
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
	t.FailNow()
}

func TestAllSyntax(t *testing.T) {
	s := &File{
		Package: &Name{
			Text: "main",
		},
		Decls: []Syntax{
			&Import{
				Path: &String{Text: `"io"`},
			},
			&Import{
				Name: &Name{Text: "strs"},
				Path: &String{Text: `"strings"`},
			},
			&ImportList{
				List: []Syntax{
					&Import{
						Path: &String{Text: `"net"`},
					},
					&Import{
						Name: &Name{Text: "temp"},
						Path: &String{Text: `"text/template"`},
					},
				},
			},
			&Const{
				Names:  []*Name{{Text: "con1"}},
				Values: []Syntax{&Int{Text: "1"}},
			},
			&Const{
				Names:  []*Name{{Text: "con2"}},
				Type:   &Name{Text: "string"},
				Values: []Syntax{&String{Text: `"s1"`}},
			},
			&Const{
				Names:  []*Name{{Text: "con3"}, {Text: "con4"}},
				Type:   &Name{Text: "string"},
				Values: []Syntax{&String{Text: `"s2"`}, &String{Text: `"s3"`}},
			},
			&ConstList{
				List: []Syntax{
					&Const{
						Names:  []*Name{{Text: "con1"}},
						Values: []Syntax{&Int{Text: "1"}},
					},
					&Const{
						Names:  []*Name{{Text: "con2"}},
						Type:   &Name{Text: "string"},
						Values: []Syntax{&String{Text: `"s1"`}},
					},
					&Const{
						Names:  []*Name{{Text: "con3"}, {Text: "con4"}},
						Type:   &Name{Text: "string"},
						Values: []Syntax{&String{Text: `"s2"`}, &String{Text: `"s3"`}},
					},
				},
			},
			&Var{
				Names:  []*Name{{Text: "v1"}},
				Values: []Syntax{&Int{Text: "2"}},
			},
			&Var{
				Names:  []*Name{{Text: "v2"}},
				Type:   &Name{Text: "string"},
				Values: []Syntax{&String{Text: `"s4"`}},
			},
			&Var{
				Names:  []*Name{{Text: "v3"}, {Text: "con4"}},
				Type:   &Name{Text: "string"},
				Values: []Syntax{&String{Text: `"s5"`}, &String{Text: `"s6"`}},
			},
			&VarList{
				List: []Syntax{
					&Var{
						Names:  []*Name{{Text: "v1"}},
						Values: []Syntax{&Int{Text: "2"}},
					},
					&Var{
						Names:  []*Name{{Text: "v2"}},
						Type:   &Name{Text: "string"},
						Values: []Syntax{&String{Text: `"s4"`}},
					},
					&Var{
						Names:  []*Name{{Text: "v3"}, {Text: "v4"}},
						Type:   &Name{Text: "string"},
						Values: []Syntax{&String{Text: `"s5"`}, &String{Text: `"s6"`}},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f1"},
			},
			&Func{
				Name: &Name{Text: "f2"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Type: &Name{Text: "int"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f3"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}},
							Type:  &Name{Text: "int"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f4"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f5"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
						&Field{
							Names: []*Name{{Text: "z"}},
							Type:  &Name{Text: "string"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f6"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
						&Field{
							Names: []*Name{{Text: "z"}},
							Type:  &Name{Text: "string"},
						},
					},
				},
				Results: &FieldList{
					List: []*Field{
						&Field{
							Type: &Name{Text: "int"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f7"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
						&Field{
							Names: []*Name{{Text: "z"}},
							Type:  &Name{Text: "string"},
						},
					},
				},
				Results: &FieldList{
					List: []*Field{
						&Field{
							Type: &Name{Text: "int"},
						},
						&Field{
							Type: &Name{Text: "int"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f8"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
						&Field{
							Names: []*Name{{Text: "z"}},
							Type:  &Name{Text: "string"},
						},
					},
				},
				Results: &FieldList{
					List: []*Field{
						&Field{
							Type: &Name{Text: "int"},
						},
						&Field{
							Type: &Name{Text: "int"},
						},
						&Field{
							Type: &Name{Text: "string"},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f9"},
				Parameters: &FieldList{
					List: []*Field{
						&Field{
							Names: []*Name{{Text: "x"}, {Text: "y"}},
							Type:  &Name{Text: "int"},
						},
						&Field{
							Names: []*Name{{Text: "z"}},
							Type:  &Name{Text: "string"},
						},
					},
				},
				Results: &FieldList{
					List: []*Field{
						&Field{
							Type: &Name{Text: "int"},
						},
						&Field{
							Type: &Name{Text: "int"},
						},
						&Field{
							Type: &Name{Text: "string"},
						},
					},
				},
				Body: &Block{
					List: []Syntax{
						&Return{
							Results: []Syntax{
								&Name{Text: "x"},
								&Name{Text: "y"},
								&Name{Text: "z"},
							},
						},
					},
				},
			},
			&Func{
				Name: &Name{Text: "f10"},
				Body: &Block{
					List: []Syntax{
						&Assign{
							Left:  []Syntax{&Name{Text: "_"}},
							Right: []Syntax{&Name{Text: "a"}},
						},
					},
				},
			},
			&Type{
				Name: &Name{Text: "t1"},
				Type: &Name{Text: "int"},
			},
			&Type{
				Name: &Name{Text: "t2"},
				Type: &Struct{
					Fields: &FieldList{
						List: []*Field{
							&Field{
								Type: &Pointer{
									X: &Name{Text: "foo"},
								},
							},
						},
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
*/
