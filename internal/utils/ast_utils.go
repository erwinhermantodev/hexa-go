package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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

// StructHasField reports whether the named struct in a Go file already has a
// field with the given name.
func StructHasField(filePath, structName, fieldName string) (bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != structName {
			return true
		}
		if st, ok := ts.Type.(*ast.StructType); ok {
			for _, field := range st.Fields.List {
				for _, name := range field.Names {
					if name.Name == fieldName {
						found = true
					}
				}
			}
		}
		return false
	})
	return found, nil
}

// structFieldMarker is the comment in routes.go above which new Router fields
// are inserted.
const structFieldMarker = "// [HANDLER-FIELDS-EXPORTED]"

// AddStructField adds a field to a struct in a Go file. When the file has the
// handler-fields marker the field is inserted above it as text (keeping the
// layout clean); otherwise it is appended through the AST. A field that
// already exists is left alone.
func AddStructField(filePath, structName, fieldName, fieldType, tag string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	var st *ast.StructType
	ast.Inspect(f, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == structName {
			st, _ = ts.Type.(*ast.StructType)
			return false
		}
		return true
	})
	if st == nil {
		return fmt.Errorf("struct %s not found in %s", structName, filePath)
	}
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			if name.Name == fieldName {
				return nil
			}
		}
	}

	// Prefer the marker when it sits inside this struct.
	start, end := fset.Position(st.Fields.Opening).Offset, fset.Position(st.Fields.Closing).Offset
	if i := strings.Index(string(src), structFieldMarker); i > start && i < end {
		lineStart := strings.LastIndex(string(src[:i]), "\n") + 1
		line := fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tag)
		patched := string(src[:lineStart]) + line + string(src[lineStart:])
		formatted, err := format.Source([]byte(patched))
		if err != nil {
			return fmt.Errorf("injected field does not compile cleanly: %v", err)
		}
		return os.WriteFile(filePath, formatted, 0644)
	}

	newField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fieldName)},
		Type:  ast.NewIdent(fieldType),
	}
	if tag != "" {
		newField.Tag = &ast.BasicLit{Kind: token.STRING, Value: tag}
	}
	st.Fields.List = append(st.Fields.List, newField)
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

// InjectCodeAST inserts the statements in stmtCode on the line above the first
// comment containing marker. The marker must sit inside a function body.
// Statements already present in that body are skipped, and a `:=` that would
// redeclare a variable with a different value is an error.
//
// The code is inserted as text and the file re-formatted, rather than spliced
// into the AST: nodes parsed from another file carry foreign positions, which
// makes go/format wrap lines at random.
func InjectCodeAST(filePath, marker, stmtCode string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	dummyFset := token.NewFileSet()
	dummyFile, err := parser.ParseFile(dummyFset, "", "package p; func dummy() {\n"+stmtCode+"\n}", 0)
	if err != nil {
		return fmt.Errorf("failed to parse injection code: %v", err)
	}
	newStmts := dummyFile.Decls[0].(*ast.FuncDecl).Body.List
	if len(newStmts) == 0 {
		return fmt.Errorf("no valid statements found in: %s", stmtCode)
	}

	var markerComment *ast.Comment
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, marker) {
				markerComment = c
				break
			}
		}
		if markerComment != nil {
			break
		}
	}
	if markerComment == nil {
		return fmt.Errorf("marker %s not found in %s", marker, filePath)
	}

	var body *ast.BlockStmt
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil &&
			markerComment.Pos() >= fn.Body.Lbrace && markerComment.Pos() <= fn.Body.Rbrace {
			body = fn.Body
			break
		}
	}
	if body == nil {
		return fmt.Errorf("marker %s in %s is not inside a function body", marker, filePath)
	}

	existing := map[string]bool{}   // normalized statement text
	declared := map[string]string{} // variable -> normalized statement that defined it
	for _, stmt := range body.List {
		text := normalize(renderNode(fset, stmt))
		existing[text] = true
		for _, name := range definedNames(stmt) {
			declared[name] = text
		}
	}

	var out strings.Builder
	for _, stmt := range newStmts {
		text := normalize(renderNode(dummyFset, stmt))
		if existing[text] {
			continue
		}
		for _, name := range definedNames(stmt) {
			if _, taken := declared[name]; taken {
				return fmt.Errorf("%s already declares %q; remove it or choose another name", filepath.Base(filePath), name)
			}
		}
		out.WriteString("\t" + renderNode(dummyFset, stmt) + "\n")
		existing[text] = true
		for _, name := range definedNames(stmt) {
			declared[name] = text
		}
	}
	if out.Len() == 0 {
		return nil
	}

	offset := fset.Position(markerComment.Pos()).Offset
	lineStart := strings.LastIndex(string(src[:offset]), "\n") + 1
	patched := string(src[:lineStart]) + out.String() + string(src[lineStart:])

	formatted, err := format.Source([]byte(patched))
	if err != nil {
		return fmt.Errorf("injected code does not compile cleanly: %v", err)
	}
	return os.WriteFile(filePath, formatted, 0644)
}

// definedNames returns the variables a `:=` statement declares.
func definedNames(stmt ast.Stmt) []string {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE {
		return nil
	}
	var names []string
	for _, lhs := range assign.Lhs {
		if id, ok := lhs.(*ast.Ident); ok && id.Name != "_" {
			names = append(names, id.Name)
		}
	}
	return names
}

// normalize strips all whitespace so statements compare equal regardless of
// how they were wrapped.
func normalize(s string) string { return strings.Join(strings.Fields(s), "") }

// saveAST formats and saves the AST back to the file
func saveAST(filePath string, fset *token.FileSet, f *ast.File) error {
	var buf bytes.Buffer
	// Create a printer with custom config to preserve comments better
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// renderNode prints a node to a canonical string for comparison.
func renderNode(fset *token.FileSet, n ast.Node) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, n); err != nil {
		return ""
	}
	return buf.String()
}
