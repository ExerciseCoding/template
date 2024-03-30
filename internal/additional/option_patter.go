package additional

import "errors"

type MyStructOption func(myStruct *MyStruct)

type MyStructOptionErr func(myStruct *MyStruct) error

type MyStruct struct {
	// 第一个部分是必须用户输入的字段
	id  		uint64
	name  		string

	// 第二个部分是可选的字段
	address 	string
	// 这里可以很多字段

	field1 int
	field2 int
}


func WithField1AndField2(field1, field2 int) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.field1 = field1
		myStruct.field2 = field2
	}
}

func WithAddress(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.address = address
	}
}


func WithAddressV1(address string) MyStructOptionErr {
	return func(myStruct *MyStruct) error {
		if address == "" {
			return errors.New("address is nil")
		}
		myStruct.address = address
		return nil
	}
}

func WithAddressV2(address string) MyStructOptionErr {
	return func(myStruct *MyStruct) error {
		if address == "" {
			panic("address is nil")
		}
		myStruct.address = address
		return nil
	}
}


// NewMyStruct 参数包含所有的必须用户指定的字段
func NewMyStruct(id uint64, name string, opts...MyStructOption) *MyStruct {
	// 构造必传的部分
	res := &MyStruct{
		id: id,
		name: name,
	}
	// 非指针
	//for _, opt := range opts {
	//	opt(&res)
	//}

	// 指针
	for _, opt := range opts {
		opt(res)
	}

	return res
}



func NewMyStructV1(id uint64, name string, opts...MyStructOptionErr) (*MyStruct, error) {
	// 构造必传的部分
	res := &MyStruct{
		id: id,
		name: name,
	}
	// 非指针
	//for _, opt := range opts {
	//	opt(&res)
	//}

	// 指针
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}
	}

	return res, nil
}



