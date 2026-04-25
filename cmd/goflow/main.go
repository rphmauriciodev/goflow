package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rphmauriciodev/goflow/internal/ingestion"
	"github.com/rphmauriciodev/goflow/internal/platform"
)

func main() {
	mux := http.NewServeMux()

	repo := platform.NewMemoryBatchRepository()

	ingestionHandler := ingestion.NewHandler(repo)
	ingestionHandler.RegisterRoutes(mux)

	port := ":3000"

	fmt.Printf("Servidor do GoFlow rodando na porta %s\n", port)

	err := http.ListenAndServe(port, mux)

	if err != nil {
		fmt.Printf("Falha ao iniciar o servidor: %s", err)
		os.Exit(1)
	}
}
