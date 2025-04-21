package services

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/raffreitas/codeflix-video-encoder/framework/utils"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, workerId int) {
	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
}
