package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/raffreitas/codeflix-video-encoder/application/repositories"
	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/raffreitas/codeflix-video-encoder/framework/queue"
	"gorm.io/gorm"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChanel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChanel,
		RabbitMQ:         rabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS.")
	}

	for qtdProcess := range concurrency {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Domain, qtdProcess)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID %v. Error during the job: %v with video: %v. Error: %v",
			jobResult.Message.DeliveryMode, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error.Error())
	} else {
		log.Printf("MessageID %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error.Error())
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult) error {
	Mutex.Lock()
	jobJson, err := json.Marshal(jobResult.Job)
	Mutex.Unlock()

	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
