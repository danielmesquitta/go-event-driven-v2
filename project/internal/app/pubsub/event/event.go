package event

type Event interface {
	// GetHeader returns the header of the event.
	GetHeader() EventHeader
}
