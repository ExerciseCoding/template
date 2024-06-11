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
	testCases := []struct {
		name string

		entity    any
		wantModel *Model
		fields    []*Field
		wantErr   error
		opts      []ModelOpt
	}{
		{
			name:   "test model",
			entity: TestModel{},
			wantModel: &Model{
				TableName: "test_model",
			},
			fields: []*Field{
				{
					ColName: "id",
					GoName:  "Id",
					Type:    reflect.TypeOf(int64(0)),
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Type:    reflect.TypeOf(&sql.NullString{}),
				},
				{
					ColName: "age",
					GoName:  "Age",
					Type:    reflect.TypeOf(int8(0)),
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
			name:    "map",
			entity:  map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "slice",
			entity:  []int{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "base types",
			entity:  0,
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
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tc.wantModel.FieldMap = fieldMap
			tc.wantModel.ColumnMap = columnMap
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name string

		entity    any
		fields    []*Field
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
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
					ColName: "id",
					GoName:  "Id",
					Type:    reflect.TypeOf(int64(0)),
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Type:    reflect.TypeOf(&sql.NullString{}),
				},
				{
					ColName: "age",
					GoName:  "Age",
					Type:    reflect.TypeOf(int8(0)),
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
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name_t",
					},
				},
			},
			fields: []*Field{
				{
					ColName: "first_name_t",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
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
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
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
				TableName: "tag_table",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
			},
		},

		{
			name:   "CustomeTableName",
			entity: &CustomeTableName{},
			wantModel: &Model{
				TableName: "custome_table_name_t",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
			},
		},

		{
			name:   "CustomeTableNamePtr",
			entity: &CustomeTableNamePtr{},
			wantModel: &Model{
				TableName: "custome_table_name_ptr_t",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
			},
		},

		{
			name:   "EmplyTableName",
			entity: &EmplyTableName{},
			wantModel: &Model{
				TableName: "emply_table_name",
				//fieldMap: map[string]*Field{
				//	"FirstName": {
				//		colName: "first_name",
				//	},
				//},
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
				},
			},
		},
	}

	r := NewRegistry()
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
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tc.wantModel.FieldMap = fieldMap
			tc.wantModel.ColumnMap = columnMap

			assert.Equal(t, tc.wantModel, m)

			typ := reflect.TypeOf(tc.entity)
			cach, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, cach)

		})
	}
}

type CustomeTableName struct {
	FirstName string
}

func (c CustomeTableName) TableName() string {
	return "custome_table_name_t"
}

type CustomeTableNamePtr struct {
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
	r := NewRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_ttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_ttt", m.TableName)
}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name    string
		field   string
		colName string

		wantColName string
		wantErr     error
	}{
		{
			name:        "column name",
			field:       "FirstName",
			colName:     "first_name_ccc",
			wantColName: "first_name_ccc",
		},

		{
			name:    "invalid column name",
			field:   "XXX",
			colName: "first_name_ccc",
			wantErr: errs.NewErrUnkownField("XXX"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRegistry()
			m, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.FieldMap[tc.field]
			require.True(t, ok)
			assert.Equal(t, tc.wantColName, fd.ColName)
		})
	}
}
