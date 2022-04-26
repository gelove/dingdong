package notify

type Notifier interface {
	Name() string
	Send() error
}
