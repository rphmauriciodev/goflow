package processing

import "time"

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
	ID             string
	BatchID        string
	Status         string
	TotalItems     int
	ProcessedItems int
	FailedItems    int
	StartTime      time.Time
	Duration       time.Duration
}

type Processor interface {
	Process(j *Job) error
}

type BatchRepository interface {
	Save(b *Batch) error
	GetByID(id string) (*Batch, error)
}

func (e *Execution) IncrementSuccess() {
	e.ProcessedItems++
}

func (e *Execution) SetDuration() {
	e.Duration = time.Since(e.StartTime)
}
