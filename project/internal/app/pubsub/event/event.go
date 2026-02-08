package event

type Event interface {
	// GetHeader returns the header of the event.
	GetHeader() EventHeader

	// SetDefaults sets the default values for the event, in case they are not set.
	SetDefaults()
}
