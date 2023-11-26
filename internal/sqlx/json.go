package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonColumn 代表存储字段的json 类型
// 主要用于没有提供默认json类型的数据库
// T 可以是结构体，也可以是切片或者map
// 理论上来说一切可以被json库所处理的类型都能被用作T
// 不建议使用指针作为T的类型
// 如果T是指针，那么在Val为nil的情况下，一定要把Valid设置为false
type JsonColumn[T any] struct {
	Val  T
	Valid bool
}

// Value 返回一个json串.类型是[]byte
func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil,nil
	}
	res, err := json.Marshal(j.Val)
	return res,err
}

// Scan 将src转化为对象
// src 的类型必须是[]byte, string或者 nil
// 如果是nil, 我们不会做任何处理
func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case nil:
		return nil
	case []byte:
		bs = val
	case string:
		bs = []byte(val)
	default:
		return fmt.Errorf("JsonColumn.Scan 不支持 src 类型 %v", src)
	}

	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
