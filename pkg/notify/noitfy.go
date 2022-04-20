package notify

type Notifier interface {
	Name() string
	Send(title, body string) error
}
