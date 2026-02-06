package main

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type EventHeader struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	OccurredAt string `json:"occurred_at"`
}

func NewEventHeader(eventName string) EventHeader {
	return EventHeader{
		ID:         uuid.NewString(),
		EventName:  eventName,
		OccurredAt: time.Now().Format(time.RFC3339),
	}
}

type ProductOutOfStock struct {
	Header    EventHeader `json:"header"`
	ProductID string      `json:"product_id"`
}

type ProductBackInStock struct {
	Header    EventHeader `json:"header"`
	ProductID string      `json:"product_id"`
	Quantity  int         `json:"quantity"`
}

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) Publisher {
	return Publisher{
		pub: pub,
	}
}

func (p Publisher) PublishProductOutOfStock(productID string) error {
	event := ProductOutOfStock{
		Header:    NewEventHeader("ProductOutOfStock"),
		ProductID: productID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(event.Header.ID, payload)

	return p.pub.Publish("product-updates", msg)
}

func (p Publisher) PublishProductBackInStock(productID string, quantity int) error {
	event := ProductBackInStock{
		Header:    NewEventHeader("ProductBackInStock"),
		ProductID: productID,
		Quantity:  quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(event.Header.ID, payload)

	return p.pub.Publish("product-updates", msg)
}
