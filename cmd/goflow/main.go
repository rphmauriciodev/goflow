package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/rphmauriciodev/goflow/internal/ingestion"
	"github.com/rphmauriciodev/goflow/internal/platform"
	"github.com/rphmauriciodev/goflow/internal/processing"
)

func main() {
	platform.InitLogger()
	pool, _ := platform.NewPostgresPool()
	defer pool.Close()

	mux := http.NewServeMux()

	repo := platform.NewPostgresBatchRepository(pool)

	processor := &processing.SimpleProcessor{}

	orchestrator := processing.NewOrchestrator(3, processor)

	batch := &processing.Batch{ID: "LOTE-123"}
	orchestrator.Start(batch)

	ingestionHandler := ingestion.NewHandler(repo)
	ingestionHandler.RegisterRoutes(mux)

	port := ":3000"

	slog.Info("servidor iniciado", "port", 3000)

	err := http.ListenAndServe(port, mux)

	if err != nil {
		slog.Error("Falha ao iniciar o servidor", "error", err)
		os.Exit(1)
	}
}
