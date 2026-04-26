package processing

import (
	"sync"
	"time"
)

type Batch struct {
	ID        string
	Source    string
	Type      string
	RawData   []byte
	CreatedAt time.Time
}

type Job struct {
	ID         string
	BatchID    string
	Payload    map[string]interface{}
	Status     string
	RetryCount int
}

type Execution struct {
	mu             sync.Mutex
	ID             string
	BatchID        string
	Status         string
	TotalItems     int
	ProcessedItems int
	FailedItems    int
	StartTime      time.Time
	Duration       time.Duration
}

type ExecutionSummary struct {
	TotalBatches   int
	TotalProcessed int
	TotalFailed    int
	StatusCounts   map[string]int
}

type Processor interface {
	Process(j *Job) error
}

type BatchRepository interface {
	Save(b *Batch) error
	GetByID(id string) (*Batch, error)
	UpdateExecutionStatus(exec *Execution) error
	GetSummary() (*ExecutionSummary, error)
}

func (e *Execution) IncrementSuccess() {
	defer e.mu.Unlock()

	e.mu.Lock()
	e.ProcessedItems++
}

func (e *Execution) SetDuration() {
	e.Duration = time.Since(e.StartTime)
}

func (e *Execution) IncrementFailure() {
	defer e.mu.Unlock()

	e.mu.Lock()
	e.FailedItems++
}

func (e *Execution) Finalize(status string) {
	defer e.mu.Unlock()

	e.mu.Lock()
	e.Status = status
	e.SetDuration()
}
