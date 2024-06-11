package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/ExerciseCoding/template/internal/orm/internal/valuer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct{
		name  string

		s QueryBuilder
		wantError error
		wantQuery *Query
	}{
		{
			name: "multiple columns",
			s: NewSelector[TestModel](db).Select(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "SELECT `first_name`,`last_name` FROM `test_model`;",
			},
		},
		{
			name: "invalid column",
			s: NewSelector[TestModel](db).Select(C("Invalid")),
			wantError: errs.NewErrUnkownField("Invalid"),
		},
		{
			name: "Avg",
			s: NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},

		{
			name: "Count",
			s: NewSelector[TestModel](db).Select(Count("Age")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},

		{
			name: "Max",
			s: NewSelector[TestModel](db).Select(Max("Age")),
			wantQuery: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},

		{
			name: "Min",
			s: NewSelector[TestModel](db).Select(Min("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},

		{
			name: "aggregate invalid columns",
			s: NewSelector[TestModel](db).Select(Avg("invalid")),
			wantError: errs.NewErrUnkownField("invalid"),
		},

		{
			name: "multiple aggregate",
			s: NewSelector[TestModel](db).Select(Avg("Age"),Count("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`),COUNT(`age`) FROM `test_model`;",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.s.Build()
			assert.Equal(t, tc.wantError, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

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

	testCases := []struct {
		name string

		s *Selector[TestModel]

		wantErr error
		wantRes *TestModel
	}{
		{
			name: "invalid query",
			s:    NewSelector[TestModel](db).Where(C("XXX").Eq(1)),

			wantErr: errs.NewErrUnkownField("XXX"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},

		{
			name: "no rows",

			s:       NewSelector[TestModel](db).Where(C("Id").Lt(2)),
			wantErr: errs.ErrNoRows,
		},
		{
			name: "data",

			s: NewSelector[TestModel](db).Where(C("Id").Eq(1)),

			wantRes: &TestModel{
				Id:        1,
				FirstName: "NingLi",
				Age:       18,
				LastName: &sql.NullString{
					Valid:  true,
					String: "MingLi",
				},
			},
		},
		{
			name:    "no column",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errs.NewErrUnkownColumn("abc"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.GetV1(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func memoryDB(t *testing.T, opts ...DBOption) *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory",
		// 仅仅用于单元测试，不会发起真实查询
		opts...)
	require.NoError(t, err)
	return db
}

func (TestModel) CreateSQL() string {
	return `CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)`
}

func BenchmarkQuerier_Get(b *testing.B) {
	db, err := Open("sqlite3", fmt.Sprintf("file:benchmark_get.db?cache=shared&mode=memory"))
	if err != nil {
		b.Fatal(err)
	}
	_, err = db.db.Exec(TestModel{}.CreateSQL())
	if err != nil {
		b.Fatal(err)
	}
	res, err := db.db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 12, "Deng", 18, "Ming")
	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}

	b.Run("unsafe", func(b *testing.B) {
		db.creator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.creator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}

	})

}
