package processing

import (
	"fmt"
	"sync"
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

	for i := 0; i < o.workerCount; i++ {
		wg.Add(1)
		go o.worker(i, jobsChan, &wg)

	}

	go func() {
		for i := 1; i <= 10; i++ {
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

	fmt.Printf("Processamento do lote %s finalizado!\n", batch.ID)
}

func (o *Orchestrator) worker(id int, jobs <-chan *Job, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d pronto para trabalhar!\n", id)

	for job := range jobs {
		fmt.Printf("Worker %d processando job %s...\n", id, job.ID)

		if err := o.processor.Process(job); err != nil {
			fmt.Printf("Worker %d erro no job %s : %v\n", id, job.ID, err)
		}
	}
}
