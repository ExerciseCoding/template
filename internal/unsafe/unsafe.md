unsafe:
- 计算地址
- 计算偏移量
- 直接操作内存

```
type User struct {
	Name    string
	age     int32
	Alias   []byte
	Address string
}
```
```
64
Name:0
age:16
Alias:24
Address:48 
```


unsafe -- uintptr 使用误区
```
type FieldMeta struct {
    // offset 
    offset uintptr
}
```

如果使用uintptr来保存对象的起始地址，那么如果发生GC了，原本的代码会直接崩溃
例如在GC前，计算的entityAddr=OxAAAA，那么GC后因为复制的原因，实际上的地址变成了OxAABB

因为GC不会维护uintptr变量，所以entityAddr还是 OxAAAA, 这个时候再用OxAAAA作为起始地址去访问字段，就不知道访问的是什么东西了

uintptr 可以用于表达相对的量
例如字段偏移量，这个字段的偏移量是不管怎么GC都不会变。 如果怕出错，那么就只在进行地址运算的时候使用uintptr,其他时候都用unsafe.Pointer

