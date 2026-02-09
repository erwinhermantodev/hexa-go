package templates

// ModuleModelTemplate generates model files for a module
const ModuleModelTemplate = `package {{.PackageName}}

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

// ModuleRepositoryTemplate generates repository implementation for a module
const ModuleRepositoryTemplate = `package {{.PackageName}}

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(entity *{{.Model.Name}}) error
	GetByID(id uint) (*{{.Model.Name}}, error)
	GetAll() ([]{{.Model.Name}}, error)
	Update(entity *{{.Model.Name}}) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(entity *{{.Model.Name}}) error {
	return r.db.Create(entity).Error
}

func (r *repository) GetByID(id uint) (*{{.Model.Name}}, error) {
	var entity {{.Model.Name}}
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) GetAll() ([]{{.Model.Name}}, error) {
	var entities []{{.Model.Name}}
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository) Update(entity *{{.Model.Name}}) error {
	return r.db.Save(entity).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&{{.Model.Name}}{}, id).Error
}
`

// ModuleServiceTemplate generates service logic for a module
const ModuleServiceTemplate = `package {{.PackageName}}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *{{.Model.Name}}Request) (*{{.Model.Name}}Response, error) {
	entity := &{{.Model.Name}}{
{{- range .Model.Fields }}
{{- if and (ne .Name "ID") (ne .Name "CreatedAt") (ne .Name "UpdatedAt") (ne .Name "DeletedAt") }}
		{{.Name}}: req.{{.Name}},
{{- end }}
{{- end }}
	}

	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	return entity.ToResponse(), nil
}

func (s *Service) GetByID(id uint) (*{{.Model.Name}}Response, error) {
	entity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return entity.ToResponse(), nil
}

func (s *Service) GetAll() ([]{{.Model.Name}}Response, error) {
	entities, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []{{.Model.Name}}Response
	for _, entity := range entities {
		responses = append(responses, *entity.ToResponse())
	}

	return responses, nil
}

func (s *Service) Update(id uint, req *{{.Model.Name}}Request) (*{{.Model.Name}}Response, error) {
	entity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

{{- range .Model.Fields }}
{{- if and (ne .Name "ID") (ne .Name "CreatedAt") (ne .Name "UpdatedAt") (ne .Name "DeletedAt") }}
	entity.{{.Name}} = req.{{.Name}}
{{- end }}
{{- end }}

	if err := s.repo.Update(entity); err != nil {
		return nil, err
	}

	return entity.ToResponse(), nil
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}
`

// ModuleHandlerTemplate generates HTTP handlers for a module
const ModuleHandlerTemplate = `package {{.PackageName}}

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"{{.Config.ModuleName}}/utils"
)

type Handler struct {
	service   *Service
	validator *utils.Validator
}

func NewHandler(service *Service, validator *utils.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// Create handles POST /{{.PackageName}}s
func (h *Handler) Create(c echo.Context) error {
	var req {{.Model.Name}}Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	res, err := h.service.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

// GetByID handles GET /{{.PackageName}}s/:id
func (h *Handler) GetByID(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	res, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
	}

	return c.JSON(http.StatusOK, res)
}

// GetAll handles GET /{{.PackageName}}s
func (h *Handler) GetAll(c echo.Context) error {
	res, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// Update handles PUT /{{.PackageName}}s/:id
func (h *Handler) Update(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req {{.Model.Name}}Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// Delete handles DELETE /{{.PackageName}}s/:id
func (h *Handler) Delete(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
`
