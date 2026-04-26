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
		UPDATE executions 
		SET processed_items = $1, 
		    failed_items = $2, 
		    status = $3, 
		    duration = $4
		WHERE id = $5
	`

	_, err = tx.Exec(ctx, query,
		exec.ProcessedItems,
		exec.FailedItems,
		exec.Status,
		exec.Duration,
		exec.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao executar update na transação: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação (commit): %w", err)
	}

	return nil
}
