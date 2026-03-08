package ticket

import (
	"context"
	"fmt"
	"tickets/internal/app/pubsub/message/cmd"
	"tickets/internal/pkg/validator"
	"tickets/internal/provider/cmdbus"
)

type SendRefundCommand struct {
	cmdBus cmdbus.CommandBus
}

func NewSendRefundCommand(cmdBus cmdbus.CommandBus) *SendRefundCommand {
	return &SendRefundCommand{cmdBus: cmdBus}
}

type SendRefundCommandInput struct {
	TicketID string `json:"ticket_id" validate:"required"`
}

func (c *SendRefundCommand) Execute(ctx context.Context, in SendRefundCommandInput) error {
	err := validator.Validate(ctx, in)
	if err != nil {
		return fmt.Errorf("error validating ticket refund input: %w", err)
	}

	command := cmd.NewRefundTicket(ctx, in.TicketID)
	err = c.cmdBus.Send(ctx, command)
	if err != nil {
		return fmt.Errorf("error sending ticket refund command: %w", err)
	}

	return nil
}
