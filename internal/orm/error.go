package orm

import "github.com/ExerciseCoding/template/internal/orm/internal/errs"

// 通过这种形式将内部错误暴漏在外面
var ErrNoRows = errs.ErrNoRows