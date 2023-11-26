serverMux struct里面的reject设成atomic类型的好处
```
type serverMux struct {
	reject atomic.Bool
	*http.ServeMux
}
```
- 能用atomic就用atomic
- 否则用读写锁
- 否则用写锁