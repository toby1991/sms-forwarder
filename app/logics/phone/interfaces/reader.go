package interfaces

type Chipper interface {
	Close() error
	Read() (n int, b []byte, err error)
	Read2() //@todo do context
	Write(b []byte) error
	Error() <-chan error
	Bytes() <-chan []byte
}

