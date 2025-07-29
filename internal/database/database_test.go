package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSaveRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	record := &RateRecord{
		TradingPair: "USDT/RUB",
		AskPrice:    100.50,
		BidPrice:    100.40,
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO rates").
		WithArgs(record.TradingPair, record.AskPrice, record.BidPrice, record.Timestamp, record.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = database.SaveRate(record)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), record.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLatestRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	expectedRecord := &RateRecord{
		ID:          1,
		TradingPair: "USDT/RUB",
		AskPrice:    100.50,
		BidPrice:    100.40,
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "trading_pair", "ask_price", "bid_price", "timestamp", "created_at"}).
		AddRow(expectedRecord.ID, expectedRecord.TradingPair, expectedRecord.AskPrice, expectedRecord.BidPrice, expectedRecord.Timestamp, expectedRecord.CreatedAt)

	mock.ExpectQuery("SELECT id, trading_pair, ask_price, bid_price, timestamp, created_at FROM rates").
		WithArgs("USDT/RUB").
		WillReturnRows(rows)

	record, err := database.GetLatestRate("USDT/RUB")
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, expectedRecord.ID, record.ID)
	assert.Equal(t, expectedRecord.TradingPair, record.TradingPair)
	assert.Equal(t, expectedRecord.AskPrice, record.AskPrice)
	assert.Equal(t, expectedRecord.BidPrice, record.BidPrice)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLatestRate_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	mock.ExpectQuery("SELECT id, trading_pair, ask_price, bid_price, timestamp, created_at FROM rates").
		WithArgs("USDT/RUB").
		WillReturnError(sql.ErrNoRows)

	record, err := database.GetLatestRate("USDT/RUB")
	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "no rate found for trading pair: USDT/RUB")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRatesByTimeRange(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()

	expectedRecords := []*RateRecord{
		{
			ID:          1,
			TradingPair: "USDT/RUB",
			AskPrice:    100.50,
			BidPrice:    100.40,
			Timestamp:   time.Now(),
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			TradingPair: "USDT/RUB",
			AskPrice:    100.60,
			BidPrice:    100.50,
			Timestamp:   time.Now(),
			CreatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "trading_pair", "ask_price", "bid_price", "timestamp", "created_at"})
	for _, record := range expectedRecords {
		rows.AddRow(record.ID, record.TradingPair, record.AskPrice, record.BidPrice, record.Timestamp, record.CreatedAt)
	}

	mock.ExpectQuery("SELECT id, trading_pair, ask_price, bid_price, timestamp, created_at FROM rates").
		WithArgs("USDT/RUB", start, end).
		WillReturnRows(rows)

	records, err := database.GetRatesByTimeRange("USDT/RUB", start, end)
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, expectedRecords[0].ID, records[0].ID)
	assert.Equal(t, expectedRecords[1].ID, records[1].ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthCheck(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	mock.ExpectPing()

	err = database.HealthCheck()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthCheck_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	database := &Database{
		db:     db,
		logger: logger,
	}

	mock.ExpectPing().WillReturnError(sql.ErrConnDone)

	err = database.HealthCheck()
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
