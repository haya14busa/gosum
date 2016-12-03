package gosum

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestAllRelatedPackage(t *testing.T) {
	const src = `package self

import (
	_ "syscall"
	_ "fmt"
	_ "os"
)

`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("src", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatal(err) // type error
	}

	pkgs := AllRelatedPackages(pkg)

	t.Run("includes self", func(t *testing.T) {
		for _, p := range pkgs {
			if p.Name() == "self" {
				return
			}
		}
		t.Error("package self not found")
	})

	t.Run("no duplication", func(t *testing.T) {
		set := make(map[string]bool)
		for _, p := range pkgs {
			set[p.Name()] = true
		}
		if len(set) != len(pkgs) {
			t.Errorf("package list includes duplication")
		}
	})

}
