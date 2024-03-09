package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)


type sqlTestSuite struct {
	suite.Suite

	//配置字样
	driver  string
	dsn    string

	//初始化字段
	db *sql.DB
}

func (s *sqlTestSuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM test_model;")
	if err != nil {
		s.T().Fatal(err)
	}
}


func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	ctx,cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS test_model(
		id INTEGER PRIMARY KEY,
		first_name TEXT NOT NULL,
		age INTEGER,
		last_name TEXT NOT NULL
	)
	`)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TestCRUD() {
	t := s.T()
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (1, 'Tom', 18, 'Jerry')")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

	// 查询的时候使用复杂参数
	re, err := db.ExecContext(context.Background(), "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+"VALUES (?,?,?,?)",2,
		FullName{FirstName: "A",LastName: "B"},18,"Jerry")
	if err != nil {
		t.Fatal(err)
	}
	affected, err = re.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

	// 查询使用JSON
	re, err = db.ExecContext(context.Background(), "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+"VALUES (?,?,?,?)",3,
		FullNameJson{FirstName: "A",LastName: "B"},18,"Jerry")
	if err != nil {
		t.Fatal(err)
	}
	affected, err = re.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

	rows, err := db.QueryContext(context.Background(),"SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model` LIMIT ?",1)
	if err != nil {
		t.Fatal()
	}
	for rows.Next() {
		tm := &TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName,&tm.Age,&tm.LastName)
		// 常见错误，缺了指针
		// errs = rows.Scan(tm.Id,tm.FirstName,tm.Age, tm.LastName)
		if err != nil {
			rows.Close()
			t.Fatal(err)
		}
		assert.Equal(t, "Tom", tm.FirstName)
	}
	rows.Close()

	// 或者Exec
	res, err = db.ExecContext(ctx, "UPDATE `test_model` SET `first_name` = 'changed' WHERE `id` = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	affected, err = res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}


	row := db.QueryRowContext(context.Background(), "SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` LIMIT 1")
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	tm := &TestModel{}

	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "changed", tm.FirstName)

}

type FullName struct {
	FirstName string
	LastName string
}

func (f FullName) Value() (driver.Value, error) {
	return f.FirstName + f.LastName, nil
}


type FullNameJson struct {
	FirstName string
	LastName string
}

func (f FullNameJson) Value() (driver.Value, error) {
	return f.FirstName + f.LastName, nil
}

type TestModel struct {
	Id   int64 `eorm:"auto_increament,primary_key"`
	FirstName string
	Age  int8
	// sql.NullString 主要用于处理那些数据库表中允许为空的字符串列，例如在 SQL 数据库中定义为 VARCHAR NULL 的列。通常，当您从数据库中查询这样的列时，结果会以 sql.NullString 的形式返回。
	LastName   *sql.NullString
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "sqlite3",
		dsn:    "file:test.db?cache=shared&mode=memory",
	})
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(0)
	fmt.Println(timer.Stop())
	<- timer.C
}