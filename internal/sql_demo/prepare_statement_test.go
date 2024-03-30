package sql_demo

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPrepare(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` where `id`=?")
	require.NoError(t, err)

	// id=1
	rows, err := stmt.QueryContext(ctx, 1)
	require.NoError(t, err)
	for rows.Next() {
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
	}
	// 整个应用关闭的时候调用
	stmt.Close()



	// prepare的缺点
	//stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` where `id` IN (?, ?, ?)")
	//stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` where `id` IN (?, ?, ?,?)")

}

type TestModel struct {
	Id   int64 `eorm:"auto_increament,primary_key"`
	FirstName string
	Age  int8
	// sql.NullString 主要用于处理那些数据库表中允许为空的字符串列，例如在 SQL 数据库中定义为 VARCHAR NULL 的列。通常，当您从数据库中查询这样的列时，结果会以 sql.NullString 的形式返回。
	LastName   *sql.NullString
}