# An extremely simple interface for ActiveJob.

Currently this package support Shoryuken. Sidekiq is not implemented yet, PRs welcome.

### Install
```shell
go get -u github.com/Honeybee-Health/activejob
```

### Usage

```go
package main

import (
    "github.com/Honeybee-Health/activejob"
)

func main() {
    job := &activejob.Job{
        QueueName: "my_worker_queue",
        JobClass:  "MyHandlerJob",
        Arguments: []interface{}{
            map[string]interface{}{
                "some": "body data",
                "number": 1
            }
        },
    }

    awsConfig := &aws.Config{Region: aws.String("us-west-2")}
    jobClient := activejob.NewShoryukenWorker(awsConfig)
    jobID, err := jobClient.Enqueue(job)
    if err != nil {
        return err
    }
}
```
