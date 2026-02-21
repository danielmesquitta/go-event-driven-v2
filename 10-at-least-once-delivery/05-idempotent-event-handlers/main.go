package main

import "context"

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

type PaymentsHandler struct {
	repo  *PaymentsRepository
	saved map[string]struct{}
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{
		repo:  repo,
		saved: make(map[string]struct{}),
	}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	if _, ok := p.saved[event.PaymentID]; ok {
		return nil
	}
	err := p.repo.SavePaymentTaken(ctx, event)
	if err != nil {
		return err
	}
	p.saved[event.PaymentID] = struct{}{}
	return nil
}

type PaymentsRepository struct {
	payments []PaymentTaken
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	return &PaymentsRepository{}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	p.payments = append(p.payments, *event)
	return nil
}
