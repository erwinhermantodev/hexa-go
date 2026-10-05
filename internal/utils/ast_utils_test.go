package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const injectSrc = `package main

func main() {
	a := 1
	_ = a

	// [INIT]

	done()
}
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "main.go")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func read(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestInjectCodeASTInsertsAboveMarkerAndStaysFormatted(t *testing.T) {
	p := writeTemp(t, injectSrc)
	if err := InjectCodeAST(p, "// [INIT]", "b := compute(a,\n\t\t3)\n_ = b"); err != nil {
		t.Fatal(err)
	}
	got := read(t, p)
	// Statements land directly above the marker, with the marker kept.
	want := "\tb := compute(a,\n\t\t3)\n\t_ = b\n\t// [INIT]\n"
	if !strings.Contains(got, want) {
		t.Errorf("unexpected output:\n%s", got)
	}
}

func TestInjectCodeASTIsIdempotent(t *testing.T) {
	p := writeTemp(t, injectSrc)
	for i := 0; i < 3; i++ {
		if err := InjectCodeAST(p, "// [INIT]", "b := compute(a)"); err != nil {
			t.Fatal(err)
		}
	}
	if n := strings.Count(read(t, p), "b := compute(a)"); n != 1 {
		t.Errorf("statement injected %d times, want 1", n)
	}
}

func TestInjectCodeASTRejectsConflictingDeclaration(t *testing.T) {
	p := writeTemp(t, injectSrc)
	if err := InjectCodeAST(p, "// [INIT]", "b := first()"); err != nil {
		t.Fatal(err)
	}
	before := read(t, p)
	err := InjectCodeAST(p, "// [INIT]", "b := second()")
	if err == nil || !strings.Contains(err.Error(), `already declares "b"`) {
		t.Errorf("got %v, want redeclaration error", err)
	}
	if read(t, p) != before {
		t.Error("file changed despite the error")
	}
}

func TestInjectCodeASTErrors(t *testing.T) {
	p := writeTemp(t, injectSrc)
	if err := InjectCodeAST(p, "// [MISSING]", "x := 1"); err == nil {
		t.Error("missing marker: expected error")
	}
	if err := InjectCodeAST(p, "// [INIT]", "x :="); err == nil {
		t.Error("invalid code: expected error")
	}
	outside := writeTemp(t, "package main\n\n// [INIT]\n\nfunc main() {}\n")
	if err := InjectCodeAST(outside, "// [INIT]", "x := 1"); err == nil {
		t.Error("marker outside a function: expected error")
	}
}

func TestAddStructFieldUsesMarkerAndSkipsExisting(t *testing.T) {
	p := writeTemp(t, "package routes\n\ntype Router struct {\n\t// [HANDLER-FIELDS-EXPORTED]\n}\n")
	for i := 0; i < 2; i++ {
		if err := AddStructField(p, "Router", "UserHandler", "*handler.UserHandler", ""); err != nil {
			t.Fatal(err)
		}
	}
	got := read(t, p)
	if strings.Count(got, "UserHandler *") != 1 {
		t.Errorf("field added more than once:\n%s", got)
	}
	if !strings.Contains(got, "\tUserHandler *handler.UserHandler\n\t// [HANDLER-FIELDS-EXPORTED]") {
		t.Errorf("unexpected layout:\n%s", got)
	}
	if err := AddStructField(p, "Missing", "X", "int", ""); err == nil {
		t.Error("unknown struct: expected error")
	}
	// Without the marker the field is appended through the AST.
	q := writeTemp(t, "package routes\n\ntype Router struct {\n\tA int\n}\n")
	if err := AddStructField(q, "Router", "B", "string", ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(read(t, q), "B string") {
		t.Errorf("AST fallback did not add the field:\n%s", read(t, q))
	}
}
