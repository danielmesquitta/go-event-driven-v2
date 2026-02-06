package event

type Event interface {
	GetHeader() EventHeader
}
