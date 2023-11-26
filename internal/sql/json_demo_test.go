package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type JsonColumn[T any] struct {
	Val 	T
	Valid 	bool //标记数据库存的是不是NULL
}




func TestJsonColumn(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	ctx,cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS user_tab(
		id INTEGER PRIMARY KEY,
		address TEXT NOT NULL
	)
	`)
	if err != nil {
		t.Fatal(err)
	}

	re, err := db.ExecContext(context.Background(), "INSERT INTO `user_tab`(`id`,`address`) VALUES (?,?)",4,JsonColumn[Address]{Val: Address{Province: "广东",City: "深圳"}})
	if err != nil {
		t.Fatal(err)
	}
	affected, err := re.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}


	row := db.QueryRowContext(context.Background(), "SELECT `id`, `address` FROM `user_tab` LIMIT 1")
	if row.Err() != nil {
		t.Fatal(row.Err())
	}

	u := User{}
	err = row.Scan(&u.Id, &u.Address)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "深圳", u.Address.Val.City)
}


type User struct {
	Id 			int64
	Address  	JsonColumn[Address]
}


type Address struct {
	Province 	string
	City    	string
}

// Value 用于查询参数
func (j JsonColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	bs := src.([]byte)
	if len(bs) == 0 {
		return nil
	}
	err := json.Unmarshal(bs, &j.Val)
	j.Valid = err == nil
	return err
}
