package orm

import (
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct{
		name  string

		entity  any
		wantModel *model
		wantErr error
	}{
		{
			name: "test model",
			entity: TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},

		},
		{
			name: "test model pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},

		},
		{
			name: "map",
			entity: map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name: "slice",
			entity: []int{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name: "base types",
			entity: 0,
			wantErr: errs.ErrPointerOnly,
		},
	}


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}