package processing

import (
	"fmt"
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

func (o *Orchestrator) Start(batch *Batch) {
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

	fmt.Printf("--- Resumo da Execução %s ---\n", exec.ID)
	fmt.Printf("Status: %s\n", exec.Status)
	fmt.Printf("Sucessos: %d | Falhas: %d\n", exec.ProcessedItems, exec.FailedItems)
	fmt.Printf("Duração: %v\n", exec.Duration)
	fmt.Printf("------------------------------\n")
}

func (o *Orchestrator) worker(id int, jobs <-chan *Job, wg *sync.WaitGroup, exec *Execution) {
	defer wg.Done()

	fmt.Printf("Worker %d pronto para trabalhar!\n", id)

	for job := range jobs {
		success := false
		maxAttempts := 3

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			fmt.Printf("Worker %d processando job %s...\n", id, job.ID)

			err := o.processor.Process(job)

			if err == nil {
				success = true
				break
			}

			job.RetryCount = attempt
			fmt.Printf("[Worker %d] Tentativa %d falhou para o Job %s: %v\n", id, attempt, job.ID, err)

			if attempt < maxAttempts {
				waitTime := time.Duration(1<<attempt) * time.Second
				time.Sleep(waitTime)
			}
		}

		if success {
			exec.IncrementSuccess()
		} else {
			fmt.Printf("[Worker %d] Job %s falhou após %d tentativas.\n", id, job.ID, maxAttempts)
			exec.IncrementFailure()
		}
	}
}
