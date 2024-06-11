package orm

import (
	"context"
	"database/sql"
)

// Querier 用于 SELECT 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)

	// 用结构体的方式，指针方式容易发生内存逃逸
	//Get(ctx context.Context) (T, error)
	//GetMulti(ctx context.Context) ([]T, error)
}

// Executor 用于INSERT, DELETE, UPDATE
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
