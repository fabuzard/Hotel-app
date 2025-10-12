package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type BookingResponse struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	RoomID       uint    `json:"room_id"`
	CheckinDate  string  `json:"checkin_date"`
	CheckoutDate string  `json:"checkout_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
}

func GetBookingByID(bookingID uint, authHeader string) (*BookingResponse, error) {
	baseURL := os.Getenv("RESERVATION_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082" // fallback for local
	}

	url := fmt.Sprintf("%s/reservations/bookings/%d", baseURL, bookingID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Authorization header directly (already contains "Bearer " prefix)
	// If it doesn't have Bearer prefix, add it
	if authHeader != "" {
		if !strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = "Bearer " + authHeader
		}
		req.Header.Set("Authorization", authHeader)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call reservation service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reservation service returned status %d", resp.StatusCode)
	}

	var booking BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		return nil, fmt.Errorf("failed to decode booking response: %w", err)
	}

	return &booking, nil
}

func UpdateBookingStatus(bookingID uint, status string, authHeader string) error {
	baseURL := os.Getenv("RESERVATION_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082" // fallback for local
	}

	url := fmt.Sprintf("%s/reservations/webhooks/bookings/%d", baseURL, bookingID)

	payload, err := json.Marshal(map[string]string{
		"status": status,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set auth header (add Bearer prefix if missing)
	if authHeader != "" {
		if !strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = "Bearer " + authHeader
		}
		req.Header.Set("Authorization", authHeader)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call reservation service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("reservation service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
