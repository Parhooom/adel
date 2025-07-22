package rabbitmq

import (
	"adel/internal/service/judge"
	"adel/internal/service/postgres"
	"encoding/json"
	"log"
	"sync"
)

type Worker struct {
	rabbitMQClient  *RabbitMQClient
	submissionStore postgres.SubmissionStore
	problemStore    postgres.ProblemStore
	numWorkers      int
	judgeClient     *judge.JudgeService
	logger          *log.Logger
	wg              *sync.WaitGroup
}

func NewRabbitMQWorker(rabbitMQClient *RabbitMQClient, submissionStore postgres.SubmissionStore, problemStore postgres.ProblemStore, numWorkers int, judgeClient *judge.JudgeService, logger *log.Logger) *Worker {
	return &Worker{
		rabbitMQClient:  rabbitMQClient,
		submissionStore: submissionStore,
		problemStore:    problemStore,
		numWorkers:      numWorkers,
		judgeClient:     judgeClient,
		logger:          logger,
		wg:              &sync.WaitGroup{},
	}
}

func (w *Worker) Start() {
	w.wg.Add(w.numWorkers)
	for i := range w.numWorkers {
		go w.runWorker(i)
	}
}

func (w *Worker) Stop() {
	w.rabbitMQClient.Close()
	w.wg.Wait()
}

func (w *Worker) runWorker(id int) {
	defer w.wg.Done()

	ch, msgs, err := w.rabbitMQClient.ConsumeSubmissionJobs()
	if err != nil {
		w.logger.Printf("worker %d: failed to start consuming: %v", id, err)
		return
	}
	defer ch.Close()

	for d := range msgs {
		w.logger.Printf("worker %d: received submission job", id)

		var job SubmissionJob
		err := json.Unmarshal(d.Body, &job)
		if err != nil {
			w.logger.Printf("worker %d: error unmarshaling job: %v. rejecting message.", id, err)
			d.Nack(false, false)
			continue
		}

		w.processSubmission(&job)

		d.Ack(false)
	}

	w.logger.Printf("worker %d: stopping", id)
}

func (w *Worker) processSubmission(job *SubmissionJob) {
	submission := &job.Submission
	problem := &job.Problem

	err := w.judgeClient.Judge(submission, problem)
	if err != nil {
		w.logger.Printf("failed to judge submission %d: %v", submission.ID, err)
		updateErr := w.submissionStore.UpdateSubmission(submission)
		if updateErr != nil {
			w.logger.Printf("failed to update submission %d with internal error status: %v", submission.ID, updateErr)
		}
		return
	}

	err = w.submissionStore.UpdateSubmission(submission)
	if err != nil {
		w.logger.Printf("failed to update submission %d, error: %v", submission.ID, err)
	}
}
