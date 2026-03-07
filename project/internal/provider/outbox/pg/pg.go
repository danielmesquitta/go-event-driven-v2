package pg

import (
	"context"
	"errors"
	"tickets/internal/app/pubsub/event"
	"tickets/internal/pkg/tx"
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
)

const ForwarderTopic = "events_to_forward"

type Outbox struct {
	*forwarder.Forwarder
}

func New(database *db.DB, eventBus eventbus.EventBus) *Outbox {
	postgresSub, err := watermillSQL.NewSubscriber(
		database.DB,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:    watermillSQL.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		},
		nil,
	)
	if err != nil {
		panic(err)
	}

	fwd, err := forwarder.NewForwarder(
		postgresSub,
		eventBus.Publisher(),
		nil,
		forwarder.Config{
			ForwarderTopic: ForwarderTopic,
			Router:         eventBus.Router(),
		},
	)
	if err != nil {
		panic(err)
	}

	return &Outbox{Forwarder: fwd}
}

var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

func (p *Outbox) Publish(ctx context.Context, ev event.Event) error {
	tx, ok := ctx.Value(tx.ContextKey).(*sqlx.Tx)
	if !ok {
		return errors.New("tx not found in context: outbox must use tx")
	}

	var (
		pub message.Publisher
		err error
	)
	pub, err = watermillSQL.NewPublisher(tx, watermillSQL.PublisherConfig{
		SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
	}, nil)
	if err != nil {
		return err
	}

	pub = eventbus.CorrelationPublisherDecorator{Publisher: pub}

	pub = forwarder.NewPublisher(pub, forwarder.PublisherConfig{
		ForwarderTopic: ForwarderTopic,
	})

	pub = eventbus.CorrelationPublisherDecorator{Publisher: pub}

	bus, err := cqrs.NewEventBusWithConfig(pub, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return err
	}

	return bus.Publish(ctx, ev)
}
