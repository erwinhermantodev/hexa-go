package utils

import (
	"strings"
	"testing"

	"github.com/erwinhermantodev/hexa-go/internal/config"
)

func TestParseField(t *testing.T) {
	tests := []struct {
		spec      string
		name, typ string
		tag       string
		validate  string
	}{
		{"Title:string::required,min=5", "Title", "string", "`json:\"title\"`", "required,min=5"},
		{"Slug:string:unique:required", "Slug", "string", "`gorm:\"unique\" json:\"slug\"`", "required"},
		{"Published:bool:default=false", "Published", "bool", "`gorm:\"default:false\" json:\"published\"`", ""},
		{"IsPublished:bool:default=false:", "IsPublished", "bool", "`gorm:\"default:false\" json:\"is_published\"`", ""},
		{"PublishedAt:*time.Time", "PublishedAt", "*time.Time", "`json:\"published_at\"`", ""},
		{"DeletedAt:gorm.DeletedAt:index", "DeletedAt", "gorm.DeletedAt", "`gorm:\"index\" json:\"deleted_at\"`", ""},
		{"ID:uint", "ID", "uint", "`gorm:\"primaryKey\" json:\"id\"`", ""},
		{"ID:uint:primaryKey;autoIncrement", "ID", "uint", "`gorm:\"primaryKey;autoIncrement\" json:\"id\"`", ""},
		{"UserID:uint::required", "UserID", "uint", "`json:\"user_id\"`", "required"},
		{"price:float64::required,gt=0", "Price", "float64", "`json:\"price\"`", "required,gt=0"},
		{"Tags:[]string", "Tags", "[]string", "`json:\"tags\"`", ""},
		{"Tag:string:size=20;not null", "Tag", "string", "`gorm:\"size:20;not null\" json:\"tag\"`", ""},
	}
	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			got, err := ParseField(tt.spec)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := config.FieldConfig{Name: tt.name, Type: tt.typ, Tag: tt.tag, Validate: tt.validate}
			if got != want {
				t.Errorf("got %+v, want %+v", got, want)
			}
		})
	}
}

func TestParseFieldErrors(t *testing.T) {
	tests := []struct {
		spec    string
		wantErr string
	}{
		{"", "expected Name:Type"},
		{"Name", "expected Name:Type"},
		{"Name:string:a:b:c", "expected Name:Type"},
		{"1Name:string", "not a valid field name"},
		{"Na me:string", "not a valid field name"},
		{"Name:", "not a valid type"},
		{"Name:str ing", "not a valid type"},
		{"Slug:string::unique", "GORM option"},
		{"Active:bool::default=true", "GORM option"},
		{"Name:string:unique,index", "separate GORM options"},
		{"Name:string::required,", "empty validation tag"},
		{"Name:string::required\"", "quotes"},
	}
	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			_, err := ParseField(tt.spec)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got err %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestParseFieldsRejectsDuplicates(t *testing.T) {
	if _, err := ParseFields([]string{"Name:string", "name:string"}); err == nil {
		t.Error("expected duplicate error")
	}
}

func TestWithDefaultFields(t *testing.T) {
	fields, _ := ParseFields([]string{"Name:string", "CreatedAt:time.Time"})
	got := WithDefaultFields(fields)

	var names []string
	for _, f := range got {
		names = append(names, f.Name)
	}
	want := "ID,Name,CreatedAt,UpdatedAt,DeletedAt"
	if strings.Join(names, ",") != want {
		t.Errorf("got %v, want %s", names, want)
	}
	// The user's own CreatedAt must win over the default.
	if got[2].Tag != "`json:\"created_at\"`" {
		t.Errorf("user-defined field was replaced: %+v", got[2])
	}
}

func TestSnakeCase(t *testing.T) {
	for in, want := range map[string]string{
		"ID": "id", "Name": "name", "IsPublished": "is_published",
		"UserID": "user_id", "HTTPServer": "http_server", "CreatedAt": "created_at",
	} {
		if got := snakeCase(in); got != want {
			t.Errorf("snakeCase(%q) = %q, want %q", in, got, want)
		}
	}
}
