// Package gosum provides utilities for Sum/Union/Variant like types.
package gosum

import "go/types"

// SumInterface represents Sum type by interface{} type.
type SumInterface struct {
	// underling interface type. It uses types.Named instead of types.Interface
	// for convenience.
	NamedInterface *types.Named

	// is "internal interface" which includes unexported methods, otherwise
	// "public" interface.
	IsInternal bool

	// Implements holds types which implement SumInterface type.
	Implements struct {
		Interfaces []*SumInterface
		Pointers   []*types.Pointer
	}

	// "Implements" only holds types in PkgScope. It's especially for "public"
	// interface.
	PkgScope []*types.Package
}

// NewSumInterface returns SumInterface from given interface type. If the given
// interface is "public", this function searches implemented types in pkgscope.
// pkgscope can be nil and it's useful to pass nil if you know the interface is
// "internal" interface.
func NewSumInterface(namedInterface *types.Named, pkgscope []*types.Package) *SumInterface {
	seen := make(map[*types.Named]bool)
	return newSumInterface(seen, namedInterface, pkgscope)
}

func newSumInterface(seen map[*types.Named]bool, namedInterface *types.Named, pkgscope []*types.Package) *SumInterface {
	if _, ok := seen[namedInterface]; ok {
		return nil
	}
	seen[namedInterface] = true

	i, ok := namedInterface.Underlying().(*types.Interface)
	if !ok || i.Empty() {
		return nil
	}

	sum := &SumInterface{
		NamedInterface: namedInterface,
		IsInternal:     IsInternalInterface(i),
	}

	definedpkg := namedInterface.Obj().Pkg()
	pkgs := pkgscope
	if definedpkg != nil {
		if sum.IsInternal {
			// package scope includes only defined package for "internal" interface.
			pkgs = []*types.Package{definedpkg}
		} else {
			// add defined package to package scope if the scope doesn't have the
			// defiend package.
			found := false
			for _, p := range pkgs {
				found = (p == definedpkg)
				if found {
					break
				}
			}
			if !found {
				pkgs = append(pkgs, definedpkg)
			}
		}
	}
	sum.PkgScope = pkgs

	// Test assignability of all distinct pairs of
	// named types (T, U) where U is an interface.
	U := namedInterface
	for _, pkg := range pkgs {
		for _, name := range pkg.Scope().Names() {
			if obj, ok := pkg.Scope().Lookup(name).(*types.TypeName); ok {
				T := obj.Type()
				if T == U {
					continue
				}
				if types.AssignableTo(T, U) { // as interface
					if i := newSumInterface(seen, T.(*types.Named), pkgscope); i != nil {
						sum.Implements.Interfaces = append(sum.Implements.Interfaces, i)
					}
				} else if !types.IsInterface(T) { // as pointer
					if ptr := types.NewPointer(T); types.AssignableTo(ptr, U) {
						sum.Implements.Pointers = append(sum.Implements.Pointers, ptr)
					}
				}
			}
		}
	}

	return sum
}

// IsInternalInterface returns true if the given interface has unexported
// method.
func IsInternalInterface(iface *types.Interface) bool {
	for i := 0; i < iface.NumMethods(); i++ {
		if !iface.Method(i).Exported() {
			return true
		}
	}
	return false
}
