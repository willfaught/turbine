package syntax

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"sort"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/willfaught/inspect"
)

func TestNodeConv(t *testing.T) {
	var standardPaths []string
	for path := range inspect.Standard {
		standardPaths = append(standardPaths, path)
	}
	sort.Strings(standardPaths)
	for _, path := range standardPaths[:1] {
		t.Run(path, func(t *testing.T) {
			var loaderPkg, err = inspect.BuildLoader.Load(path)
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
