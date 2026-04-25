package processing

import (
	"fmt"
	"time"
)

type SimpleProcessor struct{}

func (s *SimpleProcessor) Process(j *Job) error {
	time.Sleep(500 * time.Millisecond)
	j.Status = "Success"
	fmt.Printf("Job %s concluído com sucesso.\n", j.ID)
	return nil
}
