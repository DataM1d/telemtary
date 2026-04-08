package repository

import (
	"context"
	"telemetry-engine/internal/telemetry"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresRepo) StoreBatch(ctx context.Context, metrics []telemetry.OptimizedMetric) error {
	batch := &pgx.Batch{}

	for _, m := range metrics {
		batch.Queue("INSERT INTO telemetry_data (server_id, metric_value, metric_type) VALUES ($1, $2, $3)",
			m.ID, m.Value, m.Type)
	}

	br := r.Pool.SendBatch(ctx, batch)
	return br.Close()
}
