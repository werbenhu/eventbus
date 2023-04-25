package eventbus

var (
	// singleton is a pointer to a unbuffered EventBus instance, which will be created when necessary.
	singleton *EventBus
)

func init() {
	if singleton == nil {

		// If singleton is nil, we create a new instance of EventBus using the New()
		// function and assign it to the singleton variable.
		singleton = New()
	}
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

// Close closes the singleton
func Close() {
	singleton.Close()
}
