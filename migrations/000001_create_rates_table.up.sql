CREATE TABLE IF NOT EXISTS rates (
    id BIGSERIAL PRIMARY KEY,
    trading_pair VARCHAR(20) NOT NULL,
    ask_price DECIMAL(20, 8) NOT NULL,
    bid_price DECIMAL(20, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index on trading_pair and created_at for efficient queries
CREATE INDEX IF NOT EXISTS idx_rates_trading_pair_created_at ON rates(trading_pair, created_at DESC);

-- Index created_at for time-based queries
CREATE INDEX IF NOT EXISTS idx_rates_created_at ON rates(created_at DESC); 