package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewGrinexService(t *testing.T) {
	config := &GrinexConfig{
		BaseURL:   "https://grinex.io",
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.Equal(t, logger, service.logger)
	assert.NotNil(t, service.client)
}

func TestGetUSDTRate_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v2/trades", r.URL.Path)
		assert.Equal(t, "usdtrub", r.URL.Query().Get("market"))
		assert.Equal(t, "100", r.URL.Query().Get("limit"))
		assert.Equal(t, "TestAgent/1.0", r.Header.Get("User-Agent"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		// Return mock response with trades
		response := `[
			{
				"id": 199135,
				"hid": "0cd78513afe",
				"price": "81.25",
				"volume": "3003.003",
				"funds": "243993.99",
				"market": "usdtrub",
				"created_at": "2025-07-28T21:22:14+03:00"
			},
			{
				"id": 199134,
				"hid": "03445b6fa19",
				"price": "81.20",
				"volume": "15470.7692",
				"funds": "1257000.0",
				"market": "usdtrub",
				"created_at": "2025-07-28T21:19:53+03:00"
			},
			{
				"id": 199133,
				"hid": "03445b6fa18",
				"price": "81.30",
				"volume": "1000.0",
				"funds": "81300.0",
				"market": "usdtrub",
				"created_at": "2025-07-28T21:18:00+03:00"
			}
		]`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	rate, err := service.GetUSDTRate(ctx)

	require.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Equal(t, "USDT/RUB", rate.TradingPair)
	assert.Equal(t, 81.30, rate.AskPrice) // Highest price
	assert.Equal(t, 81.20, rate.BidPrice) // Lowest price

	// Check that timestamp is parsed correctly from the first trade
	expectedTime, _ := time.Parse(time.RFC3339, "2025-07-28T21:22:14+03:00")
	assert.Equal(t, expectedTime, rate.Timestamp)
}

func TestGetUSDTRate_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `[]`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	rate, err := service.GetUSDTRate(ctx)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "no trades data available")
}

func TestGetUSDTRate_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	rate, err := service.GetUSDTRate(ctx)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
}

func TestGetUSDTRate_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	rate, err := service.GetUSDTRate(ctx)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "API request failed with status 500")
}

func TestHealthCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v2/markets", r.URL.Path)
		assert.Equal(t, "TestAgent/1.0", r.Header.Get("User-Agent"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	err := service.HealthCheck(ctx)

	assert.NoError(t, err)
}

func TestHealthCheck_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &GrinexConfig{
		BaseURL:   server.URL,
		Timeout:   30 * time.Second,
		UserAgent: "TestAgent/1.0",
	}

	logger := zap.NewNop()
	service := NewGrinexService(config, logger)

	ctx := context.Background()
	err := service.HealthCheck(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "health check failed with status 500")
}

func TestCalculatePricesFromTrades(t *testing.T) {
	logger := zap.NewNop()
	service := &GrinexService{logger: logger}

	trades := []GrinexTrade{
		{Price: "81.25"},
		{Price: "81.20"},
		{Price: "81.30"},
		{Price: "81.15"},
	}

	askPrice, bidPrice, err := service.calculatePricesFromTrades(trades)

	assert.NoError(t, err)
	assert.Equal(t, 81.30, askPrice) // Highest price
	assert.Equal(t, 81.15, bidPrice) // Lowest price
}

func TestCalculatePricesFromTrades_SinglePrice(t *testing.T) {
	logger := zap.NewNop()
	service := &GrinexService{logger: logger}

	trades := []GrinexTrade{
		{Price: "81.25"},
	}

	askPrice, bidPrice, err := service.calculatePricesFromTrades(trades)

	assert.NoError(t, err)
	assert.Equal(t, 81.25, askPrice)
	assert.Equal(t, 81.25, bidPrice) // Same as ask when only one price
}

func TestCalculatePricesFromTrades_InvalidPrice(t *testing.T) {
	logger := zap.NewNop()
	service := &GrinexService{logger: logger}

	trades := []GrinexTrade{
		{Price: "invalid"},
		{Price: "81.25"},
	}

	askPrice, bidPrice, err := service.calculatePricesFromTrades(trades)

	assert.NoError(t, err)
	assert.Equal(t, 81.25, askPrice)
	assert.Equal(t, 81.25, bidPrice)
}

func TestCalculatePricesFromTrades_EmptyTrades(t *testing.T) {
	logger := zap.NewNop()
	service := &GrinexService{logger: logger}

	trades := []GrinexTrade{}

	_, _, err := service.calculatePricesFromTrades(trades)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no trades to calculate prices from")
}
