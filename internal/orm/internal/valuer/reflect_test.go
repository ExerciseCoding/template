package valuer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ExerciseCoding/template/internal/orm/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReflectValue_SetColums(t *testing.T) {
	testSetColumns(t, NewReflectValue)
}

func testSetColumns(t *testing.T, creator Creator) {
	testCases := []struct {
		name string

		entity any
		rows   func() *sqlmock.Rows

		wantErr    error
		wantEntity any
	}{
		{
			name: "set colums",

			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "LiNing", "18", "Jerry")
				return rows
			},

			wantEntity: &TestModel{
				Id:        1,
				FirstName: "LiNing",
				Age:       18,
				LastName:  sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试列的不同顺序
			name: "order",

			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "last_name", "first_name"})
				rows.AddRow("1", "18", "Jerry", "LiNing")
				return rows
			},

			wantEntity: &TestModel{
				Id:        1,
				FirstName: "LiNing",
				Age:       18,
				LastName:  sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试部分列
			name: "partial columns",

			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name"})
				rows.AddRow("1", "18", "LiNing")
				return rows
			},

			wantEntity: &TestModel{
				Id:        1,
				FirstName: "LiNing",
				Age:       18,
				//LastName:  sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}

	r := model.NewRegistry()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造rows
			mockRows := tc.rows()
			mock.ExpectQuery("SELECT XXX").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT XXX")
			require.NoError(t, err)

			rows.Next()

			m, err := r.Get(tc.entity)
			require.NoError(t, err)
			val := creator(m, tc.entity)
			err = val.SetColums(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			// 比较一下tc.entity有没有设置好数据
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  sql.NullString
}
