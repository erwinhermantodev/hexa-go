package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// AddImport adds an import path to a Go file if it doesn't already exist
func AddImport(filePath, importPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	if astutil.AddImport(fset, f, importPath) {
		return saveAST(filePath, fset, f)
	}
	return nil
}

// AddStructField adds a field to a struct in a Go file
func AddStructField(filePath, structName, fieldName, fieldType, tag string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != structName {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Check if field already exists
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				if name.Name == fieldName {
					return false
				}
			}
		}

		// Add new field
		newField := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fieldName)},
			Type:  ast.NewIdent(fieldType),
		}
		if tag != "" {
			newField.Tag = &ast.BasicLit{
				Kind:  token.STRING,
				Value: tag,
			}
		}

		st.Fields.List = append(st.Fields.List, newField)
		return false
	})

	return saveAST(filePath, fset, f)
}

// AddStatementToFunction adds a statement string to the end of a function body
func AddStatementToFunction(filePath, funcName, stmtCode string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Parse the statement code into AST nodes
	// We wrap it in a dummy package/function to parse safely
	exprTmpl := fmt.Sprintf("package p; func dummy() { %s }", stmtCode)
	dummyFset := token.NewFileSet()
	dummyFile, err := parser.ParseFile(dummyFset, "", exprTmpl, 0)
	if err != nil {
		return fmt.Errorf("failed to parse statement code: %v", err)
	}

	var newStmts []ast.Stmt
	ast.Inspect(dummyFile, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == "dummy" {
			newStmts = fn.Body.List
			return false
		}
		return true
	})

	if len(newStmts) == 0 {
		return fmt.Errorf("no valid statements found in code: %s", stmtCode)
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != funcName {
			return true
		}

		// For now, just append to the end. In the future we might want more complex positioning.
		fn.Body.List = append(fn.Body.List, newStmts...)
		found = true
		return false
	})

	if !found {
		return fmt.Errorf("function %s not found in %s", funcName, filePath)
	}

	return saveAST(filePath, fset, f)
}

// InjectCodeAST finds a marker comment and inserts a statement string before it
func InjectCodeAST(filePath, marker, stmtCode string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Parse code into statements
	exprTmpl := fmt.Sprintf("package p; func dummy() { %s }", stmtCode)
	dummyFset := token.NewFileSet()
	dummyFile, err := parser.ParseFile(dummyFset, "", exprTmpl, 0)
	if err != nil {
		return fmt.Errorf("failed to parse injection code: %v", err)
	}

	var newStmts []ast.Stmt
	ast.Inspect(dummyFile, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == "dummy" {
			newStmts = fn.Body.List
			return false
		}
		return true
	})

	if len(newStmts) == 0 {
		return fmt.Errorf("no valid statements found in: %s", stmtCode)
	}

	found := false
	// We need to find the function body that contains the marker comment
	// This is tricky because comments are often not "in" the body nodes
	// We'll use the File.Comments to find the marker and then find the nearest body

	var markerComment *ast.Comment
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, marker) {
				markerComment = c
				break
			}
		}
	}

	if markerComment == nil {
		return fmt.Errorf("marker %s not found in %s", marker, filePath)
	}

	// Inspect all function bodies to find which one includes the marker position
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if marker is within function range
		if markerComment.Pos() >= fn.Body.Lbrace && markerComment.Pos() <= fn.Body.Rbrace {
			// Insert before the identified statement (or at end of block)
			// Actually, if we want "Before marker", we need to find the statement that starts BEFORE the marker
			// but markers are usually on their own lines between statements.

			// Simple logic: insert at the first position where stmt.Pos > marker.Pos
			// This effectively places it right before the marker if the marker is above a statement.
			// Wait, if the marker is BELOW the statements we want to add, then insertIdx is correct.

			// Most hexa-go markers are BELOW the injected code:
			// repository.New...
			// // [REPOS-INIT]

			// So we want to insert BEFORE the statement that contains the marker,
			// but markers are comments, they are not statements.

			// Let's refine: find where to inject.
			newBody := make([]ast.Stmt, 0, len(fn.Body.List)+len(newStmts))
			injected := false
			for _, stmt := range fn.Body.List {
				if !injected && stmt.Pos() > markerComment.Pos() {
					newBody = append(newBody, newStmts...)
					injected = true
				}
				newBody = append(newBody, stmt)
			}
			if !injected {
				newBody = append(newBody, newStmts...)
			}

			fn.Body.List = newBody
			found = true
			return false
		}
		return true
	})

	if !found {
		// If not in a function body, maybe it's in top-level declarations?
		// For now, let's assume it's always in a function for Hexa-Go
		return fmt.Errorf("could not find function body containing marker %s", marker)
	}

	return saveAST(filePath, fset, f)
}

// saveAST formats and saves the AST back to the file
func saveAST(filePath string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	// Create a printer with custom config to preserve comments better
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
