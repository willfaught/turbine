package syntax

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestFun(t *testing.T) {
	syn := &File{
		Package: &Name{Text: "p"},
		Decls: []Declaration{
			&Import{
				After: []Context{
					&Line{},
					&Line{},
				},
				Path: &String{Text: `"fmt"`},
			},
			&Func{
				Name: &Name{Text: "F"},
				Body: &Block{
					List: []Statement{
						&Return{
							Before: []Context{&Line{}},
							After:  []Context{&Line{}},
							Results: []Expression{
								&Int{Text: "123"},
							},
						},
					},
				},
			},
		},
	}
	s, err := ToString(syn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

// TODO: Test Ellipsis markup around Ellipsis.Elem markup in Call

func TestToString_expressions(t *testing.T) {
	t.Parallel()
	u, v, w, x, y, z := &Name{Text: "u"}, &Name{Text: "v"}, &Name{Text: "w"}, &Name{Text: "x"}, &Name{Text: "y"}, &Name{Text: "z"}
	after := reflect.ValueOf([]Context{&Comment{Text: "/*a*/"}})
	before := reflect.ValueOf([]Context{&Comment{Text: "/*b*/"}})
	for exp, str := range map[Expression]string{
		&Add{X: z, Y: y}:                            "z + y",
		&And{X: z, Y: y}:                            "z && y",
		&AndNot{X: z, Y: y}:                         "z &^ y",
		&Array{Length: &Ellipsis{}, Element: z}:     "[...]z",
		&Array{Element: z}:                          "[]z",
		&Array{Length: &Int{Text: "1"}, Element: z}: "[1]z",
		&Assert{X: z}:                               "z.(type)",
		&Assert{X: z, Type: y}:                      "z.(y)",
		&BitAnd{X: z, Y: y}:                         "z & y",
		&BitOr{X: z, Y: y}:                          "z | y",
		&Call{Fun: z}:                               "z()",
		&Call{
			Fun: z,
			Args: []Expression{
				y,
				x,
			},
		}: "z(y, x)",
		&Call{
			Fun: z,
			Args: []Expression{
				&Name{Text: "y"},
				&Ellipsis{Elem: &Name{Text: "x"}},
			},
		}: "z(y, x...)",
		&Chan{Value: z}:     "chan z",
		&ChanIn{Value: z}:   "<-chan z",
		&ChanOut{Value: z}:  "chan<- z",
		&Composite{Type: z}: "z{}",
		&Composite{Type: z, Elts: []Expression{y, x}}: "z{y, x}",
		&Composite{
			Type: z,
			Elts: []Expression{
				&KeyValue{Key: y, Value: x},
				&KeyValue{Key: w, Value: v},
			},
		}: "z{y: x, w: v}",
		&Deref{X: z}:        "*z",
		&Divide{X: z, Y: y}: "z / y",
		&Ellipsis{}:         "...",
		&Ellipsis{Elem: z}:  "...z",
		&Equal{X: z, Y: y}:  "z == y",
		&Float{Text: "1.0"}: "1.0",
		&Func{}:             "func()",
		&Func{
			Params: &ParamList{
				List: []*Param{
					{Names: []*Name{z, y}, Type: x},
					{Names: []*Name{w, v}, Type: u},
				},
			},
			Results: &ParamList{
				List: []*Param{
					{Names: []*Name{z, y}, Type: x},
					{Names: []*Name{w, v}, Type: u},
				},
			},
			Body: &Block{
				List: []Statement{&Return{}},
			},
		}: "func(z, y x, w, v u) (z, y x, w, v u) { return }",
		&Greater{X: z, Y: y}:      "z > y",
		&GreaterEqual{X: z, Y: y}: "z >= y",
		&Imag{Text: "1i"}:         "1i",
		&Index{X: z, Index: y}:    "z[y]",
		&Int{Text: "1"}:           "1",
		&Interface{}:              "interface{}",
		&Interface{
			Methods: &MethodList{
				List: []*Method{
					{
						Name: z,
					},
				},
			},
		}: "interface{ z() }",
		&Interface{
			Methods: &MethodList{
				List: []*Method{
					{
						Name: &Name{Text: "m1"},
						Params: &ParamList{List: []*Param{
							{
								Names: []*Name{z, y},
								Type:  x,
							},
							{
								Names: []*Name{w, v},
								Type:  u,
							},
						}},
						Results: &ParamList{List: []*Param{
							{
								Names: []*Name{z, y},
								Type:  x,
							},
							{
								Names: []*Name{w, v},
								Type:  u,
							},
						}},
					},
					{
						Name: &Name{Text: "m2"},
						Params: &ParamList{List: []*Param{
							{
								Names: []*Name{z, y},
								Type:  x,
							},
							{
								Names: []*Name{w, v},
								Type:  u,
							},
						}},
						Results: &ParamList{List: []*Param{
							{
								Names: []*Name{z, y},
								Type:  x,
							},
							{
								Names: []*Name{w, v},
								Type:  u,
							},
						}},
					},
				},
			},
		}: "interface {\n\tm1(z, y x, w, v u) (z, y x, w, v u)\n\tm2(z, y x, w, v u) (z, y x, w, v u)\n}",
		&KeyValue{Key: z, Value: y}:           "z: y",
		&Less{X: z, Y: y}:                     "z < y",
		&LessEqual{X: z, Y: y}:                "z <= y",
		&Map{Key: z, Value: y}:                "map[z]y",
		&Multiply{X: z, Y: y}:                 "z * y",
		&Name{Text: "z"}:                      "z",
		&Negate{X: z}:                         "-z",
		&Not{X: z}:                            "!z",
		&NotEqual{X: z, Y: y}:                 "z != y",
		&Or{X: z, Y: y}:                       "z || y",
		&Paren{X: z}:                          "(z)",
		&Pointer{X: z}:                        "*z",
		&Receive{X: z}:                        "<-z",
		&Ref{X: z}:                            "&z",
		&Remainder{X: z, Y: y}:                "z % y",
		&Rune{Text: "'z'"}:                    "'z'",
		&Selector{X: z, Sel: y}:               "z.y",
		&ShiftLeft{X: z, Y: y}:                "z << y",
		&ShiftRight{X: z, Y: y}:               "z >> y",
		&Slice{X: z}:                          "z[:]",
		&Slice{X: z, Low: y}:                  "z[y:]",
		&Slice{X: z, High: y}:                 "z[:y]",
		&Slice{X: z, Low: y, High: x}:         "z[y:x]",
		&Slice{X: z, Low: y, High: x, Max: w}: "z[y:x:w]",
		&String{Text: `"z"`}:                  `"z"`,
		&String{Text: "`z`"}:                  "`z`",
		&Struct{}:                             "struct{}",
		&Struct{
			Fields: &FieldList{
				List: []*Field{
					{
						Type: z,
					},
				},
			},
		}: "struct{ z }",
		&Struct{
			Fields: &FieldList{
				List: []*Field{
					{
						Type: z,
					},
					{
						Names: []*Name{z, y},
						Type:  x,
						Tag:   &String{Text: "`tag1`"},
					},
					{
						Names: []*Name{w, v},
						Type:  u,
						Tag:   &String{Text: "`tag2`"},
					},
				},
			},
		}: "struct {\n\tz\n\tz, y x `tag1`\n\tw, v u `tag2`\n}",
		&Subtract{X: z, Y: y}: "z - y",
		&Xor{X: z, Y: y}:      "z ^ y",
	} {
		func(exp Expression, str string) {
			t.Run(fmt.Sprintf("%#v", exp), func(t *testing.T) {
				t.Parallel()
				file := &File{
					Package: &Name{Text: "p"},
					Decls: []Declaration{
						&Var{
							Names:  []*Name{{Text: "_"}},
							Values: []Expression{exp},
						},
					},
				}
				t.Run("no context", func(t *testing.T) {
					if a, err := ToString(file); err != nil {
						t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(exp), err)
					} else if e := fmt.Sprintf("package p\n\nvar _ = %s\n", str); a != e {
						t.Errorf("Syntax strings do not match\nActual:   %#v\nExpected: %#v\nSyntax: %s", a, e, pretty.Sprint(exp))
					}
				})
				t.Run("context", func(t *testing.T) {
					elem := reflect.ValueOf(exp).Elem()
					elem.FieldByName("Before").Set(before)
					elem.FieldByName("After").Set(after)
					if a, err := ToString(file); err != nil {
						t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(exp), err)
					} else if e := fmt.Sprintf("package p\n\nvar _ = /*b*/ %s /*a*/\n", str); a != e {
						t.Errorf("Syntax strings do not match\nActual:   %#v\nExpected: %#v\nSyntax: %s", a, e, pretty.Sprint(exp))
					}
				})
			})
		}(exp, str)
	}
}

func TestToString_statements(t *testing.T) {
	t.Parallel()
	u, v, w, x, y, z := &Name{Text: "u"}, &Name{Text: "v"}, &Name{Text: "w"}, &Name{Text: "x"}, &Name{Text: "y"}, &Name{Text: "z"}
	after := reflect.ValueOf([]Context{&Comment{Text: "/*a*/"}})
	before := reflect.ValueOf([]Context{&Comment{Text: "/*b*/"}})
	for state, str := range map[Statement]string{
		&AddAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}:    "z, y += x, w",
		&AndNotAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y &^= x, w",
		&Assign{Left: []Expression{z, y}, Right: []Expression{x, w}}:       "z, y = x, w",
		&BitAndAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y &= x, w",
		&BitOrAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}:  "z, y |= x, w",
		&Block{}: "{\n\t}",
		&Block{List: []Statement{&Break{}, &Return{}}}: "{\n\t\tbreak\n\t\treturn\n\t}",
		&Break{}:                       "break",
		&Break{Label: z}:               "break z",
		&Case{}:                        "default:",
		&Case{List: []Expression{z}}:   "case z:",
		&Case{Comm: &Receive{X: z}}:    "case <-z:",
		&Case{Comm: &Send{X: z, Y: y}}: "case z <- y:",
		&Continue{}:                    "continue",
		&Continue{Label: z}:            "continue z",
		&Dec{X: z}:                     "z--",
		&Defer{Call: &Call{Fun: z}}:    "defer z()",
		&Define{Left: []Expression{z, y}, Right: []Expression{x, w}}:       "z, y := x, w",
		&DivideAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y /= x, w",
		&Fallthrough{}: "fallthrough",
		&For{}:         "for {\n\t}",
		&For{Body: &Block{List: []Statement{&Return{}}}}:          "for {\n\t\treturn\n\t}",
		&For{Cond: z, Body: &Block{List: []Statement{&Return{}}}}: "for z {\n\t\treturn\n\t}",
		&For{
			Init: &Define{Left: []Expression{z}, Right: []Expression{y}},
			Cond: x,
			Post: &Inc{X: w},
			Body: &Block{List: []Statement{&Return{}}},
		}: "for z := y; x; w++ {\n\t\treturn\n\t}",
		&Go{Call: &Call{Fun: z}}: "go z()",
		&Goto{Label: z}:          "goto z",
		&If{Cond: z}:             "if z {\n\t}",
		&If{
			Init: &Define{Left: []Expression{z}, Right: []Expression{y}},
			Cond: x,
			Body: &Block{List: []Statement{&Return{}}},
		}: "if z := y; x {\n\t\treturn\n\t}",
		&Inc{X: z}: "z++",
		// TODO: Need separate test since labels are dedented one tab: &Label{Label: z, Stmt: &Inc{X: y}}: "z:\n\ty++",
		&MultiplyAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y *= x, w",
		&Range{
			Key: z,
			X:   y,
		}: "for z := range y {\n\t}",
		&Range{
			Key:   z,
			Value: y,
			X:     x,
			Body:  &Block{List: []Statement{&Return{}}},
		}: "for z, y := range x {\n\t\treturn\n\t}",
		&Receive{X: z}: "<-z",
		&RemainderAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y %= x, w",
		&Return{}:                            "return",
		&Return{Results: []Expression{z, y}}: "return z, y",
		&Select{}:                            "select {}",
		&Select{
			Body: &Block{
				List: []Statement{
					&Case{
						Comm: &Receive{X: z},
						Body: []Statement{
							&Break{},
							&Continue{},
						},
					},
					&Case{
						Comm: &Send{X: y, Y: x},
						Body: []Statement{
							&Break{},
							&Continue{},
						},
					},
					&Case{
						Body: []Statement{
							&Break{},
							&Continue{},
						},
					},
				},
			},
		}: "select {\n\tcase <-z:\n\t\tbreak\n\t\tcontinue\n\tcase y <- x:\n\t\tbreak\n\t\tcontinue\n\tdefault:\n\t\tbreak\n\t\tcontinue\n\t}",
		&Send{X: z, Y: y}: "z <- y",
		&ShiftLeftAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}:  "z, y <<= x, w",
		&ShiftRightAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y >>= x, w",
		&SubtractAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}:   "z, y -= x, w",
		&Switch{}: "switch {\n\t}",
		&Switch{Init: &Define{Left: []Expression{z}, Right: []Expression{y}}}: "switch z := y; {\n\t}",
		&Switch{Value: z}: "switch z {\n\t}",
		&Switch{Init: &Define{Left: []Expression{z}, Right: []Expression{y}}, Value: x}: "switch z := y; x {\n\t}",
		&Switch{Type: &Assert{X: z}}: "switch z.(type) {\n\t}",
		&Switch{Type: &Define{Left: []Expression{z}, Right: []Expression{&Assert{X: y}}}}: "switch z := y.(type) {\n\t}",
		&Switch{
			Init: &Define{Left: []Expression{z}, Right: []Expression{y}},
			Type: &Define{Left: []Expression{x}, Right: []Expression{&Assert{X: w}}},
		}: "switch z := y; x := w.(type) {\n\t}",
		&Switch{
			Init:  &Define{Left: []Expression{z}, Right: []Expression{y}},
			Value: x,
			Body: &Block{
				List: []Statement{
					&Case{
						Comm: &Receive{X: w},
						Body: []Statement{
							&Break{},
						},
					},
					&Case{
						Comm: &Send{X: v, Y: u},
						Body: []Statement{
							&Break{},
						},
					},
					&Case{
						Body: []Statement{
							&Break{},
						},
					},
				},
			},
		}: "switch z := y; x {\n\tcase <-w:\n\t\tbreak\n\tcase v <- u:\n\t\tbreak\n\tdefault:\n\t\tbreak\n\t}",
		&XorAssign{Left: []Expression{z, y}, Right: []Expression{x, w}}: "z, y ^= x, w",
	} {
		func(state Statement, str string) {
			t.Run(fmt.Sprintf("%#v", state), func(t *testing.T) {
				t.Parallel()
				file := &File{
					Package: &Name{Text: "p"},
					Decls: []Declaration{
						&Func{
							Name: &Name{Text: "f"},
							Body: &Block{
								List: []Statement{
									&Switch{After: []Context{&Line{}}},
									state,
								},
							},
						},
					},
				}
				t.Run("no context", func(t *testing.T) {
					if a, err := ToString(file); err != nil {
						t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(state), err)
					} else if e := fmt.Sprintf("package p\n\nfunc f() {\n\tswitch {\n\t}\n\t%s\n}\n", str); a != e {
						t.Errorf("Syntax strings do not match\nActual:   %#v\nExpected: %#v\nSyntax: %s", a, e, pretty.Sprint(state))
					}
				})
				t.Run("context", func(t *testing.T) {
					elem := reflect.ValueOf(state).Elem()
					elem.FieldByName("Before").Set(before)
					elem.FieldByName("After").Set(after)
					if a, err := ToString(file); err != nil {
						t.Errorf("Syntax: %s\nError: %v", pretty.Sprint(state), err)
					} else if e := fmt.Sprintf("package p\n\nfunc f() {\n\tswitch {\n\t}\n\t/*b*/ %s /*a*/\n}\n", str); a != e {
						t.Errorf("Syntax strings do not match\nActual:   %#v\nExpected: %#v\nSyntax: %s", a, e, pretty.Sprint(state))
					}
				})
			})
		}(state, str)
	}
}

/*
func parseFile(content string) *ast.File {
	f, err := parser.ParseFile(token.NewFileSet(), "test.go", content, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}

func TestParse(t *testing.T) {
	f := parseFile(`package p
func f() {
	switch a, b := c, d; y.(type) {
	}
}
`)
	pretty.Println(f)
	t.FailNow()
}


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
