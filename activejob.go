package activejob

import (
	"fmt"
	"time"
)

type ActiveJob interface {
	Enqueue(job *Job) (string, error)
}

type EnqueueError struct {
	JobID string
	Err   error
}

type Job struct {
	JobClass            string            `json:"job_class,omitempty"`
	JobID               string            `json:"job_id,omitempty"`
	ProviderJobID       *string           `json:"provider_job_id,omitempty"`
	QueueName           string            `json:"queue_name,omitempty"`
	Priority            *int              `json:"priority,omitempty"`
	Arguments           interface{}       `json:"arguments,omitempty"`
	Executions          int               `json:"executions,omitempty"`
	ExceptionExecutions map[string]string `json:"exception_executions,omitempty"`
	Locale              string            `json:"locale,omitempty"`
	Timezone            string            `json:"timezone,omitempty"`
	EnqueuedAt          time.Time         `json:"enqueued_at,omitempty"`
	EnqueueError        error             `json:"enqueue_error,omitempty"`
}

func (e *EnqueueError) Error() string {
	return fmt.Sprintf("enqueue error: %s - %v", e.JobID, e.Err)
}
