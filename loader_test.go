package inspect

import (
	"go/ast"
	"go/format"
	"os"
	"sort"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
)

var standardExternal = map[string]struct{}{
	"archive/tar":          struct{}{},
	"archive/zip":          struct{}{},
	"bufio":                struct{}{},
	"bytes":                struct{}{},
	"compress/flate":       struct{}{},
	"compress/gzip":        struct{}{},
	"compress/zlib":        struct{}{},
	"container/heap":       struct{}{},
	"container/list":       struct{}{},
	"container/ring":       struct{}{},
	"context":              struct{}{},
	"crypto/cipher":        struct{}{},
	"crypto/des":           struct{}{},
	"crypto/md5":           struct{}{},
	"crypto/rand":          struct{}{},
	"crypto/sha1":          struct{}{},
	"crypto/sha256":        struct{}{},
	"crypto/tls":           struct{}{},
	"crypto/x509":          struct{}{},
	"database/sql":         struct{}{},
	"debug/dwarf":          struct{}{},
	"encoding/base32":      struct{}{},
	"encoding/base64":      struct{}{},
	"encoding/binary":      struct{}{},
	"encoding/csv":         struct{}{},
	"encoding/gob":         struct{}{},
	"encoding/hex":         struct{}{},
	"encoding/json":        struct{}{},
	"encoding/pem":         struct{}{},
	"encoding/xml":         struct{}{},
	"errors":               struct{}{},
	"flag":                 struct{}{},
	"fmt":                  struct{}{},
	"go/ast":               struct{}{},
	"go/doc":               struct{}{},
	"go/parser":            struct{}{},
	"go/printer":           struct{}{},
	"go/scanner":           struct{}{},
	"go/types":             struct{}{},
	"hash":                 struct{}{},
	"hash/crc32":           struct{}{},
	"html":                 struct{}{},
	"html/template":        struct{}{},
	"image":                struct{}{},
	"image/draw":           struct{}{},
	"image/png":            struct{}{},
	"index/suffixarray":    struct{}{},
	"io":                   struct{}{},
	"io/ioutil":            struct{}{},
	"log":                  struct{}{},
	"log/syslog":           struct{}{},
	"math":                 struct{}{},
	"math/big":             struct{}{},
	"math/bits":            struct{}{},
	"math/cmplx":           struct{}{},
	"math/rand":            struct{}{},
	"mime":                 struct{}{},
	"mime/multipart":       struct{}{},
	"mime/quotedprintable": struct{}{},
	"net":                  struct{}{},
	"net/http":             struct{}{},
	"net/http/cookiejar":   struct{}{},
	"net/http/httptest":    struct{}{},
	"net/http/httptrace":   struct{}{},
	"net/http/httputil":    struct{}{},
	"net/mail":             struct{}{},
	"net/smtp":             struct{}{},
	"net/url":              struct{}{},
	"os":                   struct{}{},
	"os/exec":              struct{}{},
	"os/signal":            struct{}{},
	"path":                 struct{}{},
	"path/filepath":        struct{}{},
	"reflect":              struct{}{},
	"regexp":               struct{}{},
	"runtime":              struct{}{},
	"runtime/debug":        struct{}{},
	"runtime/trace":        struct{}{},
	"sort":                 struct{}{},
	"strconv":              struct{}{},
	"strings":              struct{}{},
	"sync":                 struct{}{},
	"sync/atomic":          struct{}{},
	"syscall":              struct{}{},
	"testing":              struct{}{},
	"text/scanner":         struct{}{},
	"text/tabwriter":       struct{}{},
	"text/template":        struct{}{},
	"time":                 struct{}{},
	"unicode":              struct{}{},
	"unicode/utf16":        struct{}{},
	"unicode/utf8":         struct{}{},
}

func TestLoaders(t *testing.T) {
	var check = func(t *testing.T, path string, p *Package, err error, find bool) {
		t.Helper()
		if find {
			if err != nil {
				t.Errorf("error: actual %v, expected nil", err)
			}
			if p == nil {
				t.Errorf("package: actual nil, expected not nil")
			} else {
				if p.Build == nil {
					t.Errorf("build: actual nil, expected not nil")
				}
				if p.Nodes == nil {
					t.Errorf("nodes: actual nil, expected not nil")
				}
				if p.NodeTypes == nil {
					t.Errorf("node types: actual nil, expected not nil")
				}
				if p.Tokens == nil {
					t.Errorf("tokens: actual nil, expected not nil")
				}
				if p.Types == nil {
					t.Errorf("types: actual nil, expected not nil")
				}
			}
		} else {
			if err != nil {
				t.Errorf("error: actual %v, expected nil", err)
			}
			if p != nil {
				t.Errorf("package: actual %#v, expected nil", p)
			}
		}
	}
	var standardPaths []string
	for path := range Standard {
		standardPaths = append(standardPaths, path)
	}
	sort.Strings(standardPaths)
	for _, test := range []struct {
		name   string
		loader interface {
			Load(string) (*Package, error)
			LoadTest(string) (*Package, error)
			LoadTestExternal(string) (*Package, error)
		}
	}{
		{"BuildLoader", BuildLoader},
		{"SourceLoader", SourceLoader},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, path := range standardPaths {
				var path = path
				t.Run(path, func(t *testing.T) {
					t.Run("Load", func(t *testing.T) {
						var pkg, err = test.loader.Load(path)
						check(t, path, pkg, err, true)
					})
					t.Run("LoadTest", func(t *testing.T) {
						var pkg, err = test.loader.LoadTest(path)
						check(t, path, pkg, err, true)
					})
					t.Run("LoadTestExternal", func(t *testing.T) {
						var pkg, err = test.loader.LoadTestExternal(path)
						var _, find = standardExternal[path]
						check(t, path, pkg, err, find)
					})
				})
			}
		})
	}
}

func TestPrint(t *testing.T) {
	var p, err = BuildLoader.Load("t")
	if err != nil {
		t.Fatal(err)
	}
	// cm := ast.NewCommentMap(p.Tokens, p.Nodes, p.Nodes.Files["/Users/will/Developer/go/src/t/t.go"].Comments)
	pretty.Println(p.Nodes)
	// pretty.Println(cm)
	// p.Nodes.Files["/Users/will/Developer/go/src/t/t.go"].Comments

}

func TestPrint2(t *testing.T) {
	var p, err = BuildLoader.Load("t")
	if err != nil {
		t.Fatal(err)
	}
	var f = p.Nodes.Files["/Users/will/Developer/go/src/t/t.go"]
	f.Comments[0] = &ast.CommentGroup{List: []*ast.Comment{{Slash: f.Comments[0].List[0].Slash, Text: "// T is u."}}}
	f.Decls[0].(*ast.FuncDecl).Doc = &ast.CommentGroup{List: []*ast.Comment{{Slash: f.Decls[0].(*ast.FuncDecl).Doc.List[0].Slash, Text: "// T is v."}}}
	spew.Dump(f.Comments[0])
	spew.Dump(f.Decls[0].(*ast.FuncDecl).Doc)
	// pretty.Println(f)
	if err := format.Node(os.Stdout, p.Tokens, f); err != nil {
		t.Error(err)
	}
}
