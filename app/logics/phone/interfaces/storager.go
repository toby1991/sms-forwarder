package interfaces

type Storager interface {
	Add(data map[string]string) error //@todo data-er
	Close() error
}
