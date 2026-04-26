package processing

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Orchestrator struct {
	workerCount int
	processor   Processor
}

func NewOrchestrator(workerCount int, p Processor) *Orchestrator {
	return &Orchestrator{
		workerCount: workerCount,
		processor:   p,
	}
}

func (o *Orchestrator) Start(batch *Batch, repo BatchRepository) {
	jobsChan := make(chan *Job, 100)

	var wg sync.WaitGroup

	exec := &Execution{
		ID:         fmt.Sprintf("exec_%d", time.Now().Unix()),
		BatchID:    batch.ID,
		Status:     "Processing",
		StartTime:  time.Now(),
		TotalItems: 10,
	}

	for i := 0; i < o.workerCount; i++ {
		wg.Add(1)
		go o.worker(i, jobsChan, &wg, exec)

	}

	go func() {
		for i := 1; i <= exec.TotalItems; i++ {
			job := &Job{
				ID:      fmt.Sprintf("%s-job-%d", batch.ID, i),
				BatchID: batch.ID,
				Status:  "Pending",
			}

			jobsChan <- job
		}
		close(jobsChan)
	}()

	wg.Wait()

	exec.Finalize("Completed")

	if err := repo.UpdateExecutionStatus(exec); err != nil {
		slog.Error("falha ao persistir status final", "error", err)
	}

	slog.Info("execução finalizada",
		"exec_id", exec.ID,
		"status", exec.Status,
		"sucessos", exec.ProcessedItems,
		"falhas", exec.FailedItems,
		"duracao_ms", exec.Duration.Milliseconds(),
	)
}

func (o *Orchestrator) worker(id int, jobs <-chan *Job, wg *sync.WaitGroup, exec *Execution) {
	defer wg.Done()

	l := slog.With(
		slog.Int("worker_id", id),
		slog.String("batch_id", exec.BatchID),
		slog.String("exec_id", exec.ID),
	)

	l.Info("worker iniciado e aguardando jobs")

	for job := range jobs {
		jl := l.With(slog.String("job_id", job.ID))

		jl.Debug("processando tentativa")
		success := false
		maxAttempts := 3

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err := o.processor.Process(job)

			if err == nil {
				success = true
				jl.Info("job processado com sucesso")
				break
			}

			job.RetryCount = attempt
			jl.Error("falha no processamento do job", "error", err, "retry_count", job.RetryCount)
			if attempt < maxAttempts {
				waitTime := time.Duration(1<<attempt) * time.Second
				time.Sleep(waitTime)
			}
		}

		if success {
			exec.IncrementSuccess()
		} else {
			jl.Error("O job falhou e excedeu as tentativas", "job", job.ID, "retry_count", job.RetryCount)
			exec.IncrementFailure()
		}
	}
}
