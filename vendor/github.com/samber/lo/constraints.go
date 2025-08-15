package lo

// Clonable defines a constraint of types having Clone() T method.
type Clonable[T any] interface {
	Clone() T
}
