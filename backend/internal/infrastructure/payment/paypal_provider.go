package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/plutov/paypal/v4"
)

var (
	ErrPayPalInit        = errors.New("failed to initialize PayPal client")
	ErrPayPalCreateOrder = errors.New("failed to create PayPal order")
	ErrPayPalGetOrder    = errors.New("failed to get PayPal order")
	ErrPayPalCapture     = errors.New("failed to capture PayPal order")
	ErrPayPalCancel      = errors.New("failed to cancel PayPal order")
)

// PayPalProvider implements PaymentProvider using PayPal
type PayPalProvider struct {
	client *paypal.Client
}

// NewPayPalProvider creates a new PayPal payment provider
func NewPayPalProvider(clientID, secret string, sandbox bool) (*PayPalProvider, error) {
	var client *paypal.Client
	var err error

	if sandbox {
		client, err = paypal.NewClient(clientID, secret, paypal.APIBaseSandBox)
	} else {
		client, err = paypal.NewClient(clientID, secret, paypal.APIBaseLive)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalInit, err)
	}

	// Get access token
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get access token: %v", ErrPayPalInit, err)
	}

	return &PayPalProvider{
		client: client,
	}, nil
}

// CreatePaymentIntent creates a new PayPal Order (equivalent to payment intent)
func (p *PayPalProvider) CreatePaymentIntent(ctx context.Context, input CreatePaymentIntentInput) (*PaymentIntent, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if input.Currency == "" {
		input.Currency = "USD"
	}

	// Convert amount from cents to decimal string
	amountStr := fmt.Sprintf("%.2f", float64(input.Amount)/100.0)

	// Create PayPal order
	order, err := p.client.CreateOrder(ctx, paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: input.Currency,
				Value:    amountStr,
			},
			Description: input.Description,
			CustomID:    input.Metadata["reservation_id"], // Store reservation ID
		},
	}, nil, &paypal.ApplicationContext{
		BrandName:          "Sorteos Platform",
		LandingPage:        "NO_PREFERENCE",
		UserAction:         "PAY_NOW",
		ShippingPreference: "NO_SHIPPING",
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalCreateOrder, err)
	}

	// Extract approval URL from links
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	// Convert metadata back to map
	metadata := make(map[string]string)
	for k, v := range input.Metadata {
		metadata[k] = v
	}
	metadata["approval_url"] = approvalURL

	return &PaymentIntent{
		ID:           order.ID,
		Amount:       input.Amount,
		Currency:     input.Currency,
		Status:       string(order.Status),
		ClientSecret: approvalURL, // Use approval URL as client secret for frontend
		Description:  input.Description,
		Metadata:     metadata,
		Created:      0, // PayPal doesn't return created timestamp in this format
	}, nil
}

// GetPaymentIntent retrieves a PayPal Order
func (p *PayPalProvider) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	order, err := p.client.GetOrder(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalGetOrder, err)
	}

	// Parse amount back to cents
	var amount int64
	if len(order.PurchaseUnits) > 0 && order.PurchaseUnits[0].Amount != nil {
		amountFloat, _ := strconv.ParseFloat(order.PurchaseUnits[0].Amount.Value, 64)
		amount = int64(amountFloat * 100)
	}

	var currency string
	if len(order.PurchaseUnits) > 0 && order.PurchaseUnits[0].Amount != nil {
		currency = order.PurchaseUnits[0].Amount.Currency
	}

	var description string
	if len(order.PurchaseUnits) > 0 {
		description = order.PurchaseUnits[0].Description
	}

	// Extract metadata
	metadata := make(map[string]string)
	if len(order.PurchaseUnits) > 0 {
		metadata["custom_id"] = order.PurchaseUnits[0].CustomID
	}

	// Extract approval URL
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}
	if approvalURL != "" {
		metadata["approval_url"] = approvalURL
	}

	return &PaymentIntent{
		ID:           order.ID,
		Amount:       amount,
		Currency:     currency,
		Status:       string(order.Status),
		ClientSecret: approvalURL,
		Description:  description,
		Metadata:     metadata,
		Created:      0,
	}, nil
}

// ConfirmPaymentIntent captures a PayPal Order (equivalent to confirming payment)
func (p *PayPalProvider) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	// Capture the order
	capture, err := p.client.CaptureOrder(ctx, paymentIntentID, paypal.CaptureOrderRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalCapture, err)
	}

	// Parse amount from the first capture in the first purchase unit
	var amount int64
	var currency string

	if len(capture.PurchaseUnits) > 0 {
		unit := capture.PurchaseUnits[0]
		// Get amount from the first payment capture
		if len(unit.Payments.Captures) > 0 {
			capturePayment := unit.Payments.Captures[0]
			amountFloat, _ := strconv.ParseFloat(capturePayment.Amount.Value, 64)
			amount = int64(amountFloat * 100)
			currency = capturePayment.Amount.Currency
		}
	}

	metadata := make(map[string]string)

	return &PaymentIntent{
		ID:           capture.ID,
		Amount:       amount,
		Currency:     currency,
		Status:       string(capture.Status),
		ClientSecret: "",
		Description:  "",
		Metadata:     metadata,
		Created:      0,
	}, nil
}

// CancelPaymentIntent cancels a PayPal Order (void/cancel is not directly supported, returns order info)
func (p *PayPalProvider) CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	// PayPal doesn't have a direct "cancel" endpoint for orders in CREATED status
	// Orders automatically expire after 3 hours if not approved
	// We'll just return the current order status
	order, err := p.client.GetOrder(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalCancel, err)
	}

	var amount int64
	var currency string
	if len(order.PurchaseUnits) > 0 && order.PurchaseUnits[0].Amount != nil {
		amountFloat, _ := strconv.ParseFloat(order.PurchaseUnits[0].Amount.Value, 64)
		amount = int64(amountFloat * 100)
		currency = order.PurchaseUnits[0].Amount.Currency
	}

	var description string
	if len(order.PurchaseUnits) > 0 {
		description = order.PurchaseUnits[0].Description
	}

	metadata := make(map[string]string)
	if len(order.PurchaseUnits) > 0 {
		metadata["custom_id"] = order.PurchaseUnits[0].CustomID
	}
	metadata["note"] = "PayPal orders auto-expire after 3 hours if not approved"

	return &PaymentIntent{
		ID:           order.ID,
		Amount:       amount,
		Currency:     currency,
		Status:       string(order.Status),
		ClientSecret: "",
		Description:  description,
		Metadata:     metadata,
		Created:      0,
	}, nil
}

// ConstructWebhookEvent constructs and verifies a PayPal webhook event
func (p *PayPalProvider) ConstructWebhookEvent(payload []byte, signature string, secret string) (*WebhookEvent, error) {
	// PayPal webhook verification is more complex and requires webhook ID
	// For now, we'll do a simple JSON unmarshal
	// In production, you should use PayPal's webhook verification API

	var event struct {
		EventType    string          `json:"event_type"`
		ResourceType string          `json:"resource_type"`
		Resource     json.RawMessage `json:"resource"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Map PayPal event types to our generic types
	// CHECKOUT.ORDER.APPROVED -> payment_intent.requires_confirmation
	// PAYMENT.CAPTURE.COMPLETED -> payment_intent.succeeded
	// PAYMENT.CAPTURE.DENIED -> payment_intent.payment_failed

	eventType := event.EventType
	switch event.EventType {
	case "CHECKOUT.ORDER.APPROVED":
		eventType = "payment_intent.requires_confirmation"
	case "PAYMENT.CAPTURE.COMPLETED":
		eventType = "payment_intent.succeeded"
	case "PAYMENT.CAPTURE.DENIED":
		eventType = "payment_intent.payment_failed"
	}

	return &WebhookEvent{
		Type: eventType,
		Data: event.Resource,
	}, nil
}
