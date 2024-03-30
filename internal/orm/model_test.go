package orm

import (
	"database/sql"
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_registry_Register(t *testing.T) {
	testCases := []struct{
		name  string

		entity  any
		wantModel *Model
		fields []*Field
		wantErr error
		opts []ModelOpt
	}{
		{
			name: "test model",
			entity: TestModel{},
			wantModel: &Model{
				tableName: "test_model",
			},
			fields: []*Field{
				{
					colName: "id",
					goName: "Id",
					typ: reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					goName: "LastName",
					typ: reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					goName: "Age",
					typ: reflect.TypeOf(int8(0)),
				},
			},
		wantErr: errs.ErrPointerOnly,
		},
		//{
		//	name: "test model pointer",
		//	entity: &TestModel{},
		//	wantModel: &Model{
		//		tableName: "test_model",
		//		fieldMap: map[string]*Field{
		//			"Id": {
		//				colName: "id",
		//			},
		//			"FirstName": {
		//				colName: "first_name",
		//			},
		//			"LastName": {
		//				colName: "last_name",
		//			},
		//			"Age": {
		//				colName: "age",
		//			},
		//		},
		//	},
		//
		//},
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

	r := registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tc.fields {
				fieldMap[f.goName] = f
				columnMap[f.colName] = f
			}
			tc.wantModel.fieldMap = fieldMap
			tc.wantModel.columnMap = columnMap
			assert.Equal(t, tc.wantModel, m)
		})
	}
}



func TestRegistry_get(t *testing.T) {
	testCases := []struct{
		name  string

		entity any
		fields []*Field
		wantModel  *Model
		wantErr error

	}{
		{
			name: "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				//fieldMap: map[string]*Field{
				//	"Id": {
				//		colName: "id",
				//	},
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//	"LastName": {
				//		colName: "last_name",
				//	},
				//	"Age": {
				//		colName: "age",
				//	},
				//},
			},
			fields: []*Field{
				{
					colName: "id",
					goName: "Id",
					typ: reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					goName: "LastName",
					typ: reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					goName: "Age",
					typ: reflect.TypeOf(int8(0)),
				},
			},
		},
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fieldMap: map[string]*Field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
			fields: []*Field{
				{
					colName: "first_name_t",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},
		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fieldMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},
		{
			name: "column only",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvaildTagContent("column"),
		},
		{
			name: "ignore tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc=abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},

		{
			name: "CustomeTableName",
			entity: &CustomeTableName{},
			wantModel: &Model{
				tableName: "custome_table_name_t",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},

		{
			name: "CustomeTableNamePtr",
			entity: &CustomeTableNamePtr{},
			wantModel: &Model{
				tableName: "custome_table_name_ptr_t",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},

		{
			name: "EmplyTableName",
			entity: &EmplyTableName{},
			wantModel: &Model{
				tableName: "emply_table_name",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName: "FirstName",
					typ: reflect.TypeOf(""),
				},
			},
		},
	}

	r := newRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tc.fields {
				fieldMap[f.goName] = f
				columnMap[f.colName] = f
			}
			tc.wantModel.fieldMap = fieldMap
			tc.wantModel.columnMap = columnMap

			assert.Equal(t, tc.wantModel, m)

			typ := reflect.TypeOf(tc.entity)
			cach, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, cach)

		})
	}
}


type  CustomeTableName struct {
	FirstName string
}

func (c CustomeTableName) TableName() string {
	return "custome_table_name_t"
}


type  CustomeTableNamePtr struct {
	FirstName string
}

func (c *CustomeTableNamePtr) TableName() string {
	return "custome_table_name_ptr_t"
}

type EmplyTableName struct {
	FirstName string
}
func (e *EmplyTableName) TableName() string {
	return ""
}


func TestModelWithTableName(t *testing.T) {
	r := newRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_ttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_ttt", m.tableName)
}


func TestModelWithColumnName(t *testing.T) {
	testCases := []struct{
		name string
		field string
		colName string

		wantColName string
		wantErr error
	}{
		{
			name: "column name",
			field: "FirstName",
			colName: "first_name_ccc",
			wantColName: "first_name_ccc",
		},

		{
			name: "invalid column name",
			field: "XXX",
			colName: "first_name_ccc",
			wantErr: errs.NewErrUnkownField("XXX"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			m, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.fieldMap[tc.field]
			require.True(t, ok)
			assert.Equal(t, tc.wantColName, fd.colName)
		})
	}
}

