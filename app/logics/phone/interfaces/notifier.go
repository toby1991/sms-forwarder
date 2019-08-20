package interfaces

type Notifier interface {
	Notify(sender, content string) error
}
