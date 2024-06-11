package valuer

import (
	"database/sql"
	"github.com/ExerciseCoding/template/internal/orm/model"
)

type Value interface {
	SetColums(rows *sql.Rows) error
}

type Creator func(model *model.Model, entity any) Value

type ValueV1 interface {
	SetColumn(entity any, rows *sql.Rows) error
}
