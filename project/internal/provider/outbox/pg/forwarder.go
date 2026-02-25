package pg

import (
	"tickets/internal/provider/db"
	"tickets/internal/provider/eventbus"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
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
