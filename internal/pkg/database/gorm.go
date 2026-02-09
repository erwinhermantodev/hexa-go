package database

import (
	"context"

	"github.com/erwinhermantodev/hexa-go/internal/pkg/interfaces"
	"gorm.io/gorm"
)

type gormWrapper struct {
	db *gorm.DB
}

func NewGORM(db *gorm.DB) interfaces.Database {
	return &gormWrapper{db: db}
}

func (w *gormWrapper) Create(ctx context.Context, value interface{}) error {
	return w.db.WithContext(ctx).Create(value).Error
}

func (w *gormWrapper) Find(ctx context.Context, dest interface{}, conds ...interface{}) error {
	return w.db.WithContext(ctx).Find(dest, conds...).Error
}

func (w *gormWrapper) First(ctx context.Context, dest interface{}, conds ...interface{}) error {
	return w.db.WithContext(ctx).First(dest, conds...).Error
}

func (w *gormWrapper) Save(ctx context.Context, value interface{}) error {
	return w.db.WithContext(ctx).Save(value).Error
}

func (w *gormWrapper) Delete(ctx context.Context, value interface{}, conds ...interface{}) error {
	return w.db.WithContext(ctx).Delete(value, conds...).Error
}

func (w *gormWrapper) Raw(ctx context.Context, sql string, values ...interface{}) error {
	return w.db.WithContext(ctx).Raw(sql, values...).Error
}

func (w *gormWrapper) Exec(ctx context.Context, sql string, values ...interface{}) error {
	return w.db.WithContext(ctx).Exec(sql, values...).Error
}
