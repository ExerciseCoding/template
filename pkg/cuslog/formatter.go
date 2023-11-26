package cuslog

type Formatter interface {
	// Maybe in async goroutines
	// write the result to buffer
	Format(entry *Entry) error
}