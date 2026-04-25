-- Tabela para os Lotes (Batches)
CREATE TABLE IF NOT EXISTS batches (
    id VARCHAR(50) PRIMARY KEY,
    source TEXT NOT NULL,
    type VARCHAR(10) NOT NULL,
    raw_data BYTEA,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para as Execuções
CREATE TABLE IF NOT EXISTS executions (
    id VARCHAR(50) PRIMARY KEY,
    batch_id VARCHAR(50) REFERENCES batches(id),
    status VARCHAR(20) NOT NULL,
    total_items INTEGER DEFAULT 0,
    processed_items INTEGER DEFAULT 0,
    failed_items INTEGER DEFAULT 0,
    start_time TIMESTAMP WITH TIME ZONE,
    duration INTERVAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);