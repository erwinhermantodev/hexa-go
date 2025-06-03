package templates

// DynamicModelTemplate generates model files
const DynamicModelTemplate = `package model

import (
{{- if or (contains .Model.Fields "time.Time") (contains .Model.Fields "gorm.DeletedAt") }}
	"time"
	"gorm.io/gorm"
{{- end }}
)

type {{.Model.Name}} struct {
{{- range .Model.Fields }}
	{{.Name}} {{.Type}} {{.Tag}}
{{- end }}
}

{{- $modelName := .Model.Name }}
{{- $lowerModelName := ToLower .Model.Name }}

// {{.Model.Name}}Request represents the request structure for creating/updating {{ToLower .Model.Name}}
type {{.Model.Name}}Request struct {
{{- range .Model.Fields }}
{{- if and (ne .Name "ID") (ne .Name "CreatedAt") (ne .Name "UpdatedAt") (ne .Name "DeletedAt") }}
	{{.Name}} {{.Type}} ` + "`json:\"{{ToLower .Name}}\"{{if .Validate}} validate:\"{{.Validate}}\"{{end}}`" + `
{{- end }}
{{- end }}
}

// {{.Model.Name}}Response represents the response structure for {{ToLower .Model.Name}}
type {{.Model.Name}}Response struct {
{{- range .Model.Fields }}
{{- if ne .Name "DeletedAt" }}
	{{.Name}} {{.Type}} ` + "`json:\"{{ToLower .Name}}\"`" + `
{{- end }}
{{- end }}
}

// ToResponse converts {{.Model.Name}} to {{.Model.Name}}Response
func (m *{{.Model.Name}}) ToResponse() *{{.Model.Name}}Response {
	return &{{.Model.Name}}Response{
{{- range .Model.Fields }}
{{- if ne .Name "DeletedAt" }}
		{{.Name}}: m.{{.Name}},
{{- end }}
{{- end }}
	}
}
`
