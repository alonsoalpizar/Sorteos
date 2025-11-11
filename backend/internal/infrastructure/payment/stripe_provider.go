package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

var (
	ErrInvalidAmount    = errors.New("invalid payment amount")
	ErrInvalidCurrency  = errors.New("invalid currency")
	ErrWebhookSignature = errors.New("invalid webhook signature")
)

// StripeProvider implements PaymentProvider using Stripe
type StripeProvider struct {
	apiKey string
}

// NewStripeProvider creates a new Stripe payment provider
func NewStripeProvider(apiKey string) *StripeProvider {
	stripe.Key = apiKey
	return &StripeProvider{
		apiKey: apiKey,
	}
}

// CreatePaymentIntent creates a new Stripe Payment Intent
func (p *StripeProvider) CreatePaymentIntent(ctx context.Context, input CreatePaymentIntentInput) (*PaymentIntent, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if input.Currency == "" {
		input.Currency = "usd"
	}

	// Convert metadata to stripe format
	metadata := make(map[string]string)
	for k, v := range input.Metadata {
		metadata[k] = v
	}

	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(input.Amount),
		Currency:    stripe.String(input.Currency),
		Description: stripe.String(input.Description),
		Metadata:    metadata,
		// Automatic payment methods (card, etc.)
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	if input.CustomerID != "" {
		params.Customer = stripe.String(input.CustomerID)
	}

	// Create the payment intent
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe create payment intent error: %w", err)
	}

	return &PaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
		Description:  pi.Description,
		Metadata:     pi.Metadata,
		Created:      pi.Created,
	}, nil
}

// GetPaymentIntent retrieves a Stripe Payment Intent
func (p *StripeProvider) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe get payment intent error: %w", err)
	}

	return &PaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
		Description:  pi.Description,
		Metadata:     pi.Metadata,
		Created:      pi.Created,
	}, nil
}

// ConfirmPaymentIntent confirms a Stripe Payment Intent
func (p *StripeProvider) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	params := &stripe.PaymentIntentConfirmParams{}
	pi, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("stripe confirm payment intent error: %w", err)
	}

	return &PaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
		Description:  pi.Description,
		Metadata:     pi.Metadata,
		Created:      pi.Created,
	}, nil
}

// CancelPaymentIntent cancels a Stripe Payment Intent
func (p *StripeProvider) CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	params := &stripe.PaymentIntentCancelParams{}
	pi, err := paymentintent.Cancel(paymentIntentID, params)
	if err != nil {
		return nil, fmt.Errorf("stripe cancel payment intent error: %w", err)
	}

	return &PaymentIntent{
		ID:           pi.ID,
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
		Description:  pi.Description,
		Metadata:     pi.Metadata,
		Created:      pi.Created,
	}, nil
}

// ConstructWebhookEvent constructs and verifies a Stripe webhook event
func (p *StripeProvider) ConstructWebhookEvent(payload []byte, signature string, secret string) (*WebhookEvent, error) {
	event, err := webhook.ConstructEvent(payload, signature, secret)
	if err != nil {
		return nil, ErrWebhookSignature
	}

	return &WebhookEvent{
		Type: string(event.Type),
		Data: event.Data.Raw,
	}, nil
}
