package redisstream

import "tickets/internal/provider/eventbus"

func (e *EventBus) AddHandler(handler eventbus.EventHandler) error {
	_, err := e.processor.AddHandler(handler)
	if err != nil {
		return err
	}

	return nil
}
