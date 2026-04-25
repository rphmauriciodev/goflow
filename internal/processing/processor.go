package processing

import (
	"errors"
	"math/rand/v2"
	"time"
)

type SimpleProcessor struct {
	maxRetries int
}

func (s *SimpleProcessor) Process(j *Job) error {
	if rand.Float32() < 0.3 {
		return errors.New("falha temporária de comunicação com API externa")
	}

	time.Sleep(100 * time.Millisecond)
	j.Status = "Success"
	return nil
}
