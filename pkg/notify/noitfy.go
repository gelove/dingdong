package notify

type Notifier interface {
	Name() string
	Send(string, string) error
}
