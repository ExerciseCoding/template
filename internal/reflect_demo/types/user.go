package types

type User struct {
	Name string
	// 因为同属一个包，所以age还可以被测试访问到
	// 如果是不同包，就访问不到
	age int
}
