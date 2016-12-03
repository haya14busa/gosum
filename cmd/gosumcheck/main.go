package main

import (
	"os"

	"github.com/haya14busa/gosum/checker"
	"honnef.co/go/lint"
	"honnef.co/go/lint/lintutil"
)

func main() {
	funcs := []lint.Func{
		checker.CheckSwitch,
	}
	lintutil.ProcessArgs("gosumcheck", funcs, os.Args[1:])
}
