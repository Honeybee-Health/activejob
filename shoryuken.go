package activejob

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	uuid "github.com/satori/go.uuid"
)

var (
	MaxMessageSize = 262144

	MaxMessageSizeExceeded = errors.New("exceeded maximum supported message size for SQS")
)

type Shoryuken struct {
	Client    sqsiface.SQSAPI
	AwsConfig *aws.Config
}

func NewShoryukenWorker(cfgs ...*aws.Config) *Shoryuken {
	sess := session.Must(session.NewSession(cfgs...))
	sqsClient := sqs.New(sess)
	return &Shoryuken{
		Client:    sqsClient,
		AwsConfig: cfgs[0],
	}
}

func (s *Shoryuken) Enqueue(job *Job) (string, error) {
	jobID := uuid.NewV4().String()
	now := time.Now()

	job.JobID = jobID
	job.EnqueuedAt = now

	j, err := json.Marshal(job)
	if err != nil {
		return "", &EnqueueError{jobID, err}
	}

	msgSize := len(j)
	if msgSize >= MaxMessageSize {
		return "", &EnqueueError{jobID, MaxMessageSizeExceeded}
	}

	queueUrl, err := s.queueUrl(job.QueueName)
	if err != nil {
		return "", &EnqueueError{jobID, err}
	}

	msg := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(string(j)),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"shoryuken_class": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("ActiveJob::QueueAdapters::ShoryukenAdapter::JobWrapper"),
			},
		},
	}

	out, err := s.Client.SendMessage(msg)
	if err != nil {
		return "", &EnqueueError{jobID, err}
	}

	return *out.MessageId, nil
}

func (s *Shoryuken) queueUrl(queueName string) (string, error) {
	i, err := s.Client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	return *i.QueueUrl, nil
}
