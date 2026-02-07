package http

import (
	"encoding/json"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/provider/pubsub"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	pubsub pubsub.PubSub
}

func getCorrelationID(c echo.Context) string {
	return c.Request().Header.Get("Correlation-ID")
}

func (h Handler) publishEvent(c echo.Context, topic event.Topic, e event.Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	msg := message.NewMessage(e.GetHeader().ID, payload)
	msg.Metadata.Set(string(event.MetadataKeyCorrelationID), getCorrelationID(c))
	msg.Metadata.Set(string(event.MetadataKeyType), string(topic))

	if err := h.pubsub.Publish(topic, msg); err != nil {
		return err
	}
	return nil
}
