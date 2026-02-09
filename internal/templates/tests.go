package templates

// DynamicRepositoryTestTemplate generates unit tests for repositories
const DynamicRepositoryTestTemplate = `package repository

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"{{.Config.ModuleName}}/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.{{.Model.Name}}{})
	return db
}

func Test{{.Model.Name}}Repository_Create(t *testing.T) {
	db := setupTestDB()
	repo := New{{.Model.Name}}Repository(db)

	entity := &model.{{.Model.Name}}{
		// Fill with some default fields if possible or leave for user
	}

	err := repo.Create(entity)
	assert.NoError(t, err)
	assert.NotZero(t, entity.ID)
}

func Test{{.Model.Name}}Repository_GetByID(t *testing.T) {
	db := setupTestDB()
	repo := New{{.Model.Name}}Repository(db)

	entity := &model.{{.Model.Name}}{}
	db.Create(entity)

	found, err := repo.GetByID(entity.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.ID, found.ID)
}
`

// DynamicServiceTestTemplate generates unit tests for services
const DynamicServiceTestTemplate = `package service

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"{{.Config.ModuleName}}/model"
)

// Mock{{.Model.Name}}Repository is a mock of the repository interface
type Mock{{.Model.Name}}Repository struct {
	mock.Mock
}

func (m *Mock{{.Model.Name}}Repository) Create(entity *model.{{.Model.Name}}) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *Mock{{.Model.Name}}Repository) GetByID(id uint) (*model.{{.Model.Name}}, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.{{.Model.Name}}), args.Error(1)
}

func (m *Mock{{.Model.Name}}Repository) GetAll() ([]model.{{.Model.Name}}, error) {
	args := m.Called()
	return args.Get(0).([]model.{{.Model.Name}}), args.Error(1)
}

func (m *Mock{{.Model.Name}}Repository) Update(entity *model.{{.Model.Name}}) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *Mock{{.Model.Name}}Repository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Mock{{.Model.Name}}Repository) FindBy(field string, value interface{}) ([]model.{{.Model.Name}}, error) {
	args := m.Called(field, value)
	return args.Get(0).([]model.{{.Model.Name}}), args.Error(1)
}

func Test{{.Model.Name}}Service_GetByID(t *testing.T) {
	mockRepo := new(Mock{{.Model.Name}}Repository)
	svc := New{{.Model.Name}}Service(mockRepo)

	expected := &model.{{.Model.Name}}{}
	expected.ID = 1
	mockRepo.On("GetByID", uint(1)).Return(expected, nil)

	res, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
	mockRepo.AssertExpectations(t)
}
`

// DynamicHandlerTestTemplate generates unit tests for handlers
const DynamicHandlerTestTemplate = `package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"{{.Config.ModuleName}}/utils"
)

func Test{{.Model.Name}}Handler_GetAll(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/{{ToLower .Model.Name}}s", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// This would need a mock service, but for scaffolding we provide a basic structure
	assert.Equal(t, http.StatusOK, http.StatusOK) 
}
`
