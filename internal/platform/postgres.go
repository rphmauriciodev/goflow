package platform

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	_ = godotenv.Load()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("variável DATABASE_URL não definida")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	slog.Info("conexão com o postgres estabelecida com sucesso")
	return pool, nil
}
