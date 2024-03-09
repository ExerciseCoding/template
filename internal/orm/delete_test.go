package orm

import (
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete_Build(t *testing.T) {
	testCases := []struct {
		name string

		builder QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no form",
			builder: &Delete[TestModel]{},

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: (&Delete[TestModel]{}).From("`test_model`"),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&Delete[TestModel]{}).From(""),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&Delete[TestModel]{}).From("`test_db`.`test_model`"),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: (&Delete[TestModel]{}).Where(),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: (&Delete[TestModel]{}).Where(C("Age").Eq(18)),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not",
			builder: (&Delete[TestModel]{}).Where(Not(C("Age").Eq(18))),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and",
			builder: (&Delete[TestModel]{}).Where((C("Age").Eq(18)).And(C("FirstName").Eq("Tom"))),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "or",
			builder: (&Delete[TestModel]{}).Where((C("Age").Eq(18)).Or(C("FirstName").Eq("Tom"))),

			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "invalid column",
			builder: (&Delete[TestModel]{}).Where((C("Age").Eq(18)).Or(C("abcd").Eq("Tom"))),
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

