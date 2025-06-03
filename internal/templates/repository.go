package templates

// DynamicRepositoryTemplate generates repository files
const DynamicRepositoryTemplate = `package repository

import (
	"{{.Config.ModuleName}}/model"
	"gorm.io/gorm"
)

type {{ToLower .Model.Name}}Repository struct {
	db *gorm.DB
}

func New{{.Model.Name}}Repository(db *gorm.DB) {{.Model.Name}}Repository {
	return &{{ToLower .Model.Name}}Repository{db: db}
}

func (r *{{ToLower .Model.Name}}Repository) Create({{ToLower .Model.Name}} *model.{{.Model.Name}}) error {
	return r.db.Create({{ToLower .Model.Name}}).Error
}

func (r *{{ToLower .Model.Name}}Repository) GetByID(id uint) (*model.{{.Model.Name}}, error) {
	var {{ToLower .Model.Name}} model.{{.Model.Name}}
	err := r.db.First(&{{ToLower .Model.Name}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{ToLower .Model.Name}}, nil
}

func (r *{{ToLower .Model.Name}}Repository) GetAll() ([]model.{{.Model.Name}}, error) {
	var {{ToLower .Model.Name}}s []model.{{.Model.Name}}
	err := r.db.Find(&{{ToLower .Model.Name}}s).Error
	return {{ToLower .Model.Name}}s, err
}

func (r *{{ToLower .Model.Name}}Repository) Update({{ToLower .Model.Name}} *model.{{.Model.Name}}) error {
	return r.db.Save({{ToLower .Model.Name}}).Error
}

func (r *{{ToLower .Model.Name}}Repository) Delete(id uint) error {
	return r.db.Delete(&model.{{.Model.Name}}{}, id).Error
}

// Add custom query methods here
func (r *{{ToLower .Model.Name}}Repository) FindBy(field string, value interface{}) ([]model.{{.Model.Name}}, error) {
	var {{ToLower .Model.Name}}s []model.{{.Model.Name}}
	err := r.db.Where(field+" = ?", value).Find(&{{ToLower .Model.Name}}s).Error
	return {{ToLower .Model.Name}}s, err
}
`
