package inspect

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/loader"
)

// Load returns the package for the import path using BuildLoader.
func Load(path string) (*Package, error) {
	return BuildLoader.Load(path)
}

// LoadTest returns the test package for the import path using BuildLoader.
func LoadTest(path string) (*Package, error) {
	return BuildLoader.LoadTest(path)
}

// LoadTestExternal returns the external test package for the import path using
// BuildLoader.
func LoadTestExternal(path string) (*Package, error) {
	return BuildLoader.LoadTestExternal(path)
}

// Loader loads build, node, token, and type information about packages.
type Loader struct {
	BuildContext    build.Context    // The build context.
	BuildImportMode build.ImportMode // The build import mode.
	ParserMode      parser.Mode      // The parser mode.
	TypesConfig     types.Config     // The types configuration.
}

// BuildLoader loads package types from builds.
var BuildLoader = &Loader{
	BuildContext:    build.Default,
	BuildImportMode: build.ImportComment,
	ParserMode:      parser.ParseComments,
	TypesConfig:     types.Config{Importer: importer.Default()},
}

// SourceLoader loads package types from source. It is slower than BuildLoader.
var SourceLoader = &Loader{
	BuildContext:    build.Default,
	BuildImportMode: build.ImportComment,
	ParserMode:      parser.ParseComments,
	TypesConfig:     types.Config{Importer: importer.For("source", nil)},
}

// Load returns the package for the import path.
func (l *Loader) Load(path string) (*Package, error) {
	var bp, err = l.BuildContext.Import(path, "", l.BuildImportMode)
	if err != nil {
		return nil, err
	}
	if len(bp.CgoFiles) == 0 {
		return l.load(bp, [][]string{bp.GoFiles}, bp.Name)
	}
	var c = l.config()
	c.Import(path)
	prog, err := c.Load()
	if err != nil {
		return nil, err
	}
	return l.pkg(bp, prog, prog.Imported[path]), nil
}

// LoadTest returns the test package for the import path.
func (l *Loader) LoadTest(path string) (*Package, error) {
	var bp, err = l.BuildContext.Import(path, "", l.BuildImportMode)
	if err != nil {
		return nil, err
	}
	if len(bp.CgoFiles) == 0 {
		return l.load(bp, [][]string{bp.GoFiles, bp.TestGoFiles}, bp.Name)
	}
	var c = l.config()
	c.ImportWithTests(path)
	prog, err := c.Load()
	if err != nil {
		return nil, err
	}
	return l.pkg(bp, prog, prog.Imported[path]), nil
}

// LoadTestExternal returns the external test package for the import path.
func (l *Loader) LoadTestExternal(path string) (*Package, error) {
	var bp, err = l.BuildContext.Import(path, "", l.BuildImportMode)
	if err != nil {
		return nil, err
	}
	if len(bp.XTestGoFiles) == 0 {
		return nil, nil
	}
	// TODO: This causes errors like
	// "/usr/local/Cellar/go/1.10.3/libexec/src/flag/flag_test.go:28:2:
	// undeclared name: ResetForTesting". The problem seems to be that the
	// external test package is attempting to use an identifier declared in the
	// test code of the imported package being used. Since exported declarations
	// in test code aren't normally importable by other packages, I can
	// understand this being considered an error. However, lots of standard
	// library code seems to expect this to work, and indeed it somehow does
	// when parsed by other code some other way. I can't figure out how to get
	// it to work. golang.org/x/tools/go/loader only seems to parse
	// XTestGoFiles, like I do here.
	//
	// if len(bp.CgoFiles) == 0 {
	// 	return l.load(bp, [][]string{bp.XTestGoFiles}, bp.Name+"_test")
	// }
	var c = l.config()
	c.ImportWithTests(path)
	prog, err := c.Load()
	if err != nil {
		return nil, err
	}
	return l.pkg(bp, prog, prog.Created[0]), nil
}

func (l *Loader) config() *loader.Config {
	return &loader.Config{
		Build:               &l.BuildContext,
		ParserMode:          l.ParserMode,
		TypeChecker:         l.TypesConfig,
		TypeCheckFuncBodies: func(string) bool { return !l.TypesConfig.IgnoreFuncBodies },
	}
}

func (l *Loader) load(bp *build.Package, fss [][]string, pkgName string) (*Package, error) {
	var fset = token.NewFileSet()
	var lookup = map[string]struct{}{}
	for _, fs := range fss {
		for _, f := range fs {
			lookup[f] = struct{}{}
		}
	}
	var naps, err = parser.ParseDir(fset, bp.Dir, func(fi os.FileInfo) bool {
		var _, ok = lookup[fi.Name()]
		return ok
	}, l.ParserMode)
	if err != nil {
		return nil, err
	}
	var ap = naps[pkgName]
	var fs []*ast.File
	for _, f := range ap.Files {
		fs = append(fs, f)
	}
	var info = &types.Info{
		Defs:       map[*ast.Ident]types.Object{},
		Implicits:  map[ast.Node]types.Object{},
		Scopes:     map[ast.Node]*types.Scope{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
		Types:      map[ast.Expr]types.TypeAndValue{},
		Uses:       map[*ast.Ident]types.Object{},
	}
	tp, err := l.TypesConfig.Check(bp.ImportPath, fset, fs, info)
	if err != nil {
		return nil, err
	}
	return &Package{
		Build:     bp,
		Nodes:     ap,
		NodeTypes: info,
		Tokens:    fset,
		Types:     tp,
	}, nil
}

func (l *Loader) pkg(bp *build.Package, prog *loader.Program, pi *loader.PackageInfo) *Package {
	var files map[string]*ast.File
	if len(pi.Files) > 0 {
		files = map[string]*ast.File{}
		for _, f := range pi.Files {
			files[prog.Fset.Position(f.Name.NamePos).Filename] = f
		}
	}
	return &Package{
		Build: bp,
		Nodes: &ast.Package{
			Files: files,
			Name:  pi.Pkg.Name(),
		},
		NodeTypes: &pi.Info,
		Tokens:    prog.Fset,
		Types:     pi.Pkg,
	}
}

// Package is build, node, token, and type information for a package.
type Package struct {
	Build     *build.Package // The build.
	Nodes     *ast.Package   // The nodes.
	NodeTypes *types.Info    // The node types.
	Tokens    *token.FileSet // The files.
	Types     *types.Package // The types.
}
