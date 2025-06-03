package config

// ProjectConfig represents the configuration for a Go project
type ProjectConfig struct {
	Name        string
	ModuleName  string
	Description string
	Author      string
	Models      []ModelConfig
	Services    []string
}

// ModelConfig represents configuration for a model
type ModelConfig struct {
	Name       string
	Fields     []FieldConfig
	HasRepo    bool
	HasService bool
	HasHandler bool
}

// FieldConfig represents configuration for a model field
type FieldConfig struct {
	Name     string
	Type     string
	Tag      string
	Validate string
}

// DefaultUserModel returns the default User model configuration
func DefaultUserModel() ModelConfig {
	return ModelConfig{
		Name: "User",
		Fields: []FieldConfig{
			{Name: "ID", Type: "uint", Tag: "`gorm:\"primaryKey\" json:\"id\"`"},
			{Name: "Name", Type: "string", Tag: "`gorm:\"not null\" json:\"name\"`", Validate: "required,min=2,max=100"},
			{Name: "Email", Type: "string", Tag: "`gorm:\"unique;not null\" json:\"email\"`", Validate: "required,email"},
			{Name: "Password", Type: "string", Tag: "`gorm:\"not null\" json:\"-\"`", Validate: "required,min=6"},
			{Name: "IsActive", Type: "bool", Tag: "`gorm:\"default:true\" json:\"is_active\"`"},
			{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\"`"},
			{Name: "UpdatedAt", Type: "time.Time", Tag: "`json:\"updated_at\"`"},
			{Name: "DeletedAt", Type: "gorm.DeletedAt", Tag: "`gorm:\"index\" json:\"-\"`"},
		},
		HasRepo: true, HasService: true, HasHandler: true,
	}
}

// DefaultModelFields returns default fields for any model
func DefaultModelFields() []FieldConfig {
	return []FieldConfig{
		{Name: "ID", Type: "uint", Tag: "`gorm:\"primaryKey\" json:\"id\"`"},
		{Name: "CreatedAt", Type: "time.Time", Tag: "`json:\"created_at\"`"},
		{Name: "UpdatedAt", Type: "time.Time", Tag: "`json:\"updated_at\"`"},
		{Name: "DeletedAt", Type: "gorm.DeletedAt", Tag: "`gorm:\"index\" json:\"-\"`"},
	}
}
