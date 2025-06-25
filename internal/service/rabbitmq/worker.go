package rabbitmq

import (
	"adel/internal/service/judge"
	"adel/internal/service/postgres"
	"log"
	"strconv"
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

	ch, msgs, err := w.rabbitMQClient.ConsumeSubmissionID()
	if err != nil {
		w.logger.Printf("worker %d: failed to start consuming: %v", id, err)
		return
	}
	defer ch.Close()

	for d := range msgs {
		w.logger.Printf("worker %d: received submission ID: %s", id, d.Body)

		submissionID, err := strconv.ParseInt(string(d.Body), 10, 64)
		if err != nil {
			w.logger.Printf("worker %d: error parsing submission ID: %v. rejecting message.", id, err)
			d.Nack(false, false)
			continue
		}

		w.processSubmission(submissionID)

		d.Ack(false)
	}

	w.logger.Printf("worker %d: stopping", id)
}

func (w *Worker) processSubmission(submissionID int64) {
	submission, err := w.submissionStore.GetSubmissionByID(submissionID)
	if err != nil {
		w.logger.Printf("failed to get submission %d from store: %v", submissionID, err)
		return
	}

	problem, err := w.problemStore.GetProblemByID(submission.ProblemID)
	if err != nil {
		w.logger.Printf("failed to get problem %d for submission %d from store: %v", submission.ProblemID, submissionID, err)
		return
	}

	err = w.judgeClient.Judge(submission, problem)
	if err != nil {
		w.logger.Printf("failed to judge submission %d: %v", submissionID, err)

		updateErr := w.submissionStore.UpdateSubmission(submission)
		if updateErr != nil {
			w.logger.Printf("failed to update submission %d with internal error status: %v", submissionID, updateErr)
		}

		return
	}

	err = w.submissionStore.UpdateSubmission(submission)
	if err != nil {
		w.logger.Printf("failed to update submission %d, error: %v", submissionID, err)
	}
}
