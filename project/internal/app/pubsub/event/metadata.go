package event

type MetadataKey string

const (
	MetadataKeyCorrelationID MetadataKey = "correlation_id"
	MetadataKeyType          MetadataKey = "type"
)
