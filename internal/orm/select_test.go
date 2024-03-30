package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelect_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name string

		builder QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no form",
			builder: NewSelector[TestModel](db),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](db).From("`test_model`"),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel](db).From(""),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel](db).From("`test_db`.`test_model`"),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](db).Where(),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Age").Eq(18)),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(C("Age").Eq(18))),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and",
			builder: NewSelector[TestModel](db).Where((C("Age").Eq(18)).And(C("FirstName").Eq("Tom"))),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "or",
			builder: NewSelector[TestModel](db).Where((C("Age").Eq(18)).Or(C("FirstName").Eq("Tom"))),

			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},

		{
			name:    "invalid column",
			builder: NewSelector[TestModel](db).Where((C("Age").Eq(18)).Or(C("abcd").Eq("Tom"))),
			wantErr: errs.NewErrUnkownField("abcd"),
		},

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}




func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	fmt.Println(mock)
	require.NoError(t, err)

	// 对应于query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	// 对应于no rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "NingLi", "18", "MingLi")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// no column
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("abc", "NingLi", "18", "MingLi")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct{
		name 	string

		s *Selector[TestModel]

		wantErr error
		wantRes *TestModel
	}{
		{
			name: "invalid query",
			s: NewSelector[TestModel](db).Where(C("XXX").Eq(1)),

			wantErr: errs.NewErrUnkownField("XXX"),
		},
		{
			name: "query error",
			s: NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},

		{
			name: "no rows",

			s: NewSelector[TestModel](db).Where(C("Id").Lt(2)),
			wantErr: errs.ErrNoRows,
		},
		{
			name: "data",

			s: NewSelector[TestModel](db).Where(C("Id").Eq(1)),

			wantRes: &TestModel{
				Id:        1,
				FirstName: "NingLi",
				Age:       18,
				LastName:  &sql.NullString{
					Valid: true,
					String: "MingLi",
				},
			},
		},
		{
			name: "no column",
			s: NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errs.NewErrUnkownColumn("abc"),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type TestModel struct {
	Id        	int64
	FirstName 	string
	Age       	int8
	LastName  	*sql.NullString
}



func memoryDB(t *testing.T, opts...DBOption) *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory",
		// 仅仅用于单元测试，不会发起真实查询
		opts...)
	require.NoError(t, err)
	return db
}