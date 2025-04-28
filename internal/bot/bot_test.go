package bot

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchStockQuote_Success(t *testing.T) {
	// Create a fake HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.US,2025-04-28,19:30:00,170.00,175.00,169.00,172.00,1000000
`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	// Replace the real API base URL
	oldAPIBaseURL := apiBaseURL
	apiBaseURL = mockServer.URL
	defer func() { apiBaseURL = oldAPIBaseURL }() // Restore after test

	price, err := fetchStockQuote("AAPL")

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	if price != "172.00" {
		t.Errorf("expected price 172.00 but got %s", price)
	}
}

func TestFetchStockQuote_InvalidStock(t *testing.T) {
	// Fake server returns N/D (no data)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `Symbol,Date,Time,Open,High,Low,Close,Volume
INVALID.US,N/D,N/D,N/D,N/D,N/D,N/D,N/D
`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	oldAPIBaseURL := apiBaseURL
	apiBaseURL = mockServer.URL
	defer func() { apiBaseURL = oldAPIBaseURL }()

	_, err := fetchStockQuote("INVALID")

	if err == nil || !strings.Contains(err.Error(), "invalid stock code") {
		t.Errorf("expected error for invalid stock code, got: %v", err)
	}
}
