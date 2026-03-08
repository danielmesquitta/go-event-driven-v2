package redisstream

import "tickets/internal/provider/cmdbus"

func (e *CommandBus) AddHandler(handler cmdbus.CommandHandler) error {
	_, err := e.processor.AddHandler(handler)
	if err != nil {
		return err
	}

	return nil
}
