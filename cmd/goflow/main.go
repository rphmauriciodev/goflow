package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rphmauriciodev/goflow/internal/ingestion"
	"github.com/rphmauriciodev/goflow/internal/platform"
	"github.com/rphmauriciodev/goflow/internal/processing"
)

func main() {
	platform.InitLogger()
	mux := http.NewServeMux()
	slog.Info("Motor do GoFlow a aquecer...")

	if err := godotenv.Load(); err != nil {
		slog.Warn("ficheiro .env não encontrado, a usar variáveis do sistema")
	}

	pool, err := platform.NewPostgresPool()
	if err != nil {
		slog.Error("falha crítica ao ligar ao banco de dados", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := platform.NewPostgresBatchRepository(pool)

	processor := &processing.SimpleProcessor{}

	orchestrator := processing.NewOrchestrator(5, processor)

	ingestionHandler := ingestion.NewHandler(repo, orchestrator)
	ingestionHandler.RegisterRoutes(mux)

	port := ":8080"

	slog.Info("servidor iniciado", "port", port)

	err = http.ListenAndServe(port, mux)

	if err != nil {
		slog.Error("Falha ao iniciar o servidor", "error", err)
		os.Exit(1)
	}
}
