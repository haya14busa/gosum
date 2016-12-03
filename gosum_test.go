package gosum

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/types/typeutil"
)

func TestNewSumInterface(t *testing.T) {
	const src = `package self

import (
	"go/ast"
)

type A interface {
	isA()
}

type B interface {
	A
	isB()
}

type C struct{}

func (*C) isA() {}

type B1 struct{}

func (*B1) isA() {}
func (*B1) isB() {}

type Dep interface {
	ast.Node
}

type Empty interface {}
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

	pkgscope := typeutil.Dependencies(pkg)

	t.Run("empty scope", func(t *testing.T) {
		for _, name := range pkg.Scope().Names() {
			if obj, ok := pkg.Scope().Lookup(name).(*types.TypeName); ok && types.IsInterface(obj.Type()) {
				got := NewSumInterface(obj.Type().(*types.Named), nil)
				if got == nil {
					continue
				}
				if !(len(got.PkgScope) == 1 && got.PkgScope[0].Name() == "self") {
					t.Errorf("PkgScope should be [src], got %v", got.PkgScope)
				}
				if obj.Name() == "Dep" && got.Implements.Pointers != nil {
					t.Errorf("Dep: Implements.Pointers == %v, want nil", got.Implements.Pointers)
				}
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		for _, name := range pkg.Scope().Names() {
			if obj, ok := pkg.Scope().Lookup(name).(*types.TypeName); ok && types.IsInterface(obj.Type()) {
				got := NewSumInterface(obj.Type().(*types.Named), pkgscope)
				// skip empty interface
				if got == nil {
					continue
				}
				switch obj.Name() {
				case "A":
					if !got.IsInternal {
						t.Error("interface A should be internal interface, but it's public")
					}
					if !(len(got.PkgScope) == 1 && got.PkgScope[0].Name() == "self") {
						t.Errorf("PkgScope should be [src], got %v", got.PkgScope)
					}
					{ // SumInterface.Implements.Interfaces
						want := []string{"B"}
						gotifaces := make([]string, len(got.Implements.Interfaces))
						for i, iface := range got.Implements.Interfaces {
							gotifaces[i] = iface.NamedInterface.Obj().Name()
						}
						if !reflect.DeepEqual(gotifaces, want) {
							t.Errorf("Implements.Interfaces: got %v, want %v", gotifaces, want)
						}
					}
					{ // SumInterface.Implements.Pointers
						want := []string{"*src.B1", "*src.C"}
						gotps := make([]string, len(got.Implements.Pointers))
						for i, ptr := range got.Implements.Pointers {
							gotps[i] = ptr.String()
						}
						if !reflect.DeepEqual(gotps, want) {
							t.Errorf("Implements.Pointers: got %v, want %v", gotps, want)
						}
					}
				case "B":
					if !got.IsInternal {
						t.Error("interface B should be internal interface, but it's public")
					}
					if !(len(got.PkgScope) == 1 && got.PkgScope[0].Name() == "self") {
						t.Errorf("PkgScope should be [src], got %v", got.PkgScope)
					}
					{ // SumInterface.Implements.Interfaces
						want := []string{}
						gotifaces := make([]string, len(got.Implements.Interfaces))
						for i, iface := range got.Implements.Interfaces {
							gotifaces[i] = iface.NamedInterface.Obj().Name()
						}
						if !reflect.DeepEqual(gotifaces, want) {
							t.Errorf("Implements.Interfaces: got %v, want %v", gotifaces, want)
						}
					}
					{ // SumInterface.Implements.Pointers
						want := []string{"*src.B1"}
						gotps := make([]string, len(got.Implements.Pointers))
						for i, ptr := range got.Implements.Pointers {
							gotps[i] = ptr.String()
						}
						if !reflect.DeepEqual(gotps, want) {
							t.Errorf("Implements.Pointers: got %v, want %v", gotps, want)
						}
					}
				case "Dep":
					if got.IsInternal {
						t.Error("interface Dep should be public interface, but it's internal")
					}
					if !reflect.DeepEqual(got.PkgScope, pkgscope) {
						t.Errorf("PkgScope should be %v, got %v", pkgscope, got.PkgScope)
					}
					covered := make(map[string]bool)
					for _, iface := range got.Implements.Interfaces {
						for _, ptr := range iface.Implements.Pointers {
							covered[ptr.String()] = true
						}
					}
					if len(covered) != len(got.Implements.Pointers) {
						msg := "all covered Pointer by each Implements.Interfaces should be same with Implements.Pointers:\n got %v\nwant %v"
						t.Errorf(msg, covered, got.Implements.Pointers)
					}
				default:
					t.Errorf("got unexpected interface: %v", obj)
				}
			}
		}
	})

}

func TestIsInternalInterface(t *testing.T) {
	const src = `package isinternal

type publicInterface interface {
	Public()
}

type publicInterface2 interface {
	publicInterface
}

type internalInterface1 interface {
	internal()
}

type internalInterface2 interface {
	Public()
	internal()
}

type internalInterface3 interface {
	publicInterface
	internal()
}
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

	wants := map[string]bool{
		"publicInterface":    false,
		"publicInterface2":   false,
		"internalInterface1": true,
		"internalInterface2": true,
		"internalInterface3": true,
	}

	for _, name := range pkg.Scope().Names() {
		if obj, ok := pkg.Scope().Lookup(name).(*types.TypeName); ok && types.IsInterface(obj.Type()) {
			typ := obj.Type().Underlying().(*types.Interface)
			got := IsInternalInterface(typ)
			want, ok := wants[name]
			if !ok {
				t.Errorf("IsInternalInterface(%s) == %v, want is not prepared", name, got)
			} else if got != want {
				t.Errorf("IsInternalInterface(%s) == %v, want %v", name, got, want)
			}
		}
	}
}
