package imports

import (
	"go/ast"
	"go/token"
)

func removeLines(fset *token.FileSet, f *ast.File) {
	for i, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if len(d.Specs) == 0 {
			// Empty import block, remove it.
			f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		}

		if !d.Lparen.IsValid() {
			// Not a block: sorted by default.
			continue
		}

		lastLine := fset.File(d.Lparen).Line(d.Lparen)
		for _, s := range d.Specs {
			newLine := fset.File(s.Pos()).Line(s.Pos())
			emptyLines := newLine - lastLine - 1
			for i := 0; i < emptyLines; i++ {
				fset.File(s.Pos()).MergeLine(lastLine)
			}
			lastLine = fset.File(s.Pos()).Line(s.Pos())
		}
	}
}
