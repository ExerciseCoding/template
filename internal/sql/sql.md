### driver.Value 和 sql.Scanner接口
场景：SQL默认支持的类型就是基础类型，[]byte和string，如果需要自定义类型，例如支持json类型
```
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
	
	上面查询中使用了结构体：FullName{FirstName: "A",LastName: "B"}，默认的SQL是不支持的，需要实现Value方法才行
	type FullName struct {
        FirstName string
        LastName string
    }

    func (f FullName) Value() (driver.Value, error) {
        return f.FirstName + f.LastName, nil
    }
    
    //使用JSON
    re, err := db.ExecContext(context.Background(), "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+"VALUES (?,?,?,?)",2,
		FullNameJson{FirstName: "A",LastName: "B"},18,"Jerry")
    type FullNameJson struct {
        FirstName string
        LastName string
    }
    func (f FullNameJson) Value() (driver.Value, error) {
        return f.FirstName + f.LastName, nil
    }
```


driver.Value 接口：(Go类型到数据库类型)
- driver.Value 接口用于将自定义数据类型转换为数据库可接受的值。
该接口定义了一个方法 Value() (Value, error)，其中 Value 是一个通用类型，通常是一个可以表示数据库支持的原始数据类型的值（如整数、字符串、二进制数据等）。
自定义类型必须实现 driver.Value 接口，以便能够将其值传递给数据库驱动程序，例如 MySQL、PostgreSQL 等。
sql.Scanner 接口：(数据库类型到Go类型)
- sql.Scanner 接口用于将从数据库检索的数据值扫描到自定义数据类型中。
该接口定义了一个方法 Scan(src interface{}) error，其中 src 是从数据库检索的原始值，通常是一个 driver.Value 类型。
自定义类型必须实现 sql.Scanner 接口，以便能够将从数据库中检索的数据值转换为自定义类型。