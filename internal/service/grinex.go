package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// GrinexConfig holds configuration for the Grinex API
type GrinexConfig struct {
	BaseURL   string
	UserAgent string
	Timeout   time.Duration
}

// Rate represents a trading rate from Grinex
type Rate struct {
	TradingPair string
	AskPrice    float64
	BidPrice    float64
	Timestamp   time.Time
}

// GrinexTrade represents a trade from Grinex API
type GrinexTrade struct {
	ID        int64  `json:"id"`
	HID       string `json:"hid"`
	Price     string `json:"price"`
	Volume    string `json:"volume"`
	Funds     string `json:"funds"`
	Market    string `json:"market"`
	CreatedAt string `json:"created_at"`
}

type GrinexService struct {
	config *GrinexConfig
	client *http.Client
	logger *zap.Logger
}

func NewGrinexService(config *GrinexConfig, logger *zap.Logger) *GrinexService {
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &GrinexService{
		config: config,
		client: client,
		logger: logger,
	}
}

// GetUSDTRate fetches the current USDT rate from Grinex using recent trades
func (g *GrinexService) GetUSDTRate(ctx context.Context) (*Rate, error) {
	url := fmt.Sprintf("%s/api/v2/trades", g.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters for USDT/RUB market
	q := req.URL.Query()
	q.Add("market", "usdtrub")
	q.Add("limit", "100")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", g.config.UserAgent)
	req.Header.Set("Accept", "application/json")

	g.logger.Info("Fetching USDT rate from Grinex", zap.String("url", req.URL.String()))

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var trades []GrinexTrade
	if err := json.Unmarshal(body, &trades); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(trades) == 0 {
		return nil, fmt.Errorf("no trades data available")
	}

	// Calculate ask and bid prices from recent trades
	askPrice, bidPrice, err := g.calculatePricesFromTrades(trades)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate prices from trades: %w", err)
	}

	// Get the latest trade timestamp
	latestTrade := trades[0] // Assuming trades are sorted by time descending
	timestamp, err := time.Parse(time.RFC3339, latestTrade.CreatedAt)
	if err != nil {
		timestamp = time.Now() // Fallback to current time
	}

	rate := &Rate{
		TradingPair: "USDT/RUB",
		AskPrice:    askPrice,
		BidPrice:    bidPrice,
		Timestamp:   timestamp,
	}

	g.logger.Info("Successfully fetched USDT rate",
		zap.Float64("ask_price", rate.AskPrice),
		zap.Float64("bid_price", rate.BidPrice),
		zap.Time("timestamp", rate.Timestamp),
		zap.Int("trades_count", len(trades)),
	)

	return rate, nil
}

// calculatePricesFromTrades calculates ask and bid prices from recent trades
func (g *GrinexService) calculatePricesFromTrades(trades []GrinexTrade) (askPrice, bidPrice float64, err error) {
	if len(trades) == 0 {
		return 0, 0, fmt.Errorf("no trades to calculate prices from")
	}

	// Sort trades by price to find highest (ask) and lowest (bid) recent prices
	var prices []float64
	for _, trade := range trades {
		price, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			g.logger.Warn("Failed to parse trade price", zap.String("price", trade.Price), zap.Error(err))
			continue
		}
		prices = append(prices, price)
	}

	if len(prices) == 0 {
		return 0, 0, fmt.Errorf("no valid prices found in trades")
	}

	// Sort prices in descending order
	sort.Sort(sort.Reverse(sort.Float64Slice(prices)))

	// Use the highest price as ask and lowest as bid
	askPrice = prices[0]             // Highest price
	bidPrice = prices[len(prices)-1] // Lowest price

	// If we have very few trades, use the same price for both
	if len(prices) < 2 {
		bidPrice = askPrice
	}

	return askPrice, bidPrice, nil
}

// HealthCheck performs a health check on the Grinex API
func (g *GrinexService) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v2/markets", g.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("User-Agent", g.config.UserAgent)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}
