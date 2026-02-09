package interfaces

import "context"

type Database interface {
	Create(ctx context.Context, value interface{}) error
	Find(ctx context.Context, dest interface{}, conds ...interface{}) error
	First(ctx context.Context, dest interface{}, conds ...interface{}) error
	Save(ctx context.Context, value interface{}) error
	Delete(ctx context.Context, value interface{}, conds ...interface{}) error
	Raw(ctx context.Context, sql string, values ...interface{}) error
	Exec(ctx context.Context, sql string, values ...interface{}) error
}
