package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/rphmauriciodev/goflow/internal/platform"
)

func main() {
	_ = godotenv.Load()
	pool, _ := platform.NewPostgresPool()
	defer pool.Close()

	repo := platform.NewPostgresBatchRepository(pool)
	stats, err := repo.GetSummary()
	if err != nil {
		fmt.Printf("Erro ao carregar dashboard: %v\n", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

	fmt.Println("\n📊 GOFLOW DASHBOARD - STATUS REALTIME")
	fmt.Println("======================================")
	fmt.Fprintf(w, "MÉTRICA\tVALOR\n")
	fmt.Fprintf(w, "Total de Lotes\t%d\n", stats.TotalBatches)
	fmt.Fprintf(w, "Itens Processados\t%d ✅\n", stats.TotalProcessed)
	fmt.Fprintf(w, "Itens Falhos\t%d ❌\n", stats.TotalFailed)

	fmt.Fprintln(w, "\nSTATUS\tQUANTIDADE")
	for status, count := range stats.StatusCounts {
		fmt.Fprintf(w, "%s\t%d\n", status, count)
	}

	w.Flush()
	fmt.Println("======================================\n")
}
