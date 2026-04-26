package platform

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rphmauriciodev/goflow/internal/processing"
)

type PostgresBatchRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresBatchRepository(pool *pgxpool.Pool) *PostgresBatchRepository {
	return &PostgresBatchRepository{
		pool: pool,
	}
}

func (r *PostgresBatchRepository) Save(b *processing.Batch) error {
	query := `
		INSERT INTO batches (id, source, type, raw_data, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			source = EXCLUDED.source,
			type = EXCLUDED.type,
			raw_data = EXCLUDED.raw_data;
	`

	_, err := r.pool.Exec(context.Background(), query,
		b.ID,
		b.Source,
		b.Type,
		b.RawData,
		b.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("falha ao persistir lote: %w", err)
	}

	return nil
}

func (r *PostgresBatchRepository) GetByID(id string) (*processing.Batch, error) {
	query := `SELECT id, source, type, raw_data, created_at FROM batches WHERE id = $1`

	var b processing.Batch
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&b.ID, &b.Source, &b.Type, &b.RawData, &b.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("lote %s não encontrado: %w", id, err)
	}

	return &b, nil
}

func (r *PostgresBatchRepository) UpdateExecutionStatus(exec *processing.Execution) error {
	ctx := context.Background()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer tx.Rollback(ctx)

	query := `
		INSERT INTO executions (id, batch_id, status, total_items, processed_items, failed_items, start_time, duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			processed_items = EXCLUDED.processed_items,
			failed_items = EXCLUDED.failed_items,
			duration = EXCLUDED.duration;
	`
	_, err = r.pool.Exec(ctx, query,
		exec.ID,
		exec.BatchID,
		exec.Status,
		exec.TotalItems,
		exec.ProcessedItems,
		exec.FailedItems,
		exec.StartTime,
		exec.Duration,
	)
	if err != nil {
		return fmt.Errorf("erro ao executar update na transação: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação (commit): %w", err)
	}

	return nil
}

func (r *PostgresBatchRepository) GetSummary() (*processing.ExecutionSummary, error) {
	summary := &processing.ExecutionSummary{
		StatusCounts: make(map[string]int),
	}

	query := `
		SELECT 
			COUNT(DISTINCT batch_id) as total_batches,
			COALESCE(SUM(processed_items),0) as total_processed,
			COALESCE(SUM(failed_items),0) as total_failed
		FROM executions
	`
	err := r.pool.QueryRow(context.Background(), query).Scan(
		&summary.TotalBatches,
		&summary.TotalProcessed,
		&summary.TotalFailed,
	)
	if err != nil {
		return nil, err
	}

	statusQuery := "SELECT status, COUNT(*) FROM executions GROUP BY status"
	rows, err := r.pool.Query(context.Background(), statusQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		summary.StatusCounts[status] = count
	}

	return summary, nil
}
