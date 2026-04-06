package interfaces

// Notifier is an interface that defines a method for sending notifications with a message.
type Notifier interface {

	// Notify sends a notification with the given message. It returns an error if the notification fails to send.
	Notify(message string) error

	// GetHost returns the host or URL associated with the notifier.
	GetHost() string
}
