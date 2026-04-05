package interfaces

// Notifier is an interface that defines a method for sending notifications with a message.
type Notifier interface {
	Notify(message string) error
}