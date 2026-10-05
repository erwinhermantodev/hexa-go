package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/erwinhermantodev/hexa-go/internal/config"
)

// Field spec grammar, used by both the -f flag and the interactive prompt:
//
//	Name:Type[:GormOptions[:Validation]]
//
//	Name         Go field name; the first letter is upper-cased.
//	Type         string, int, uint, int64, float64, bool, time.Time, gorm.DeletedAt,
//	             a pointer to any of those (*time.Time), or a custom type name.
//	GormOptions  options for the gorm tag, separated by ";". A value is written
//	             with "=" because ":" separates the slots: unique;default=false
//	Validation   go-playground/validator tags, separated by ",": required,min=5
//
// Examples:
//
//	Title:string::required,min=5
//	Slug:string:unique:required
//	Published:bool:default=false
//	PublishedAt:*time.Time

const fieldSpecHelp = "expected Name:Type[:GormOptions[:Validation]]"

var (
	typePattern  = regexp.MustCompile(`^(\*|\[\])?[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)?$`)
	namePattern  = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	camelBoundry = regexp.MustCompile(`([a-z0-9])([A-Z])`)

	// gormOnly lists option names that belong in the GORM slot. Seeing one in
	// the validation slot is almost always a misplaced argument.
	gormOnly = map[string]bool{
		"unique": true, "index": true, "uniqueIndex": true, "primaryKey": true,
		"autoIncrement": true, "default": true, "not null": true, "column": true, "size": true,
	}
)

// ParseField parses one field spec into a FieldConfig.
func ParseField(spec string) (config.FieldConfig, error) {
	spec = strings.TrimSpace(spec)
	parts := strings.Split(spec, ":")
	if len(parts) < 2 || len(parts) > 4 {
		return config.FieldConfig{}, fmt.Errorf("invalid field %q: %s", spec, fieldSpecHelp)
	}

	name := strings.TrimSpace(parts[0])
	typ := strings.TrimSpace(parts[1])
	var gormSpec, validate string
	if len(parts) > 2 {
		gormSpec = strings.TrimSpace(parts[2])
	}
	if len(parts) > 3 {
		validate = strings.TrimSpace(parts[3])
	}

	if !namePattern.MatchString(name) {
		return config.FieldConfig{}, fmt.Errorf("invalid field %q: %q is not a valid field name", spec, name)
	}
	name = upperFirst(name)

	if !typePattern.MatchString(typ) {
		return config.FieldConfig{}, fmt.Errorf("invalid field %q: %q is not a valid type", spec, typ)
	}

	for _, s := range []string{gormSpec, validate} {
		if strings.ContainsAny(s, "`\"") {
			return config.FieldConfig{}, fmt.Errorf("invalid field %q: tags must not contain quotes or backticks", spec)
		}
	}

	gormOpts, err := parseGormOptions(spec, gormSpec)
	if err != nil {
		return config.FieldConfig{}, err
	}
	if name == "ID" && !hasOption(gormOpts, "primaryKey") {
		gormOpts = append([]string{"primaryKey"}, gormOpts...)
	}

	if err := checkValidation(spec, validate); err != nil {
		return config.FieldConfig{}, err
	}

	tag := fmt.Sprintf("json:\"%s\"", snakeCase(name))
	if len(gormOpts) > 0 {
		tag = fmt.Sprintf("gorm:\"%s\" %s", strings.Join(gormOpts, ";"), tag)
	}

	return config.FieldConfig{
		Name:     name,
		Type:     typ,
		Tag:      "`" + tag + "`",
		Validate: validate,
	}, nil
}

// ParseFields parses several specs and rejects duplicate field names.
func ParseFields(specs []string) ([]config.FieldConfig, error) {
	var fields []config.FieldConfig
	seen := map[string]bool{}
	for _, spec := range specs {
		f, err := ParseField(spec)
		if err != nil {
			return nil, err
		}
		if seen[f.Name] {
			return nil, fmt.Errorf("duplicate field %q", f.Name)
		}
		seen[f.Name] = true
		fields = append(fields, f)
	}
	return fields, nil
}

// WithDefaultFields prepends the standard ID/CreatedAt/UpdatedAt/DeletedAt
// fields that the generated repository, service and handler rely on, skipping
// any the user already defined.
func WithDefaultFields(fields []config.FieldConfig) []config.FieldConfig {
	defined := map[string]bool{}
	for _, f := range fields {
		defined[f.Name] = true
	}

	var out []config.FieldConfig
	var trailing []config.FieldConfig
	for _, d := range config.DefaultModelFields() {
		if defined[d.Name] {
			continue
		}
		if d.Name == "ID" {
			out = append(out, d)
		} else {
			trailing = append(trailing, d)
		}
	}
	out = append(out, fields...)
	return append(out, trailing...)
}

// parseGormOptions turns "unique;default=false" into ["unique", "default:false"].
func parseGormOptions(spec, s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	var opts []string
	for _, opt := range strings.Split(s, ";") {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		if strings.Contains(opt, ",") {
			return nil, fmt.Errorf("invalid field %q: separate GORM options with \";\", not \",\" (got %q)", spec, opt)
		}
		opts = append(opts, strings.Replace(opt, "=", ":", 1))
	}
	return opts, nil
}

func hasOption(opts []string, name string) bool {
	for _, o := range opts {
		if strings.SplitN(o, ":", 2)[0] == name {
			return true
		}
	}
	return false
}

// checkValidation rejects GORM options placed in the validation slot.
func checkValidation(spec, validate string) error {
	if validate == "" {
		return nil
	}
	for _, tag := range strings.Split(validate, ",") {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("invalid field %q: empty validation tag", spec)
		}
		name := strings.SplitN(tag, "=", 2)[0]
		if gormOnly[name] {
			return fmt.Errorf("invalid field %q: %q is a GORM option; put it in the third slot (Name:Type:%s)", spec, name, tag)
		}
	}
	return nil
}

func upperFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// snakeCase converts IsPublished to is_published and UserID to user_id.
func snakeCase(s string) string {
	s = camelBoundry.ReplaceAllString(s, "${1}_${2}")
	// Split a trailing acronym from the following word: HTTPServer -> HTTP_Server.
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) && unicode.IsUpper(runes[i-1]) {
			b.WriteRune('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
