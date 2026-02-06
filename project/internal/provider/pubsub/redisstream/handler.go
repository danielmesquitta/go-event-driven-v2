package redisstream

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"tickets/internal/app/pubsub/event"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (p *PubSub) AddHandler(
	subscribeTopic event.Topic,
	publishTopic event.Topic,
	handlerFunc message.HandlerFunc,
) *message.Handler {
	handlerName := handlerNameFromFunc(handlerFunc)

	sub, err := p.newSubscriber(fmt.Sprintf("%s-%s-group", handlerName, subscribeTopic))
	if err != nil {
		panic(err)
	}

	return p.router.AddHandler(handlerName, string(subscribeTopic), sub, string(publishTopic), p.publisher, handlerFunc)
}

func handlerNameFromFunc(fn any) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// "github.com/.../pubsub.(*Router).handleIssueReceipt-fm"
	//  → "handleIssueReceipt"

	// Get the part after the last '.'
	if idx := strings.LastIndex(name, "."); idx != -1 {
		name = name[idx+1:]
	}
	// Strip the "-fm" suffix from method values
	name = strings.TrimSuffix(name, "-fm")

	return name
}
