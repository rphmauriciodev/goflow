package platform

import (
	"errors"
	"sync"

	"github.com/rphmauriciodev/goflow/internal/processing"
)

type MemoryBatchRepository struct {
	mu      sync.RWMutex
	batches map[string]*processing.Batch
}

func NewMemoryBatchRepository() *MemoryBatchRepository {
	return &MemoryBatchRepository{
		batches: make(map[string]*processing.Batch),
	}
}

func (r *MemoryBatchRepository) Save(b *processing.Batch) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.batches[b.ID] = b

	return nil
}

func (r *MemoryBatchRepository) GetByID(id string) (*processing.Batch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	batch, ok := r.batches[id]
	if !ok {
		return nil, errors.New("lote não encontrado")
	}
	return batch, nil
}
