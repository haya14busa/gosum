## Sum/Union/Variant Type in Go and Static Check Tool of switch-case handling  

[![GoDoc](https://godoc.org/github.com/haya14busa/gosum?status.svg)](https://godoc.org/github.com/haya14busa/gosum)
[![LICENSE](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## gosumcheck

### Installation 

```
$ go get -u github.com/haya14busa/gosum/cmd/gosumcheck
```

### Example

```
$ gosumcheck go/ast/...
/usr/lib/go/src/go/ast/commentmap.go:233:3: uncovered cases for ast.Node type switch:
        - *ast.ChanType
        - *ast.Ident
        - *ast.SelectorExpr
        - *ast.TypeAssertExpr
        - *ast.CompositeLit
        - *ast.FieldList
        - *ast.Package
        - *ast.ArrayType
        - *ast.ParenExpr
        - *ast.BinaryExpr
        - *ast.UnaryExpr
        - *ast.BadExpr
        - *ast.FuncLit
        - *ast.CommentGroup
        - *ast.IndexExpr
        - *ast.MapType
        - *ast.StructType
        - *ast.BasicLit
        - *ast.Ellipsis
        - *ast.InterfaceType
        - *ast.Comment
        - *ast.FuncType
        - *ast.SliceExpr
        - *ast.StarExpr
        - *ast.CallExpr
        - *ast.KeyValueExpr
/usr/lib/go/src/go/ast/filter.go:158:2: uncovered cases for ast.Spec type switch:
        - *ast.ImportSpec
/usr/lib/go/src/go/ast/filter.go:209:2: uncovered cases for ast.Decl type switch:
        - *ast.BadDecl
```

See `gosumcheck -h` and [Sum/Union/Variant Type in Go and Static Check Tool of switch-case handling](https://medium.com/@haya14busa/sum-union-variant-type-in-go-and-static-check-tool-of-switch-case-handling-3bfc61618b1e) for detail.
