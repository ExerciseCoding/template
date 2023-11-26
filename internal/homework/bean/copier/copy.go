package copier

// Copier 复制数据
type Copier[Src any, Dst any] interface {
	CopyTo(Src,Dst) error
	Copy(Src) (Dst, error)
}