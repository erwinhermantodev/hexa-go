package templates

// DynamicHandlerTemplate generates handler files for models
const DynamicHandlerTemplate = `package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"{{.Config.ModuleName}}/model"
	"{{.Config.ModuleName}}/service"
	"{{.Config.ModuleName}}/utils"
)

type {{.Model.Name}}Handler struct {
	{{ToLower .Model.Name}}Service *service.{{.Model.Name}}Service
	validator *utils.Validator
}

func New{{.Model.Name}}Handler({{ToLower .Model.Name}}Service *service.{{.Model.Name}}Service, validator *utils.Validator) *{{.Model.Name}}Handler {
	return &{{.Model.Name}}Handler{
		{{ToLower .Model.Name}}Service: {{ToLower .Model.Name}}Service,
		validator: validator,
	}
}

// Create{{.Model.Name}} creates a new {{ToLower .Model.Name}}
// @Summary Create {{ToLower .Model.Name}}
// @Description Create a new {{ToLower .Model.Name}} with the provided data
// @Tags {{ToLower .Model.Name}}
// @Accept json
// @Produce json
// @Param {{ToLower .Model.Name}} body model.{{.Model.Name}}Request true "{{.Model.Name}} data"
// @Success 201 {object} model.{{.Model.Name}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{ToLower .Model.Name}}s [post]
func (h *{{.Model.Name}}Handler) Create{{.Model.Name}}(c echo.Context) error {
	var req model.{{.Model.Name}}Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body", 
			"details": err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed", 
			"details": err.Error(),
		})
	}

	{{ToLower .Model.Name}}, err := h.{{ToLower .Model.Name}}Service.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create {{ToLower .Model.Name}}", 
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "{{.Model.Name}} created successfully",
		"data":    {{ToLower .Model.Name}},
	})
}

// Get{{.Model.Name}} retrieves a {{ToLower .Model.Name}} by ID
// @Summary Get {{ToLower .Model.Name}} by ID
// @Description Get a single {{ToLower .Model.Name}} by its ID
// @Tags {{ToLower .Model.Name}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Model.Name}} ID"
// @Success 200 {object} model.{{.Model.Name}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /{{ToLower .Model.Name}}s/{id} [get]
func (h *{{.Model.Name}}Handler) Get{{.Model.Name}}(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		})
	}

	{{ToLower .Model.Name}}, err := h.{{ToLower .Model.Name}}Service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "{{.Model.Name}} not found", 
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "{{.Model.Name}} retrieved successfully",
		"data":    {{ToLower .Model.Name}},
	})
}

// GetAll{{.Model.Name}}s retrieves all {{ToLower .Model.Name}}s
// @Summary Get all {{ToLower .Model.Name}}s
// @Description Get a list of all {{ToLower .Model.Name}}s
// @Tags {{ToLower .Model.Name}}
// @Accept json
// @Produce json
// @Success 200 {array} model.{{.Model.Name}}Response
// @Failure 500 {object} map[string]interface{}
// @Router /{{ToLower .Model.Name}}s [get]
func (h *{{.Model.Name}}Handler) GetAll{{.Model.Name}}s(c echo.Context) error {
	{{ToLower .Model.Name}}s, err := h.{{ToLower .Model.Name}}Service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve {{ToLower .Model.Name}}s", 
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "{{.Model.Name}}s retrieved successfully",
		"data":    {{ToLower .Model.Name}}s,
		"count":   len({{ToLower .Model.Name}}s),
	})
}

// Update{{.Model.Name}} updates a {{ToLower .Model.Name}} by ID
// @Summary Update {{ToLower .Model.Name}}
// @Description Update a {{ToLower .Model.Name}} with the provided data
// @Tags {{ToLower .Model.Name}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Model.Name}} ID"
// @Param {{ToLower .Model.Name}} body model.{{.Model.Name}}Request true "Updated {{ToLower .Model.Name}} data"
// @Success 200 {object} model.{{.Model.Name}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{ToLower .Model.Name}}s/{id} [put]
func (h *{{.Model.Name}}Handler) Update{{.Model.Name}}(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		})
	}

	var req model.{{.Model.Name}}Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body", 
			"details": err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed", 
			"details": err.Error(),
		})
	}

	{{ToLower .Model.Name}}, err := h.{{ToLower .Model.Name}}Service.Update(uint(id), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update {{ToLower .Model.Name}}", 
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "{{.Model.Name}} updated successfully",
		"data":    {{ToLower .Model.Name}},
	})
}

// Delete{{.Model.Name}} deletes a {{ToLower .Model.Name}} by ID
// @Summary Delete {{ToLower .Model.Name}}
// @Description Delete a {{ToLower .Model.Name}} by its ID
// @Tags {{ToLower .Model.Name}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Model.Name}} ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{ToLower .Model.Name}}s/{id} [delete]
func (h *{{.Model.Name}}Handler) Delete{{.Model.Name}}(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid ID format",
		})
	}

	if err := h.{{ToLower .Model.Name}}Service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete {{ToLower .Model.Name}}", 
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "{{.Model.Name}} deleted successfully",
	})
}
`

// CustomHandlerTemplate generates standalone handler files
const CustomHandlerTemplate = `package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"{{.Config.ModuleName}}/service"
	"{{.Config.ModuleName}}/utils"
)

type {{.HandlerName}}Handler struct {
	{{ToLower .HandlerName}}Service *service.{{.HandlerName}}Service
	validator *utils.Validator
}

func New{{.HandlerName}}Handler({{ToLower .HandlerName}}Service *service.{{.HandlerName}}Service, validator *utils.Validator) *{{.HandlerName}}Handler {
	return &{{.HandlerName}}Handler{
		{{ToLower .HandlerName}}Service: {{ToLower .HandlerName}}Service,
		validator: validator,
	}
}

// Add your handler methods here
// Example:
// func (h *{{.HandlerName}}Handler) HandleSomething(c echo.Context) error {
//     return c.JSON(http.StatusOK, map[string]interface{}{"message": "Hello from {{.HandlerName}}Handler"})
// }

func (h *{{.HandlerName}}Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"handler": "{{.HandlerName}}Handler",
	})
}
`
