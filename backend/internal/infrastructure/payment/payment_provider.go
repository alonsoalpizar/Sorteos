package payment

import (
	"context"
)

// PaymentProvider defines the interface for payment processing
type PaymentProvider interface {
	// CreatePaymentIntent creates a new payment intent
	CreatePaymentIntent(ctx context.Context, input CreatePaymentIntentInput) (*PaymentIntent, error)

	// GetPaymentIntent retrieves a payment intent by ID
	GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error)

	// ConfirmPaymentIntent confirms a payment intent
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error)

	// CancelPaymentIntent cancels a payment intent
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error)

	// ConstructWebhookEvent constructs and verifies a webhook event
	ConstructWebhookEvent(payload []byte, signature string, secret string) (*WebhookEvent, error)
}

// CreatePaymentIntentInput represents input for creating a payment intent
type CreatePaymentIntentInput struct {
	Amount      int64             // Amount in cents (e.g., 1000 = $10.00)
	Currency    string            // e.g., "usd"
	Description string            // Payment description
	Metadata    map[string]string // Custom metadata
	CustomerID  string            // Optional Stripe customer ID
}

// PaymentIntent represents a payment intent
type PaymentIntent struct {
	ID           string
	Amount       int64
	Currency     string
	Status       string // requires_payment_method, requires_confirmation, requires_action, processing, succeeded, canceled
	ClientSecret string
	Description  string
	Metadata     map[string]string
	Created      int64
}

// WebhookEvent represents a webhook event from the payment provider
type WebhookEvent struct {
	Type string      // e.g., "payment_intent.succeeded"
	Data interface{} // Event data (type varies by event type)
}
