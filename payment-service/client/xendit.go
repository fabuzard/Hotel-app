package client

import (
	"context"
	"fmt"
	"os"
	"time"

	xendit "github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

// CreateXenditPaymentURL creates a payment invoice and returns the payment URL
// Simple function - just input booking details, get payment URL
// payerEmail is optional - pass empty string if not available
func CreateXenditPaymentURL(bookingID uint, amount float64, payerEmail string) (string, error) {
	// Get Xendit secret key from environment
	secretKey := os.Getenv("XENDIT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("XENDIT_SECRET_KEY environment variable is not set")
	}

	// Initialize Xendit client
	xndClient := xendit.NewClient(secretKey)

	// Generate unique external ID
	externalID := fmt.Sprintf("booking-%d-%d", bookingID, time.Now().Unix())

	// Create invoice request
	invoiceReq := *invoice.NewCreateInvoiceRequest(externalID, amount)

	// Only set payer email if provided
	if payerEmail != "" {
		invoiceReq.SetPayerEmail(payerEmail)
	}

	invoiceReq.SetDescription(fmt.Sprintf("Payment for Booking #%d", bookingID))

	// Call Xendit API
	ctx := context.Background()
	resp, httpResp, err := xndClient.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(invoiceReq).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create Xendit invoice: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		return "", fmt.Errorf("xendit API returned status %d", httpResp.StatusCode)
	}

	// Get invoice URL from response
	if resp.InvoiceUrl == "" {
		return "", fmt.Errorf("invoice URL not returned from Xendit")
	}

	return resp.InvoiceUrl, nil
}
