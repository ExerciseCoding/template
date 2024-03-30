package sql_demo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonColumn[T any] struct {
	Val T

	// NULL的问题
	Valid bool
}


func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}

	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch data := src.(type) {
	case string:
		// 可以考虑额外处理空字符串
		bs = []byte(data)
	case []byte:
		// 可以额外考虑处理 []byte{}
		bs = data
	case nil:
		// 说明数据库里面存的就是NULL
		return nil
	default:
		return errors.New("不支持类型")

	}
	err := json.Unmarshal(bs, &j.Val)
	if err == nil {
		j.Valid = true
	}
	return err
}