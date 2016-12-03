package gosum

import "go/types"

// AllRelatedPackages returns all related package with the given package.
// Returned package list includes the given package too.
func AllRelatedPackages(pkg *types.Package) []*types.Package {
	pkgs := make(map[*types.Package]bool)
	visitPkgs(pkgs, pkg)
	ps := make([]*types.Package, 0, len(pkgs))
	for p := range pkgs {
		ps = append(ps, p)
	}
	return ps
}

func visitPkgs(pkgs map[*types.Package]bool, p *types.Package) {
	if _, ok := pkgs[p]; ok {
		return
	}
	pkgs[p] = true
	for _, pkg := range p.Imports() {
		visitPkgs(pkgs, pkg)
	}
}
