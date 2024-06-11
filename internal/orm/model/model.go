package orm

import (
	"github.com/ExerciseCoding/template/internal/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagKeyColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

type Model struct {
	TableName string
	// 字段名-> 字段定义
	FieldMap map[string]*Field

	// 列名 -> 字段定义
	ColumnMap map[string]*Field
}

type ModelOpt func(m *Model) error

type Field struct {
	// 字段名
	GoName string

	// 列名
	ColName string

	// 代表类型
	Type reflect.Type

	Offset uintptr
}

func ModelWithTableName(tableName string) ModelOpt {
	return func(m *Model) error {
		m.TableName = tableName
		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnkownField(field)

		}
		fd.ColName = colName
		return nil
	}
}

// var models = map[reflect.Type]*model{}
//var defaultRegistry = &registry{
//	models: map[reflect.Type]*model{},
//}

// registry 代表的是元数据的注册中心
type registry struct {
	// 读写锁

	// models map[reflect.Type]*model
	models sync.Map
}

func NewRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	// 缺点：多个go 协程操作这个时会出现覆盖
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}

	return m.(*Model), nil
}

// 解决registry map并发安全的解法1

//type registry struct {
//	// 读写锁
//	lock sync.RWMutex
//	models map[reflect.Type]*model
//}
// 处理并发安全：double check
//func (r *registry) get1(val any) (*model, error) {
//	typ := reflect.TypeOf(val)
//	r.lock.RLock()
//	m, ok := r.models[typ]
//	r.lock.RUnlock()
//	if ok {
//		return m, nil
//	}
//	r.lock.Lock()
//	defer r.lock.Unlock()
//	m, ok = r.models[typ]
//	if ok {
//		return m, nil
//	}
//
//	m, err := r.parseModel(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models[typ] = m
//
//	return m, nil
//}

// 只支持一级指针
func (r *registry) Register(entity any, opts ...ModelOpt) (*Model, error) {
	typ := reflect.TypeOf(entity)
	//if typ.Kind() == reflect.Pointer {
	//	typ = typ.Elem()
	//}
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemTyp := typ.Elem()
	numField := elemTyp.NumField()
	fieldMap := make(map[string]*Field, numField)
	columnMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := elemTyp.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagKeyColumn]
		if colName == "" {
			// 用户没有设置
			colName = underscoreName(fd.Name)
		}
		fdMeta := &Field{
			GoName:  fd.Name,
			ColName: colName,
			// 字段类型
			Type:   fd.Type,
			Offset: fd.Offset,
		}
		fieldMap[fd.Name] = fdMeta
		columnMap[colName] = fdMeta
	}

	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(elemTyp.Name())
	}
	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
	}
	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)
	return res, nil
}

type User struct {
	Id uint64 `orm:"column=id, xxx=bbb"`
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")

	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvaildTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))

		} else {
			buf = append(buf, byte(v))
		}
	}

	return string(buf)
}
