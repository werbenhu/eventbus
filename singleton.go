package eventbus

var (
	// singleton is a pointer to a unbuffered EventBus instance, which will be created when necessary.
	singleton *EventBus
)

func init() {
	ResetSingleton()
}

// ResetSingleton resets the singleton object. If the singleton object is not nil,
// it first closes the old singleton, and then creates a new singleton instance.
func ResetSingleton() {
	if singleton != nil {
		singleton.Close()
	}
	singleton = New()
}

// Unsubscribe removes handler defined for a topic.
// Returns error if there are no handlers subscribed to the topic.
func Unsubscribe(topic string, handler any) error {
	return singleton.Unsubscribe(topic, handler)
}

// Subscribe subscribes to a topic, return an error if the handler is not a function.
// The handler must have two parameters: the first parameter must be a string,
// and the type of the handler's second parameter must be consistent with the type of the payload in `Publish()`
func Subscribe(topic string, handler any) error {
	return singleton.Subscribe(topic, handler)
}

// Publish triggers the handlers defined for a topic. The `payload` argument will be passed to the handler.
// The type of the payload must correspond to the second parameter of the handler in `Subscribe()`.
func Publish(topic string, payload any) error {
	return singleton.Publish(topic, payload)
}

// PublishSync is a synchronous version of Publish that triggers the handlers defined for a topic with the given payload.
// The type of the payload must correspond to the second parameter of the handler in `Subscribe()`.
func PublishSync(topic string, payload any) error {
	return singleton.Publish(topic, payload)
}

// Close closes the singleton instance of EventBus.
func Close() {
	if singleton != nil {
		singleton.Close()
		singleton = nil
	}
}
