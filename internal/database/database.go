package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type RateRecord struct {
	ID          int64
	TradingPair string
	AskPrice    float64
	BidPrice    float64
	Timestamp   time.Time
	CreatedAt   time.Time
}

type Database struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewDatabase(dsn string, logger *zap.Logger) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:     db,
		logger: logger,
	}, nil
}

func (d *Database) SaveRate(record *RateRecord) error {
	query := `
		INSERT INTO rates (trading_pair, ask_price, bid_price, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := d.db.QueryRow(
		query,
		record.TradingPair,
		record.AskPrice,
		record.BidPrice,
		record.Timestamp,
		record.CreatedAt,
	).Scan(&record.ID)

	if err != nil {
		return fmt.Errorf("failed to save rate: %w", err)
	}

	d.logger.Info("Rate saved to database",
		zap.String("trading_pair", record.TradingPair),
		zap.Float64("ask_price", record.AskPrice),
		zap.Float64("bid_price", record.BidPrice),
		zap.Time("timestamp", record.Timestamp),
	)

	return nil
}

func (d *Database) GetLatestRate(tradingPair string) (*RateRecord, error) {
	query := `
		SELECT id, trading_pair, ask_price, bid_price, timestamp, created_at
		FROM rates
		WHERE trading_pair = $1
		ORDER BY created_at DESC
		LIMIT 1`

	record := &RateRecord{}
	err := d.db.QueryRow(query, tradingPair).Scan(
		&record.ID,
		&record.TradingPair,
		&record.AskPrice,
		&record.BidPrice,
		&record.Timestamp,
		&record.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no rate found for trading pair: %s", tradingPair)
		}
		return nil, fmt.Errorf("failed to get latest rate: %w", err)
	}

	return record, nil
}

func (d *Database) GetRatesByTimeRange(tradingPair string, start, end time.Time) ([]*RateRecord, error) {
	query := `
		SELECT id, trading_pair, ask_price, bid_price, timestamp, created_at
		FROM rates
		WHERE trading_pair = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC`

	rows, err := d.db.Query(query, tradingPair, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query rates: %w", err)
	}
	defer rows.Close()

	var records []*RateRecord
	for rows.Next() {
		record := &RateRecord{}
		err := rows.Scan(
			&record.ID,
			&record.TradingPair,
			&record.AskPrice,
			&record.BidPrice,
			&record.Timestamp,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rate record: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return records, nil
}

func (d *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return d.db.PingContext(ctx)
}

func (d *Database) Close() error {
	return d.db.Close()
}
