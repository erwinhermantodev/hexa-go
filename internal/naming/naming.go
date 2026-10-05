// Package naming is the single place that turns a user-supplied name
// ("order-item", "orderItem", "Order Item") into the identifiers, file names
// and URL segments the generator needs.
package naming

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"
	"unicode"
)

var validName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*([ _-][A-Za-z0-9]+)*$`)

// reservedPackages are names a generated module may not take, because they
// would clash with imports in the generated main.go or routes.go.
var reservedPackages = map[string]bool{
	"main": true, "utils": true, "routes": true, "handler": true, "service": true,
	"repository": true, "model": true, "docs": true, "grpc": true, "echo": true,
	"gorm": true, "log": true, "http": true, "os": true, "fmt": true, "context": true,
	"time": true, "signal": true, "syscall": true, "zerolog": true, "middleware": true,
}

// irregular holds plurals the suffix rules below get wrong.
var irregular = map[string]string{
	"person": "people", "child": "children", "man": "men", "woman": "women",
	"mouse": "mice", "foot": "feet", "tooth": "teeth", "goose": "geese",
	"datum": "data", "criterion": "criteria", "analysis": "analyses",
	"status": "statuses", "bus": "buses",
}

// uncountable nouns keep their form; the plural handler would otherwise equal
// the singular one and clash.
var uncountable = map[string]bool{
	"equipment": true, "information": true, "metadata": true, "news": true,
	"series": true, "species": true, "sheep": true, "fish": true,
}

// Validate reports whether s can be used as a model, service, handler or
// module name.
func Validate(s string) error {
	if !validName.MatchString(s) {
		return fmt.Errorf("invalid name %q: use letters and digits, optionally separated by '-', '_' or spaces, starting with a letter", s)
	}
	if token.IsKeyword(strings.ToLower(s)) {
		return fmt.Errorf("invalid name %q: Go keyword", s)
	}
	return nil
}

// ValidatePackage additionally checks that s is usable as a Go package name
// inside the generated project.
func ValidatePackage(s string) error {
	if err := Validate(s); err != nil {
		return err
	}
	if reservedPackages[Package(s)] {
		return fmt.Errorf("invalid module name %q: %q would clash with a package already used by the generated project", s, Package(s))
	}
	return nil
}

// words splits a name into words on separators and camel-case boundaries,
// keeping acronyms together: "HTTPServer" -> [HTTP Server], "user_id" -> [user id].
func words(s string) []string {
	var out []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			out = append(out, string(cur))
			cur = nil
		}
	}
	runes := []rune(s)
	for i, r := range runes {
		switch {
		case r == ' ' || r == '_' || r == '-':
			flush()
		case unicode.IsUpper(r) && len(cur) > 0:
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower) {
				flush()
			}
			cur = append(cur, r)
		default:
			cur = append(cur, r)
		}
	}
	flush()
	return out
}

// Pascal returns the exported Go identifier: "order-item" -> "OrderItem".
// Names that are already PascalCase are returned unchanged.
func Pascal(s string) string {
	var b strings.Builder
	for _, w := range words(s) {
		r := []rune(w)
		if unicode.IsUpper(r[0]) && allUpperOrDigit(r) && len(r) > 1 {
			b.WriteString(w) // acronym such as ID or HTTP
			continue
		}
		b.WriteRune(unicode.ToUpper(r[0]))
		b.WriteString(string(r[1:]))
	}
	return b.String()
}

func allUpperOrDigit(r []rune) bool {
	for _, c := range r {
		if !unicode.IsUpper(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// Camel returns the unexported identifier: "OrderItem" -> "orderItem",
// "IDCard" -> "idCard".
func Camel(s string) string {
	ws := words(s)
	if len(ws) == 0 {
		return ""
	}
	rest := Pascal(strings.Join(ws[1:], " "))
	return strings.ToLower(ws[0]) + rest
}

// Snake returns the file/column style name: "OrderItem" -> "order_item".
func Snake(s string) string { return joinLower(s, "_") }

// Kebab returns the URL style name: "OrderItem" -> "order-item".
func Kebab(s string) string { return joinLower(s, "-") }

// Package returns the Go package name: "OrderItem" -> "orderitem".
func Package(s string) string { return joinLower(s, "") }

func joinLower(s, sep string) string {
	ws := words(s)
	for i, w := range ws {
		ws[i] = strings.ToLower(w)
	}
	return strings.Join(ws, sep)
}

// Plural pluralises the last word of a name and returns it in Pascal form:
// "category" -> "Categories", "OrderItem" -> "OrderItems", "Address" -> "Addresses".
func Plural(s string) string {
	ws := words(s)
	if len(ws) == 0 {
		return s
	}
	last := ws[len(ws)-1]
	lower := strings.ToLower(last)

	var plural string
	switch {
	case len([]rune(last)) > 1 && allUpperOrDigit([]rune(last)):
		plural = last + "s" // acronym: ID -> IDs
	case uncountable[lower]:
		plural = lower
	case irregular[lower] != "":
		plural = irregular[lower]
	case strings.HasSuffix(lower, "y") && len(lower) > 1 && !isVowel(rune(lower[len(lower)-2])):
		plural = lower[:len(lower)-1] + "ies"
	case hasAnySuffix(lower, "s", "x", "z", "ch", "sh"):
		plural = lower + "es"
	default:
		plural = lower + "s"
	}

	ws[len(ws)-1] = plural
	return Pascal(strings.Join(ws, " "))
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, x := range suffixes {
		if strings.HasSuffix(s, x) {
			return true
		}
	}
	return false
}

func isVowel(r rune) bool { return strings.ContainsRune("aeiou", r) }

// TrimSuffix removes a trailing component word the user may have typed out of
// habit: TrimSuffix("NotificationService", "Service") -> "Notification".
// A name that is only the suffix is returned unchanged.
func TrimSuffix(name, suffix string) string {
	if len(name) > len(suffix) && strings.EqualFold(name[len(name)-len(suffix):], suffix) {
		return name[:len(name)-len(suffix)]
	}
	return name
}
